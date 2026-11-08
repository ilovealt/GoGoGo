package postgres

import (
	"context"

	"gogogo/internal/model"
)

func (s *Store) CreatePoint(ctx context.Context, point *model.MeasurementPoint) error {
	if err := normalizePoint(point); err != nil {
		return err
	}

	if point.EquipmentID != nil {
		if _, err := s.GetEquipment(ctx, *point.EquipmentID); err != nil {
			return err
		}
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.MeasurementPoint{}).
		Where("code = ? AND deleted_at IS NULL", point.Code).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("point code already exists")
	}

	if point.Metadata == nil {
		point.Metadata = map[string]any{}
	}
	at := now()
	point.CreatedAt = at
	point.UpdatedAt = at

	err := s.db.WithContext(ctx).Create(point).Error
	return duplicateAs(err, "point code already exists")
}

func (s *Store) ListPoints(ctx context.Context, equipmentID *uint64, pointType string, status *uint8) ([]model.MeasurementPoint, error) {
	query := s.db.WithContext(ctx).Model(&model.MeasurementPoint{}).Where("deleted_at IS NULL")
	if equipmentID != nil {
		query = query.Where("equipment_id = ?", *equipmentID)
	}
	if pointType = trim(pointType); pointType != "" {
		query = query.Where("point_type = ?", pointType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	points := make([]model.MeasurementPoint, 0)
	if err := query.Order("id ASC").Find(&points).Error; err != nil {
		return nil, err
	}
	return points, nil
}

func (s *Store) GetPoint(ctx context.Context, id uint64) (model.MeasurementPoint, error) {
	var point model.MeasurementPoint
	err := s.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&point).Error
	if err != nil {
		return model.MeasurementPoint{}, recordNotFound(err, "point not found")
	}
	return point, nil
}

func (s *Store) UpdatePoint(ctx context.Context, point model.MeasurementPoint) error {
	existing, err := s.GetPoint(ctx, point.ID)
	if err != nil {
		return err
	}

	if err := normalizePoint(&point); err != nil {
		return err
	}

	if point.EquipmentID != nil {
		if _, err := s.GetEquipment(ctx, *point.EquipmentID); err != nil {
			return err
		}
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.MeasurementPoint{}).
		Where("code = ? AND id <> ? AND deleted_at IS NULL", point.Code, point.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("point code already exists")
	}

	if point.Metadata == nil {
		point.Metadata = map[string]any{}
	}
	point.CreatedAt = existing.CreatedAt
	point.DeletedAt = existing.DeletedAt
	point.UpdatedAt = now()

	err = s.db.WithContext(ctx).Save(&point).Error
	return duplicateAs(err, "point code already exists")
}

func (s *Store) SoftDeletePoint(ctx context.Context, id uint64) error {
	at := now()
	result := s.db.WithContext(ctx).Model(&model.MeasurementPoint{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": at, "updated_at": at})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return notFound("point not found")
	}
	return nil
}

func (s *Store) ListPointsByPlant(ctx context.Context, plantID uint64) ([]model.MeasurementPoint, error) {
	points := make([]model.MeasurementPoint, 0)
	err := s.db.WithContext(ctx).
		Model(&model.MeasurementPoint{}).
		Select("measurement_points.*").
		Joins("JOIN equipments ON equipments.id = measurement_points.equipment_id AND equipments.deleted_at IS NULL").
		Joins("JOIN areas ON areas.id = equipments.area_id").
		Where("areas.plant_id = ? AND measurement_points.deleted_at IS NULL", plantID).
		Order("measurement_points.id ASC").
		Find(&points).Error
	if err != nil {
		return nil, err
	}
	return points, nil
}
