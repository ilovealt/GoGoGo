package service

import (
	"context"
	"strings"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/repository"
)

type DataSourceService struct {
	store repository.Store
}

func NewDataSourceService(store repository.Store) *DataSourceService {
	return &DataSourceService{store: store}
}

func (s *DataSourceService) Create(ctx context.Context, req request.SourceCreate) (model.DataSource, error) {
	status := valueOr(req.Status, model.StatusEnabled)
	if err := requireStatus(status); err != nil {
		return model.DataSource{}, err
	}
	config := req.Config
	if config == nil {
		config = map[string]any{}
	}
	source := model.DataSource{
		Name:     strings.TrimSpace(req.Name),
		Protocol: strings.ToLower(strings.TrimSpace(req.Protocol)),
		Endpoint: strings.TrimSpace(req.Endpoint),
		Config:   config,
		Status:   status,
	}
	if err := s.store.CreateDataSource(ctx, &source); err != nil {
		return model.DataSource{}, err
	}
	return source, nil
}

func (s *DataSourceService) List(ctx context.Context) ([]model.DataSource, error) {
	return s.store.ListDataSources(ctx)
}

func (s *DataSourceService) Get(ctx context.Context, id uint64) (model.DataSource, error) {
	return s.store.GetDataSource(ctx, id)
}

func (s *DataSourceService) Update(ctx context.Context, id uint64, req request.SourceUpdate) (model.DataSource, error) {
	source, err := s.store.GetDataSource(ctx, id)
	if err != nil {
		return model.DataSource{}, err
	}
	if req.Name != nil {
		source.Name = strings.TrimSpace(*req.Name)
	}
	if req.Protocol != nil {
		source.Protocol = strings.ToLower(strings.TrimSpace(*req.Protocol))
	}
	if req.Endpoint != nil {
		source.Endpoint = strings.TrimSpace(*req.Endpoint)
	}
	if req.Config != nil {
		source.Config = *req.Config
		if source.Config == nil {
			source.Config = map[string]any{}
		}
	}
	if req.Status != nil {
		if err := requireStatus(*req.Status); err != nil {
			return model.DataSource{}, err
		}
		source.Status = *req.Status
	}
	if source.Name == "" || source.Endpoint == "" {
		return model.DataSource{}, apperrors.New(apperrors.ErrInvalid, "data source name and endpoint are required")
	}
	if err := s.store.UpdateDataSource(ctx, source); err != nil {
		return model.DataSource{}, err
	}
	return source, nil
}

func (s *DataSourceService) Delete(ctx context.Context, id uint64) error {
	return s.store.DeleteDataSource(ctx, id)
}

// Test 目前只做参数校验和模拟结果，不发起真实网络连接。
func (s *DataSourceService) Test(ctx context.Context, id uint64) (map[string]any, error) {
	source, err := s.store.GetDataSource(ctx, id)
	if err != nil {
		return nil, err
	}
	if source.Endpoint == "" {
		return nil, apperrors.New(apperrors.ErrInvalid, "endpoint is empty")
	}
	return map[string]any{
		"ok":         true,
		"protocol":   source.Protocol,
		"endpoint":   source.Endpoint,
		"latency_ms": 0,
		"message":    "connectivity test is mocked",
	}, nil
}
