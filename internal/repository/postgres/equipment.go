package postgres

import (
	"context"

	"gogogo/internal/model"
)

func (s *Store) CreateEquipment(ctx context.Context, equipment *model.Equipment) error {
	equipment.Code = trim(equipment.Code)
	equipment.Name = trim(equipment.Name)
	equipment.EquipmentType = trim(equipment.EquipmentType)
	equipment.Description = trim(equipment.Description)
	if equipment.Code == "" || equipment.Name == "" {
		return invalid("equipment code and name are required")
	}
	if err := validateStatus(equipment.Status); err != nil {
		return err
	}

	if _, err := s.GetArea(ctx, equipment.AreaID); err != nil {
		return err
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Equipment{}).
		Where("area_id = ? AND code = ? AND deleted_at IS NULL", equipment.AreaID, equipment.Code).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("equipment code already exists in this area")
	}

	at := now()
	equipment.CreatedAt = at
	equipment.UpdatedAt = at

	err := s.db.WithContext(ctx).Create(equipment).Error
	return duplicateAs(err, "equipment code already exists in this area")
}

func (s *Store) ListEquipments(ctx context.Context, areaID *uint64, keyword string) ([]model.Equipment, error) {
	query := s.db.WithContext(ctx).Model(&model.Equipment{}).Where("deleted_at IS NULL")
	if areaID != nil {
		query = query.Where("area_id = ?", *areaID)
	}
	if keyword = trim(keyword); keyword != "" {
		pattern := likePattern(keyword)
		query = query.Where("(lower(code) LIKE ? ESCAPE '\\' OR lower(name) LIKE ? ESCAPE '\\')", pattern, pattern)
	}

	equipments := make([]model.Equipment, 0)
	if err := query.Order("id ASC").Find(&equipments).Error; err != nil {
		return nil, err
	}
	return equipments, nil
}

func (s *Store) GetEquipment(ctx context.Context, id uint64) (model.Equipment, error) {
	var equipment model.Equipment
	err := s.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&equipment).Error
	if err != nil {
		return model.Equipment{}, recordNotFound(err, "equipment not found")
	}
	return equipment, nil
}

func (s *Store) UpdateEquipment(ctx context.Context, equipment model.Equipment) error {
	existing, err := s.GetEquipment(ctx, equipment.ID)
	if err != nil {
		return err
	}

	equipment.Code = trim(equipment.Code)
	equipment.Name = trim(equipment.Name)
	equipment.EquipmentType = trim(equipment.EquipmentType)
	equipment.Description = trim(equipment.Description)
	if equipment.Code == "" || equipment.Name == "" {
		return invalid("equipment code and name are required")
	}
	if err := validateStatus(equipment.Status); err != nil {
		return err
	}

	if _, err := s.GetArea(ctx, equipment.AreaID); err != nil {
		return err
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Equipment{}).
		Where("area_id = ? AND code = ? AND id <> ? AND deleted_at IS NULL",
			equipment.AreaID, equipment.Code, equipment.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("equipment code already exists in this area")
	}

	equipment.CreatedAt = existing.CreatedAt
	equipment.DeletedAt = existing.DeletedAt
	equipment.UpdatedAt = now()

	err = s.db.WithContext(ctx).Save(&equipment).Error
	return duplicateAs(err, "equipment code already exists in this area")
}

func (s *Store) SoftDeleteEquipment(ctx context.Context, id uint64) error {
	at := now()
	result := s.db.WithContext(ctx).Model(&model.Equipment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": at, "updated_at": at})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return notFound("equipment not found")
	}
	return nil
}
