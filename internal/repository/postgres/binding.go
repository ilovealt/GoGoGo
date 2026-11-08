package postgres

import (
	"context"

	"gogogo/internal/model"
)

func (s *Store) validateBinding(ctx context.Context, binding model.PointBinding) error {
	if _, err := s.GetPoint(ctx, binding.PointID); err != nil {
		return err
	}
	if _, err := s.GetDataSource(ctx, binding.SourceID); err != nil {
		return err
	}
	if len(binding.AddressConfig) == 0 {
		return invalid("address_config is required")
	}
	return nil
}

func (s *Store) CreateBinding(ctx context.Context, binding *model.PointBinding) error {
	if err := s.validateBinding(ctx, *binding); err != nil {
		return err
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.PointBinding{}).
		Where("point_id = ? AND source_id = ?", binding.PointID, binding.SourceID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("binding already exists for this point and source")
	}

	if binding.ConversionConfig == nil {
		binding.ConversionConfig = map[string]any{}
	}
	at := now()
	binding.CreatedAt = at
	binding.UpdatedAt = at

	err := s.db.WithContext(ctx).Create(binding).Error
	return duplicateAs(err, "binding already exists for this point and source")
}

func (s *Store) ListBindings(ctx context.Context, pointID, sourceID *uint64) ([]model.PointBinding, error) {
	query := s.db.WithContext(ctx).Model(&model.PointBinding{})
	if pointID != nil {
		query = query.Where("point_id = ?", *pointID)
	}
	if sourceID != nil {
		query = query.Where("source_id = ?", *sourceID)
	}

	bindings := make([]model.PointBinding, 0)
	if err := query.Order("id ASC").Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

func (s *Store) GetBinding(ctx context.Context, id uint64) (model.PointBinding, error) {
	var binding model.PointBinding
	if err := s.db.WithContext(ctx).First(&binding, id).Error; err != nil {
		return model.PointBinding{}, recordNotFound(err, "binding not found")
	}
	return binding, nil
}

func (s *Store) UpdateBinding(ctx context.Context, binding model.PointBinding) error {
	existing, err := s.GetBinding(ctx, binding.ID)
	if err != nil {
		return err
	}
	if err := s.validateBinding(ctx, binding); err != nil {
		return err
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.PointBinding{}).
		Where("point_id = ? AND source_id = ? AND id <> ?", binding.PointID, binding.SourceID, binding.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("binding already exists for this point and source")
	}

	if binding.ConversionConfig == nil {
		binding.ConversionConfig = map[string]any{}
	}
	binding.CreatedAt = existing.CreatedAt
	binding.UpdatedAt = now()

	err = s.db.WithContext(ctx).Save(&binding).Error
	return duplicateAs(err, "binding already exists for this point and source")
}

func (s *Store) DeleteBinding(ctx context.Context, id uint64) error {
	result := s.db.WithContext(ctx).Delete(&model.PointBinding{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return notFound("binding not found")
	}
	return nil
}

func (s *Store) ListActiveBindings(ctx context.Context) ([]model.BindingView, error) {
	views := make([]model.BindingView, 0)
	err := s.db.WithContext(ctx).
		Table("v_active_point_bindings").
		Order("binding_id ASC").
		Scan(&views).Error
	if err != nil {
		return nil, err
	}
	return views, nil
}
