package service

import (
	"context"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/repository"
)

type BindingService struct {
	store repository.Store
}

func NewBindingService(store repository.Store) *BindingService {
	return &BindingService{store: store}
}

func (s *BindingService) Create(ctx context.Context, req request.BindingCreate) (model.PointBinding, error) {
	if req.PointID == 0 || req.SourceID == 0 {
		return model.PointBinding{}, apperrors.New(apperrors.ErrInvalid, "point_id and source_id are required")
	}
	if len(req.AddressConfig) == 0 {
		return model.PointBinding{}, apperrors.New(apperrors.ErrInvalid, "address_config is required")
	}
	conversion := req.ConversionConfig
	if conversion == nil {
		conversion = map[string]any{}
	}
	binding := model.PointBinding{
		PointID:          req.PointID,
		SourceID:         req.SourceID,
		AddressConfig:    req.AddressConfig,
		ConversionConfig: conversion,
		Enabled:          valueOr(req.Enabled, true),
	}
	if err := s.store.CreateBinding(ctx, &binding); err != nil {
		return model.PointBinding{}, err
	}
	return binding, nil
}

func (s *BindingService) List(ctx context.Context, pointID, sourceID *uint64) ([]model.PointBinding, error) {
	return s.store.ListBindings(ctx, pointID, sourceID)
}

func (s *BindingService) Get(ctx context.Context, id uint64) (model.PointBinding, error) {
	return s.store.GetBinding(ctx, id)
}

func (s *BindingService) Update(ctx context.Context, id uint64, req request.BindingUpdate) (model.PointBinding, error) {
	binding, err := s.store.GetBinding(ctx, id)
	if err != nil {
		return model.PointBinding{}, err
	}
	if req.AddressConfig != nil {
		binding.AddressConfig = *req.AddressConfig
	}
	if req.ConversionConfig != nil {
		binding.ConversionConfig = *req.ConversionConfig
	}
	if req.Enabled != nil {
		binding.Enabled = *req.Enabled
	}
	if len(binding.AddressConfig) == 0 {
		return model.PointBinding{}, apperrors.New(apperrors.ErrInvalid, "address_config is required")
	}
	if binding.ConversionConfig == nil {
		binding.ConversionConfig = map[string]any{}
	}
	if err := s.store.UpdateBinding(ctx, binding); err != nil {
		return model.PointBinding{}, err
	}
	return binding, nil
}

func (s *BindingService) Delete(ctx context.Context, id uint64) error {
	return s.store.DeleteBinding(ctx, id)
}

func (s *BindingService) Active(ctx context.Context) ([]model.BindingView, error) {
	return s.store.ListActiveBindings(ctx)
}
