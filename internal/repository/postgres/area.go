package postgres

import (
	"context"

	"gogogo/internal/model"
)

func (s *Store) CreateArea(ctx context.Context, area *model.Area) error {
	area.Code = trim(area.Code)
	area.Name = trim(area.Name)
	area.Description = trim(area.Description)
	if area.Code == "" || area.Name == "" {
		return invalid("area code and name are required")
	}
	if err := validateStatus(area.Status); err != nil {
		return err
	}

	if _, err := s.GetPlant(ctx, area.PlantID); err != nil {
		return err
	}
	if area.ParentID != nil {
		parent, err := s.GetArea(ctx, *area.ParentID)
		if err != nil {
			return err
		}
		if parent.PlantID != area.PlantID {
			return invalid("parent area must belong to the same plant")
		}
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Area{}).
		Where("plant_id = ? AND code = ?", area.PlantID, area.Code).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("area code already exists in this plant")
	}

	at := now()
	area.CreatedAt = at
	area.UpdatedAt = at

	err := s.db.WithContext(ctx).Create(area).Error
	return duplicateAs(err, "area code already exists in this plant")
}

func (s *Store) ListAreasByPlant(ctx context.Context, plantID uint64) ([]model.Area, error) {
	areas := make([]model.Area, 0)
	if err := s.db.WithContext(ctx).
		Where("plant_id = ?", plantID).
		Order("id ASC").
		Find(&areas).Error; err != nil {
		return nil, err
	}
	return areas, nil
}

func (s *Store) GetArea(ctx context.Context, id uint64) (model.Area, error) {
	var area model.Area
	if err := s.db.WithContext(ctx).First(&area, id).Error; err != nil {
		return model.Area{}, recordNotFound(err, "area not found")
	}
	return area, nil
}

func (s *Store) UpdateArea(ctx context.Context, area model.Area) error {
	existing, err := s.GetArea(ctx, area.ID)
	if err != nil {
		return err
	}

	area.Code = trim(area.Code)
	area.Name = trim(area.Name)
	area.Description = trim(area.Description)
	if area.Code == "" || area.Name == "" {
		return invalid("area code and name are required")
	}
	if err := validateStatus(area.Status); err != nil {
		return err
	}

	if _, err := s.GetPlant(ctx, area.PlantID); err != nil {
		return err
	}
	if area.ParentID != nil {
		if *area.ParentID == area.ID {
			return invalid("area cannot be its own parent")
		}
		parent, err := s.GetArea(ctx, *area.ParentID)
		if err != nil {
			return err
		}
		if parent.PlantID != area.PlantID {
			return invalid("parent area must belong to the same plant")
		}
		hasCycle, err := s.areaHasCycle(ctx, area.ID, *area.ParentID)
		if err != nil {
			return err
		}
		if hasCycle {
			return invalid("area hierarchy contains a cycle")
		}
	}

	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Area{}).
		Where("plant_id = ? AND code = ? AND id <> ?", area.PlantID, area.Code, area.ID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return conflict("area code already exists in this plant")
	}

	area.CreatedAt = existing.CreatedAt
	area.UpdatedAt = now()

	err = s.db.WithContext(ctx).Save(&area).Error
	return duplicateAs(err, "area code already exists in this plant")
}

func (s *Store) DeleteArea(ctx context.Context, id uint64) error {
	if _, err := s.GetArea(ctx, id); err != nil {
		return err
	}

	var children int64
	if err := s.db.WithContext(ctx).Model(&model.Area{}).
		Where("parent_id = ?", id).Count(&children).Error; err != nil {
		return err
	}
	if children > 0 {
		return conflict("area still has child areas")
	}

	// 外键不识别软删除，只要还有设备行引用（即使已软删除）就不能删除区域。
	var equipments int64
	if err := s.db.WithContext(ctx).Model(&model.Equipment{}).
		Where("area_id = ?", id).Count(&equipments).Error; err != nil {
		return err
	}
	if equipments > 0 {
		return conflict("area still has equipments")
	}

	result := s.db.WithContext(ctx).Delete(&model.Area{}, id)
	if result.Error != nil {
		return foreignKeyAs(result.Error, "area is still referenced")
	}
	if result.RowsAffected == 0 {
		return notFound("area not found")
	}
	return nil
}

// areaHasCycle 沿 parent_id 向上遍历，判断把 areaID 挂到 parentID 下是否成环。
func (s *Store) areaHasCycle(ctx context.Context, areaID, parentID uint64) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Raw(`
WITH RECURSIVE ancestors AS (
    SELECT id, parent_id, 1 AS depth
    FROM areas
    WHERE id = ?
    UNION ALL
    SELECT a.id, a.parent_id, anc.depth + 1
    FROM areas a
    JOIN ancestors anc ON a.id = anc.parent_id
    WHERE anc.depth < 100
)
SELECT COUNT(*) FROM ancestors WHERE id = ?`, parentID, areaID).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
