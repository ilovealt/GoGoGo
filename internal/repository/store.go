package repository

import (
	"context"
	"time"

	"gogogo/internal/model"
)

// Store 是仓储层接口，由 postgres 包提供 GORM/PostgreSQL 实现。
type Store interface {
	// 用户
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id uint64) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	CountUsers(ctx context.Context) (int, error)

	// 刷新令牌
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (model.RefreshToken, error)
	RotateRefreshToken(ctx context.Context, oldID uint64, next *model.RefreshToken) error
	RevokeRefreshToken(ctx context.Context, id uint64, at time.Time) error
	RevokeRefreshFamily(ctx context.Context, familyID string, at time.Time) error
	RevokeUserRefreshTokens(ctx context.Context, userID uint64, at time.Time) error
	DeleteExpiredRefreshTokens(ctx context.Context, now time.Time) (int, error)

	// 工厂
	CreatePlant(ctx context.Context, plant *model.Plant) error
	ListPlants(ctx context.Context) ([]model.Plant, error)
	GetPlant(ctx context.Context, id uint64) (model.Plant, error)
	UpdatePlant(ctx context.Context, plant model.Plant) error
	DeletePlant(ctx context.Context, id uint64) error

	// 区域
	CreateArea(ctx context.Context, area *model.Area) error
	ListAreasByPlant(ctx context.Context, plantID uint64) ([]model.Area, error)
	GetArea(ctx context.Context, id uint64) (model.Area, error)
	UpdateArea(ctx context.Context, area model.Area) error
	DeleteArea(ctx context.Context, id uint64) error

	// 设备
	CreateEquipment(ctx context.Context, equipment *model.Equipment) error
	ListEquipments(ctx context.Context, areaID *uint64, keyword string) ([]model.Equipment, error)
	GetEquipment(ctx context.Context, id uint64) (model.Equipment, error)
	UpdateEquipment(ctx context.Context, equipment model.Equipment) error
	SoftDeleteEquipment(ctx context.Context, id uint64) error

	// 测点
	CreatePoint(ctx context.Context, point *model.MeasurementPoint) error
	ListPoints(ctx context.Context, equipmentID *uint64, pointType string, status *uint8) ([]model.MeasurementPoint, error)
	GetPoint(ctx context.Context, id uint64) (model.MeasurementPoint, error)
	UpdatePoint(ctx context.Context, point model.MeasurementPoint) error
	SoftDeletePoint(ctx context.Context, id uint64) error
	ListPointsByPlant(ctx context.Context, plantID uint64) ([]model.MeasurementPoint, error)

	// 数据源
	CreateDataSource(ctx context.Context, source *model.DataSource) error
	ListDataSources(ctx context.Context) ([]model.DataSource, error)
	GetDataSource(ctx context.Context, id uint64) (model.DataSource, error)
	UpdateDataSource(ctx context.Context, source model.DataSource) error
	DeleteDataSource(ctx context.Context, id uint64) error

	// 测点绑定
	CreateBinding(ctx context.Context, binding *model.PointBinding) error
	ListBindings(ctx context.Context, pointID, sourceID *uint64) ([]model.PointBinding, error)
	GetBinding(ctx context.Context, id uint64) (model.PointBinding, error)
	UpdateBinding(ctx context.Context, binding model.PointBinding) error
	DeleteBinding(ctx context.Context, id uint64) error
	ListActiveBindings(ctx context.Context) ([]model.BindingView, error)

	// 时序数据
	AddSamples(ctx context.Context, samples []model.PointSample) error
	ListSamples(ctx context.Context, pointID uint64, start, end time.Time) ([]model.PointSample, error)
	GetLatest(ctx context.Context, pointID uint64) (model.PointLatestValue, error)
	ListLatestByPlant(ctx context.Context, plantID uint64) ([]model.PointLatestView, error)
}
