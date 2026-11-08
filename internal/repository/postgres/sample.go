package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gogogo/internal/model"
)

const sampleBatchSize = 500

type sampleKey struct {
	PointID uint64
	Ts      int64
}

func (s *Store) AddSamples(ctx context.Context, samples []model.PointSample) error {
	if len(samples) == 0 {
		return nil
	}

	at := now()
	pointIDs := make([]uint64, 0, len(samples))
	sourceIDs := make([]uint64, 0, len(samples))
	seenPoints := make(map[uint64]struct{}, len(samples))
	seenSources := make(map[uint64]struct{}, len(samples))
	for _, sample := range samples {
		if _, ok := seenPoints[sample.PointID]; !ok {
			seenPoints[sample.PointID] = struct{}{}
			pointIDs = append(pointIDs, sample.PointID)
		}
		if sample.SourceID != nil {
			if _, ok := seenSources[*sample.SourceID]; !ok {
				seenSources[*sample.SourceID] = struct{}{}
				sourceIDs = append(sourceIDs, *sample.SourceID)
			}
		}
	}

	pointRows := make([]model.MeasurementPoint, 0, len(pointIDs))
	if err := s.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", pointIDs).
		Find(&pointRows).Error; err != nil {
		return err
	}
	if len(pointRows) != len(pointIDs) {
		return notFound("point not found")
	}
	points := make(map[uint64]model.MeasurementPoint, len(pointRows))
	for _, point := range pointRows {
		points[point.ID] = point
	}

	sources := make(map[uint64]struct{}, len(sourceIDs))
	if len(sourceIDs) > 0 {
		sourceRows := make([]model.DataSource, 0, len(sourceIDs))
		if err := s.db.WithContext(ctx).
			Where("id IN ?", sourceIDs).
			Find(&sourceRows).Error; err != nil {
			return err
		}
		if len(sourceRows) != len(sourceIDs) {
			return notFound("data source not found")
		}
		for _, source := range sourceRows {
			sources[source.ID] = struct{}{}
		}
	}

	// 同一 (point_id, ts) 在后到的数据中覆盖先到的数据，与内存实现一致。
	prepared := make([]model.PointSample, 0, len(samples))
	indexByKey := make(map[sampleKey]int, len(samples))
	for _, sample := range samples {
		point, ok := points[sample.PointID]
		if !ok {
			return notFound("point not found")
		}
		if err := validateSample(&sample, point); err != nil {
			return err
		}
		if sample.SourceID != nil {
			if _, ok := sources[*sample.SourceID]; !ok {
				return notFound("data source not found")
			}
		}
		if sample.CollectedAt.IsZero() {
			sample.CollectedAt = at
		}
		sample.Ts = sample.Ts.UTC()

		key := sampleKey{PointID: sample.PointID, Ts: sample.Ts.UnixNano()}
		if index, ok := indexByKey[key]; ok {
			prepared[index] = sample
		} else {
			indexByKey[key] = len(prepared)
			prepared = append(prepared, sample)
		}
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for start := 0; start < len(prepared); start += sampleBatchSize {
			end := start + sampleBatchSize
			if end > len(prepared) {
				end = len(prepared)
			}
			batch := prepared[start:end]
			// 同一测点同一时间戳只保留第一条，冲突时忽略。
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&batch).Error; err != nil {
				return err
			}
		}

		// 每个测点只保留时间戳最大的样本，时间相同则后到的胜出。
		latest := make([]model.PointLatestValue, 0)
		latestIndex := make(map[uint64]int, len(prepared))
		for _, sample := range prepared {
			if index, ok := latestIndex[sample.PointID]; ok {
				if sample.Ts.Before(latest[index].Ts) {
					continue
				}
				latest[index] = latestValue(sample, at)
				continue
			}
			latestIndex[sample.PointID] = len(latest)
			latest = append(latest, latestValue(sample, at))
		}

		if len(latest) == 0 {
			return nil
		}
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "point_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"ts", "value_type", "value_num", "value_text", "value_bool",
				"quality", "source_id", "updated_at",
			}),
			Where: clause.Where{Exprs: []clause.Expression{
				clause.Expr{SQL: "point_latest_values.ts <= EXCLUDED.ts"},
			}},
		}).Create(&latest).Error
	})
}

func latestValue(sample model.PointSample, at time.Time) model.PointLatestValue {
	return model.PointLatestValue{
		PointID:   sample.PointID,
		Ts:        sample.Ts,
		ValueType: sample.ValueType,
		ValueNum:  sample.ValueNum,
		ValueText: sample.ValueText,
		ValueBool: sample.ValueBool,
		Quality:   sample.Quality,
		SourceID:  sample.SourceID,
		UpdatedAt: at,
	}
}

func (s *Store) ListSamples(ctx context.Context, pointID uint64, start, end time.Time) ([]model.PointSample, error) {
	samples := make([]model.PointSample, 0)
	err := s.db.WithContext(ctx).
		Where("point_id = ? AND ts >= ? AND ts <= ?", pointID, start.UTC(), end.UTC()).
		Order("ts ASC").
		Find(&samples).Error
	if err != nil {
		return nil, err
	}
	return samples, nil
}

func (s *Store) GetLatest(ctx context.Context, pointID uint64) (model.PointLatestValue, error) {
	var latest model.PointLatestValue
	err := s.db.WithContext(ctx).Where("point_id = ?", pointID).First(&latest).Error
	if err != nil {
		return model.PointLatestValue{}, recordNotFound(err, "latest value not found")
	}
	return latest, nil
}

func (s *Store) ListLatestByPlant(ctx context.Context, plantID uint64) ([]model.PointLatestView, error) {
	views := make([]model.PointLatestView, 0)
	err := s.db.WithContext(ctx).
		Table("point_latest_values AS plv").
		Select(`mp.id AS point_id,
			mp.code AS point_code,
			mp.name AS point_name,
			mp.point_type,
			mp.unit,
			plv.ts,
			plv.value_type,
			plv.value_num,
			plv.value_text,
			plv.value_bool,
			plv.quality,
			plv.source_id`).
		Joins("JOIN measurement_points mp ON mp.id = plv.point_id").
		Joins("JOIN equipments e ON e.id = mp.equipment_id").
		Joins("JOIN areas a ON a.id = e.area_id").
		Where("mp.deleted_at IS NULL AND mp.status = ? AND e.deleted_at IS NULL AND a.plant_id = ?",
			model.StatusEnabled, plantID).
		Order("mp.id ASC").
		Scan(&views).Error
	if err != nil {
		return nil, err
	}
	return views, nil
}
