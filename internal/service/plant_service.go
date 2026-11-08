package service

import (
	"context"
	"strings"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/repository"
)

type PlantService struct {
	store repository.Store
}

func NewPlantService(store repository.Store) *PlantService {
	return &PlantService{store: store}
}

func (s *PlantService) Create(ctx context.Context, req request.PlantCreate) (model.Plant, error) {
	status := valueOr(req.Status, model.StatusEnabled)
	if err := requireStatus(status); err != nil {
		return model.Plant{}, err
	}
	timezone := strings.TrimSpace(req.Timezone)
	if timezone == "" {
		timezone = "UTC"
	}
	plant := model.Plant{
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Timezone:    timezone,
		Status:      status,
	}
	if err := s.store.CreatePlant(ctx, &plant); err != nil {
		return model.Plant{}, err
	}
	return plant, nil
}

func (s *PlantService) List(ctx context.Context) ([]model.Plant, error) {
	return s.store.ListPlants(ctx)
}

func (s *PlantService) Get(ctx context.Context, id uint64) (model.Plant, error) {
	return s.store.GetPlant(ctx, id)
}

func (s *PlantService) Update(ctx context.Context, id uint64, req request.PlantUpdate) (model.Plant, error) {
	plant, err := s.store.GetPlant(ctx, id)
	if err != nil {
		return model.Plant{}, err
	}
	if req.Code != nil {
		plant.Code = strings.TrimSpace(*req.Code)
	}
	if req.Name != nil {
		plant.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		plant.Description = strings.TrimSpace(*req.Description)
	}
	if req.Timezone != nil {
		plant.Timezone = strings.TrimSpace(*req.Timezone)
	}
	if req.Status != nil {
		if err := requireStatus(*req.Status); err != nil {
			return model.Plant{}, err
		}
		plant.Status = *req.Status
	}
	if plant.Code == "" || plant.Name == "" {
		return model.Plant{}, apperrors.New(apperrors.ErrInvalid, "plant code and name are required")
	}
	if err := s.store.UpdatePlant(ctx, plant); err != nil {
		return model.Plant{}, err
	}
	return plant, nil
}

func (s *PlantService) Delete(ctx context.Context, id uint64) error {
	return s.store.DeletePlant(ctx, id)
}
