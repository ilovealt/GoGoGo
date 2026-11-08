package service

import (
	"context"
	"strings"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/repository"
)

type PointService struct {
	store repository.Store
}

func NewPointService(store repository.Store) *PointService {
	return &PointService{store: store}
}

func (s *PointService) Create(ctx context.Context, req request.PointCreate) (model.MeasurementPoint, error) {
	valueType := valueOr(req.ValueType, model.ValueTypeNumeric)
	if valueType < model.ValueTypeNumeric || valueType > model.ValueTypeBool {
		return model.MeasurementPoint{}, apperrors.New(apperrors.ErrInvalid, "value_type must be 1, 2 or 3")
	}
	precision := valueOr(req.Precision, uint8(2))
	status := valueOr(req.Status, model.StatusEnabled)
	if err := requireStatus(status); err != nil {
		return model.MeasurementPoint{}, err
	}
	metadata := req.Metadata
	if metadata == nil {
		metadata = map[string]any{}
	}
	point := model.MeasurementPoint{
		EquipmentID: req.EquipmentID,
		Code:        strings.TrimSpace(req.Code),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		PointType:   strings.TrimSpace(req.PointType),
		ValueType:   valueType,
		Unit:        strings.TrimSpace(req.Unit),
		Precision:   precision,
		MinValue:    req.MinValue,
		MaxValue:    req.MaxValue,
		Deadband:    valueOr(req.Deadband, 0),
		Metadata:    metadata,
		Status:      status,
	}
	if err := s.store.CreatePoint(ctx, &point); err != nil {
		return model.MeasurementPoint{}, err
	}
	return point, nil
}

func (s *PointService) List(ctx context.Context, equipmentID *uint64, pointType string, status *uint8) ([]model.MeasurementPoint, error) {
	return s.store.ListPoints(ctx, equipmentID, strings.TrimSpace(pointType), status)
}

func (s *PointService) Get(ctx context.Context, id uint64) (model.MeasurementPoint, error) {
	return s.store.GetPoint(ctx, id)
}

func (s *PointService) Update(ctx context.Context, id uint64, req request.PointUpdate) (model.MeasurementPoint, error) {
	point, err := s.store.GetPoint(ctx, id)
	if err != nil {
		return model.MeasurementPoint{}, err
	}
	if req.ClearEquipment {
		point.EquipmentID = nil
	} else if req.EquipmentID != nil {
		point.EquipmentID = req.EquipmentID
	}
	if req.Code != nil {
		point.Code = strings.TrimSpace(*req.Code)
	}
	if req.Name != nil {
		point.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		point.Description = strings.TrimSpace(*req.Description)
	}
	if req.PointType != nil {
		point.PointType = strings.TrimSpace(*req.PointType)
	}
	if req.ValueType != nil {
		point.ValueType = *req.ValueType
	}
	if req.Unit != nil {
		point.Unit = strings.TrimSpace(*req.Unit)
	}
	if req.Precision != nil {
		point.Precision = *req.Precision
	}
	if req.ClearMinValue {
		point.MinValue = nil
	} else if req.MinValue != nil {
		point.MinValue = req.MinValue
	}
	if req.ClearMaxValue {
		point.MaxValue = nil
	} else if req.MaxValue != nil {
		point.MaxValue = req.MaxValue
	}
	if req.Deadband != nil {
		point.Deadband = *req.Deadband
	}
	if req.Metadata != nil {
		point.Metadata = *req.Metadata
	}
	if req.Status != nil {
		if err := requireStatus(*req.Status); err != nil {
			return model.MeasurementPoint{}, err
		}
		point.Status = *req.Status
	}
	if point.Code == "" || point.Name == "" || point.PointType == "" {
		return model.MeasurementPoint{}, apperrors.New(apperrors.ErrInvalid, "point code, name and point_type are required")
	}
	if err := s.store.UpdatePoint(ctx, point); err != nil {
		return model.MeasurementPoint{}, err
	}
	return point, nil
}

func (s *PointService) Delete(ctx context.Context, id uint64) error {
	return s.store.SoftDeletePoint(ctx, id)
}
