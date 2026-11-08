package postgres

import (
	"context"

	"gogogo/internal/model"
)

func (s *Store) CreateDataSource(ctx context.Context, source *model.DataSource) error {
	if err := normalizeDataSource(source); err != nil {
		return err
	}
	if source.Config == nil {
		source.Config = map[string]any{}
	}

	at := now()
	source.CreatedAt = at
	source.UpdatedAt = at

	err := s.db.WithContext(ctx).Create(source).Error
	return duplicateAs(err, "data source name already exists")
}

func (s *Store) ListDataSources(ctx context.Context) ([]model.DataSource, error) {
	sources := make([]model.DataSource, 0)
	if err := s.db.WithContext(ctx).Order("id ASC").Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

func (s *Store) GetDataSource(ctx context.Context, id uint64) (model.DataSource, error) {
	var source model.DataSource
	if err := s.db.WithContext(ctx).First(&source, id).Error; err != nil {
		return model.DataSource{}, recordNotFound(err, "data source not found")
	}
	return source, nil
}

func (s *Store) UpdateDataSource(ctx context.Context, source model.DataSource) error {
	existing, err := s.GetDataSource(ctx, source.ID)
	if err != nil {
		return err
	}

	if err := normalizeDataSource(&source); err != nil {
		return err
	}
	if source.Config == nil {
		source.Config = map[string]any{}
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.DataSource{}).
		Where("name = ? AND id <> ?", source.Name, source.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("data source name already exists")
	}

	source.CreatedAt = existing.CreatedAt
	source.UpdatedAt = now()

	err = s.db.WithContext(ctx).Save(&source).Error
	return duplicateAs(err, "data source name already exists")
}

func (s *Store) DeleteDataSource(ctx context.Context, id uint64) error {
	if _, err := s.GetDataSource(ctx, id); err != nil {
		return err
	}

	// point_bindings 由外键级联删除；point_samples/point_latest_values 的 source_id
	// 由外键 ON DELETE SET NULL 置空。
	result := s.db.WithContext(ctx).Delete(&model.DataSource{}, id)
	if result.Error != nil {
		return foreignKeyAs(result.Error, "data source is still referenced")
	}
	if result.RowsAffected == 0 {
		return notFound("data source not found")
	}
	return nil
}
