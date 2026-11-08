package bootstrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"gogogo/internal/config"
)

// OpenDatabase 按配置建立 GORM/PostgreSQL 连接并配置连接池。
func OpenDatabase(cfg config.DatabaseConfig, logLevel string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:         gormlogger.Default.LogMode(gormLogLevel(logLevel)),
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func gormLogLevel(level string) gormlogger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return gormlogger.Info
	case "error":
		return gormlogger.Error
	default:
		return gormlogger.Warn
	}
}
