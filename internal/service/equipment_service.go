package service

import (
	"context"
	"strings"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/repository"
)

type EquipmentService struct {
	store repository.Store
}

func NewEquipmentService(store repository.Store) *EquipmentService {
	return &EquipmentService{store: store}
}

func (s *EquipmentService) Create(ctx context.Context, req request.EquipmentCreate) (model.Equipment, error) {
	status := valueOr(req.Status, model.StatusEnabled)
	if err := requireStatus(status); err != nil {
		return model.Equipment{}, err
	}
	if req.AreaID == 0 {
		return model.Equipment{}, apperrors.New(apperrors.ErrInvalid, "area_id is required")
	}
	equipment := model.Equipment{
		AreaID:        req.AreaID,
		Code:          strings.TrimSpace(req.Code),
		Name:          strings.TrimSpace(req.Name),
		EquipmentType: strings.TrimSpace(req.EquipmentType),
		Description:   strings.TrimSpace(req.Description),
		Status:        status,
	}
	if err := s.store.CreateEquipment(ctx, &equipment); err != nil {
		return model.Equipment{}, err
	}
	return equipment, nil
}

func (s *EquipmentService) List(ctx context.Context, areaID *uint64, keyword string) ([]model.Equipment, error) {
	return s.store.ListEquipments(ctx, areaID, strings.TrimSpace(keyword))
}

func (s *EquipmentService) Get(ctx context.Context, id uint64) (model.Equipment, error) {
	return s.store.GetEquipment(ctx, id)
}

func (s *EquipmentService) Update(ctx context.Context, id uint64, req request.EquipmentUpdate) (model.Equipment, error) {
	equipment, err := s.store.GetEquipment(ctx, id)
	if err != nil {
		return model.Equipment{}, err
	}
	if req.AreaID != nil {
		equipment.AreaID = *req.AreaID
	}
	if req.Code != nil {
		equipment.Code = strings.TrimSpace(*req.Code)
	}
	if req.Name != nil {
		equipment.Name = strings.TrimSpace(*req.Name)
	}
	if req.EquipmentType != nil {
		equipment.EquipmentType = strings.TrimSpace(*req.EquipmentType)
	}
	if req.Description != nil {
		equipment.Description = strings.TrimSpace(*req.Description)
	}
	if req.Status != nil {
		if err := requireStatus(*req.Status); err != nil {
			return model.Equipment{}, err
		}
		equipment.Status = *req.Status
	}
	if equipment.Code == "" || equipment.Name == "" {
		return model.Equipment{}, apperrors.New(apperrors.ErrInvalid, "equipment code and name are required")
	}
	if err := s.store.UpdateEquipment(ctx, equipment); err != nil {
		return model.Equipment{}, err
	}
	return equipment, nil
}

func (s *EquipmentService) Delete(ctx context.Context, id uint64) error {
	return s.store.SoftDeleteEquipment(ctx, id)
}
