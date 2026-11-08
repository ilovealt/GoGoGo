package service

import (
	"context"
	"strings"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/repository"
)

type AreaService struct {
	store repository.Store
}

func NewAreaService(store repository.Store) *AreaService {
	return &AreaService{store: store}
}

func (s *AreaService) Create(ctx context.Context, req request.AreaCreate) (model.Area, error) {
	status := valueOr(req.Status, model.StatusEnabled)
	if err := requireStatus(status); err != nil {
		return model.Area{}, err
	}
	if req.PlantID == 0 {
		return model.Area{}, apperrors.New(apperrors.ErrInvalid, "plant_id is required")
	}
	area := model.Area{
		PlantID:     req.PlantID,
		ParentID:    req.ParentID,
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Status:      status,
	}
	if err := s.store.CreateArea(ctx, &area); err != nil {
		return model.Area{}, err
	}
	return area, nil
}

func (s *AreaService) ListByPlant(ctx context.Context, plantID uint64) ([]model.Area, error) {
	if _, err := s.store.GetPlant(ctx, plantID); err != nil {
		return nil, err
	}
	return s.store.ListAreasByPlant(ctx, plantID)
}

func (s *AreaService) Get(ctx context.Context, id uint64) (model.Area, error) {
	return s.store.GetArea(ctx, id)
}

func (s *AreaService) Update(ctx context.Context, id uint64, req request.AreaUpdate) (model.Area, error) {
	area, err := s.store.GetArea(ctx, id)
	if err != nil {
		return model.Area{}, err
	}
	if req.ClearParent {
		area.ParentID = nil
	} else if req.ParentID != nil {
		area.ParentID = req.ParentID
	}
	if req.Code != nil {
		area.Code = strings.TrimSpace(*req.Code)
	}
	if req.Name != nil {
		area.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		area.Description = strings.TrimSpace(*req.Description)
	}
	if req.Status != nil {
		if err := requireStatus(*req.Status); err != nil {
			return model.Area{}, err
		}
		area.Status = *req.Status
	}
	if area.Code == "" || area.Name == "" {
		return model.Area{}, apperrors.New(apperrors.ErrInvalid, "area code and name are required")
	}
	if err := s.store.UpdateArea(ctx, area); err != nil {
		return model.Area{}, err
	}
	return area, nil
}

func (s *AreaService) Delete(ctx context.Context, id uint64) error {
	return s.store.DeleteArea(ctx, id)
}
