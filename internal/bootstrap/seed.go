package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"gogogo/internal/model"
	"gogogo/internal/pkg/utils"
	"gogogo/internal/repository"
)

type SeedData struct {
	HasDemo      bool
	DemoPointIDs []uint64
	DemoSourceID uint64
}

// Seed 创建默认用户和一组演示数据。已有工厂时跳过演示数据，避免进程重启后重复写入。
func Seed(ctx context.Context, store repository.Store, logger *slog.Logger, withDemo bool) (*SeedData, error) {
	count, err := store.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		if err := createDefaultUsers(ctx, store); err != nil {
			return nil, err
		}
		logger.Warn("seeded default users: admin/admin123 and viewer/viewer123; change the passwords")
	}

	if !withDemo {
		return &SeedData{}, nil
	}

	plants, err := store.ListPlants(ctx)
	if err != nil {
		return nil, err
	}
	if len(plants) > 0 {
		return &SeedData{}, nil
	}

	plant := model.Plant{
		Code:        "P001",
		Name:        "Demo Plant",
		Description: "Seed data for the demo",
		Timezone:    "UTC",
		Status:      model.StatusEnabled,
	}
	if err := store.CreatePlant(ctx, &plant); err != nil {
		return nil, fmt.Errorf("seed plant: %w", err)
	}

	area := model.Area{
		PlantID:     plant.ID,
		Code:        "A001",
		Name:        "Assembly Area",
		Description: "Seed data for the demo",
		Status:      model.StatusEnabled,
	}
	if err := store.CreateArea(ctx, &area); err != nil {
		return nil, fmt.Errorf("seed area: %w", err)
	}

	equipment := model.Equipment{
		AreaID:        area.ID,
		Code:          "E001",
		Name:          "Compressor",
		EquipmentType: "compressor",
		Status:        model.StatusEnabled,
	}
	if err := store.CreateEquipment(ctx, &equipment); err != nil {
		return nil, fmt.Errorf("seed equipment: %w", err)
	}

	minTemp, maxTemp := 0.0, 100.0
	tempPoint := model.MeasurementPoint{
		EquipmentID: &equipment.ID,
		Code:        "TEMP-001",
		Name:        "Temperature",
		PointType:   "temperature",
		ValueType:   model.ValueTypeNumeric,
		Unit:        "C",
		Precision:   1,
		MinValue:    &minTemp,
		MaxValue:    &maxTemp,
		Deadband:    0.1,
		Metadata:    map[string]any{"demo": true},
		Status:      model.StatusEnabled,
	}
	if err := store.CreatePoint(ctx, &tempPoint); err != nil {
		return nil, fmt.Errorf("seed temperature point: %w", err)
	}

	minPressure, maxPressure := 0.0, 10.0
	pressurePoint := model.MeasurementPoint{
		EquipmentID: &equipment.ID,
		Code:        "PRESS-001",
		Name:        "Pressure",
		PointType:   "pressure",
		ValueType:   model.ValueTypeNumeric,
		Unit:        "bar",
		Precision:   2,
		MinValue:    &minPressure,
		MaxValue:    &maxPressure,
		Deadband:    0.01,
		Metadata:    map[string]any{"demo": true},
		Status:      model.StatusEnabled,
	}
	if err := store.CreatePoint(ctx, &pressurePoint); err != nil {
		return nil, fmt.Errorf("seed pressure point: %w", err)
	}

	source := model.DataSource{
		Name:     "Demo MQTT",
		Protocol: "mqtt",
		Endpoint: "mqtt://localhost:1883",
		Config:   map[string]any{"topic": "demo"},
		Status:   model.StatusEnabled,
	}
	if err := store.CreateDataSource(ctx, &source); err != nil {
		return nil, fmt.Errorf("seed data source: %w", err)
	}

	bindings := []model.PointBinding{
		{
			PointID:          tempPoint.ID,
			SourceID:         source.ID,
			AddressConfig:    map[string]any{"topic": "demo/temperature"},
			ConversionConfig: map[string]any{},
			Enabled:          true,
		},
		{
			PointID:          pressurePoint.ID,
			SourceID:         source.ID,
			AddressConfig:    map[string]any{"topic": "demo/pressure"},
			ConversionConfig: map[string]any{},
			Enabled:          true,
		},
	}
	for index := range bindings {
		if err := store.CreateBinding(ctx, &bindings[index]); err != nil {
			return nil, fmt.Errorf("seed binding: %w", err)
		}
	}

	now := time.Now().UTC()
	samples := make([]model.PointSample, 0, 60)
	for i := 29; i >= 0; i-- {
		ts := now.Add(-time.Duration(i) * time.Minute)
		tempValue := 22 + 3*math.Sin(float64(29-i)/5)
		pressureValue := 5 + 0.5*math.Cos(float64(29-i)/4)
		samples = append(samples,
			model.PointSample{
				PointID:     tempPoint.ID,
				Ts:          ts,
				ValueType:   model.ValueTypeNumeric,
				ValueNum:    &tempValue,
				Quality:     model.QualityGood,
				SourceID:    &source.ID,
				CollectedAt: now,
			},
			model.PointSample{
				PointID:     pressurePoint.ID,
				Ts:          ts,
				ValueType:   model.ValueTypeNumeric,
				ValueNum:    &pressureValue,
				Quality:     model.QualityGood,
				SourceID:    &source.ID,
				CollectedAt: now,
			},
		)
	}
	if err := store.AddSamples(ctx, samples); err != nil {
		return nil, fmt.Errorf("seed samples: %w", err)
	}

	logger.Info("seeded demo data",
		"plant_id", plant.ID,
		"area_id", area.ID,
		"equipment_id", equipment.ID,
		"temperature_point_id", tempPoint.ID,
		"pressure_point_id", pressurePoint.ID,
		"source_id", source.ID,
	)

	return &SeedData{
		HasDemo:      true,
		DemoPointIDs: []uint64{tempPoint.ID, pressurePoint.ID},
		DemoSourceID: source.ID,
	}, nil
}

func createDefaultUsers(ctx context.Context, store repository.Store) error {
	adminHash, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}
	admin := model.User{
		Username:     "admin",
		PasswordHash: adminHash,
		Nickname:     "Administrator",
		Role:         model.RoleAdmin,
		Status:       model.StatusEnabled,
	}
	if err := store.CreateUser(ctx, &admin); err != nil {
		return fmt.Errorf("seed admin user: %w", err)
	}

	viewerHash, err := utils.HashPassword("viewer123")
	if err != nil {
		return err
	}
	viewer := model.User{
		Username:     "viewer",
		PasswordHash: viewerHash,
		Nickname:     "Viewer",
		Role:         model.RoleViewer,
		Status:       model.StatusEnabled,
	}
	if err := store.CreateUser(ctx, &viewer); err != nil {
		return fmt.Errorf("seed viewer user: %w", err)
	}
	return nil
}
