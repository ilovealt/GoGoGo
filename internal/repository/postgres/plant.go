package postgres

import (
	"context"

	"gogogo/internal/model"
)

func (s *Store) CreatePlant(ctx context.Context, plant *model.Plant) error {
	plant.Code = trim(plant.Code)
	plant.Name = trim(plant.Name)
	plant.Description = trim(plant.Description)
	if plant.Code == "" || plant.Name == "" {
		return invalid("plant code and name are required")
	}
	if plant.Timezone == "" {
		plant.Timezone = "UTC"
	}
	if err := validateStatus(plant.Status); err != nil {
		return err
	}

	at := now()
	plant.CreatedAt = at
	plant.UpdatedAt = at

	err := s.db.WithContext(ctx).Create(plant).Error
	return duplicateAs(err, "plant code already exists")
}

func (s *Store) ListPlants(ctx context.Context) ([]model.Plant, error) {
	plants := make([]model.Plant, 0)
	if err := s.db.WithContext(ctx).Order("id ASC").Find(&plants).Error; err != nil {
		return nil, err
	}
	return plants, nil
}

func (s *Store) GetPlant(ctx context.Context, id uint64) (model.Plant, error) {
	var plant model.Plant
	if err := s.db.WithContext(ctx).First(&plant, id).Error; err != nil {
		return model.Plant{}, recordNotFound(err, "plant not found")
	}
	return plant, nil
}

func (s *Store) UpdatePlant(ctx context.Context, plant model.Plant) error {
	existing, err := s.GetPlant(ctx, plant.ID)
	if err != nil {
		return err
	}

	plant.Code = trim(plant.Code)
	plant.Name = trim(plant.Name)
	plant.Description = trim(plant.Description)
	if plant.Code == "" || plant.Name == "" {
		return invalid("plant code and name are required")
	}
	if plant.Timezone == "" {
		plant.Timezone = "UTC"
	}
	if err := validateStatus(plant.Status); err != nil {
		return err
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Plant{}).
		Where("code = ? AND id <> ?", plant.Code, plant.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("plant code already exists")
	}

	plant.CreatedAt = existing.CreatedAt
	plant.UpdatedAt = now()

	err = s.db.WithContext(ctx).Save(&plant).Error
	return duplicateAs(err, "plant code already exists")
}

func (s *Store) DeletePlant(ctx context.Context, id uint64) error {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Area{}).
		Where("plant_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("plant still has areas")
	}

	result := s.db.WithContext(ctx).Delete(&model.Plant{}, id)
	if result.Error != nil {
		return foreignKeyAs(result.Error, "plant still has areas")
	}
	if result.RowsAffected == 0 {
		return notFound("plant not found")
	}
	return nil
}
