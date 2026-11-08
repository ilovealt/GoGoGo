package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gogogo/internal/bootstrap"
	"gogogo/internal/collector"
	"gogogo/internal/config"
	"gogogo/internal/handler"
	"gogogo/internal/middleware"
	"gogogo/internal/repository"
	"gogogo/internal/repository/postgres"
	"gogogo/internal/router"
	"gogogo/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}

	logger := bootstrap.NewLogger(cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := bootstrap.OpenDatabase(cfg.Database, cfg.LogLevel)
	if err != nil {
		logger.Error("open database failed", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("get database handle failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			logger.Error("close database failed", "error", err)
		}
	}()

	if err := bootstrap.EnsureSamplePartitions(ctx, db, cfg.Database.PartitionMonthsAhead); err != nil {
		logger.Error("ensure point_samples partitions failed", "error", err)
		os.Exit(1)
	}

	store := postgres.New(db)
	seed, err := bootstrap.Seed(ctx, store, logger, cfg.SeedDemo)
	if err != nil {
		logger.Error("seed data failed", "error", err)
		os.Exit(1)
	}

	authService := service.NewAuthService(store, cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	userService := service.NewUserService(store)
	plantService := service.NewPlantService(store)
	areaService := service.NewAreaService(store)
	equipmentService := service.NewEquipmentService(store)
	pointService := service.NewPointService(store)
	sourceService := service.NewDataSourceService(store)
	bindingService := service.NewBindingService(store)
	sampleService := service.NewSampleService(store)

	handlers := router.Handlers{
		Auth:      handler.NewAuthHandler(authService),
		User:      handler.NewUserHandler(userService),
		Plant:     handler.NewPlantHandler(plantService),
		Area:      handler.NewAreaHandler(areaService),
		Equipment: handler.NewEquipmentHandler(equipmentService),
		Point:     handler.NewPointHandler(pointService),
		Source:    handler.NewSourceHandler(sourceService),
		Binding:   handler.NewBindingHandler(bindingService),
		Sample:    handler.NewSampleHandler(sampleService),
		Health:    handler.NewHealthHandler(),
	}

	engine := router.New(
		handlers,
		middleware.Auth(store, cfg.JWTSecret),
		cfg.CORSOrigins,
		logger,
	)

	if cfg.DemoCollector && seed.HasDemo {
		demo := collector.NewDemo(store, cfg.DemoInterval, logger)
		go demo.Run(ctx)
		logger.Info("demo collector started", "interval", cfg.DemoInterval.String(), "points", seed.DemoPointIDs)
	}

	go cleanupRefreshTokens(ctx, store, logger)
	go maintainSamplePartitions(ctx, store, cfg.Database.PartitionMonthsAhead, logger)

	if err := bootstrap.Run(ctx, cfg.Addr, engine, logger); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func cleanupRefreshTokens(ctx context.Context, store repository.Store, logger *slog.Logger) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			deleted, err := store.DeleteExpiredRefreshTokens(ctx, time.Now().UTC())
			if err != nil {
				logger.Error("cleanup refresh tokens", "error", err)
				continue
			}
			if deleted > 0 {
				logger.Info("cleaned up expired refresh tokens", "count", deleted)
			}
		}
	}
}

// maintainSamplePartitions 每天检查一次分区，保证跨月后仍有可用分区。
func maintainSamplePartitions(ctx context.Context, store *postgres.Store, monthsAhead int, logger *slog.Logger) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := bootstrap.EnsureSamplePartitions(ctx, store.DB(), monthsAhead); err != nil {
				logger.Error("ensure point_samples partitions", "error", err)
			}
		}
	}
}
