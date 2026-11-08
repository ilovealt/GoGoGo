package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

const devJWTSecret = "dev-only-secret-change-me-0123456789abcdef"

type Config struct {
	Addr          string
	JWTSecret     string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	LogLevel      string
	SeedDemo      bool
	DemoCollector bool
	DemoInterval  time.Duration
	CORSOrigins   []string
	Database      DatabaseConfig
}

// DatabaseConfig 保存 PostgreSQL 连接参数。
// 密码只通过 DB_PASSWORD 环境变量注入，不写入配置文件。
type DatabaseConfig struct {
	Host                 string
	Port                 int
	Name                 string
	User                 string
	Password             string
	SSLMode              string
	DSNOverride          string
	MaxOpenConns         int
	MaxIdleConns         int
	ConnMaxLifetime      time.Duration
	ConnMaxIdleTime      time.Duration
	PartitionMonthsAhead int
}

// DSN 返回 pgx 连接串；DB_DSN 优先级最高。
func (d DatabaseConfig) DSN() string {
	if strings.TrimSpace(d.DSNOverride) != "" {
		return strings.TrimSpace(d.DSNOverride)
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(d.Host, strconv.Itoa(d.Port)),
		Path:   "/" + d.Name,
	}
	if d.User != "" {
		u.User = url.UserPassword(d.User, d.Password)
	}
	query := url.Values{}
	query.Set("sslmode", d.SSLMode)
	query.Set("TimeZone", "UTC")
	u.RawQuery = query.Encode()
	return u.String()
}

type fileConfig struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	JWT struct {
		Secret     string `yaml:"secret"`
		AccessTTL  string `yaml:"access_ttl"`
		RefreshTTL string `yaml:"refresh_ttl"`
	} `yaml:"jwt"`
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Demo struct {
		Seed      *bool  `yaml:"seed"`
		Collector *bool  `yaml:"collector"`
		Interval  string `yaml:"interval"`
	} `yaml:"demo"`
	CORS struct {
		Origins []string `yaml:"origins"`
	} `yaml:"cors"`
	Database struct {
		Host                 string `yaml:"host"`
		Port                 *int   `yaml:"port"`
		Name                 string `yaml:"name"`
		User                 string `yaml:"user"`
		SSLMode              string `yaml:"sslmode"`
		MaxOpenConns         *int   `yaml:"max_open_conns"`
		MaxIdleConns         *int   `yaml:"max_idle_conns"`
		ConnMaxLifetime      string `yaml:"conn_max_lifetime"`
		ConnMaxIdleTime      string `yaml:"conn_max_idle_time"`
		PartitionMonthsAhead *int   `yaml:"partition_months_ahead"`
	} `yaml:"database"`
}

func Load() (Config, error) {
	cfg := Config{
		Addr:          ":8080",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
		LogLevel:      "info",
		SeedDemo:      true,
		DemoCollector: true,
		DemoInterval:  5 * time.Second,
		CORSOrigins:   []string{"*"},
		Database: DatabaseConfig{
			Host:                 "localhost",
			Port:                 5432,
			Name:                 "gogogo",
			User:                 "postgres",
			SSLMode:              "disable",
			MaxOpenConns:         50,
			MaxIdleConns:         10,
			ConnMaxLifetime:      time.Hour,
			ConnMaxIdleTime:      10 * time.Minute,
			PartitionMonthsAhead: 2,
		},
	}

	configFile := envOr("CONFIG_FILE", "config/config.yaml")
	if data, err := os.ReadFile(configFile); err == nil {
		var file fileConfig
		if err := yaml.Unmarshal(data, &file); err != nil {
			return Config{}, fmt.Errorf("parse config file %s: %w", configFile, err)
		}
		if strings.TrimSpace(file.Server.Addr) != "" {
			cfg.Addr = strings.TrimSpace(file.Server.Addr)
		}
		if strings.TrimSpace(file.JWT.Secret) != "" {
			cfg.JWTSecret = strings.TrimSpace(file.JWT.Secret)
		}
		if strings.TrimSpace(file.JWT.AccessTTL) != "" {
			parsed, err := time.ParseDuration(file.JWT.AccessTTL)
			if err != nil {
				return Config{}, fmt.Errorf("invalid jwt.access_ttl: %w", err)
			}
			cfg.AccessTTL = parsed
		}
		if strings.TrimSpace(file.JWT.RefreshTTL) != "" {
			parsed, err := time.ParseDuration(file.JWT.RefreshTTL)
			if err != nil {
				return Config{}, fmt.Errorf("invalid jwt.refresh_ttl: %w", err)
			}
			cfg.RefreshTTL = parsed
		}
		if strings.TrimSpace(file.Log.Level) != "" {
			cfg.LogLevel = strings.TrimSpace(file.Log.Level)
		}
		if file.Demo.Seed != nil {
			cfg.SeedDemo = *file.Demo.Seed
		}
		if file.Demo.Collector != nil {
			cfg.DemoCollector = *file.Demo.Collector
		}
		if strings.TrimSpace(file.Demo.Interval) != "" {
			parsed, err := time.ParseDuration(file.Demo.Interval)
			if err != nil {
				return Config{}, fmt.Errorf("invalid demo.interval: %w", err)
			}
			cfg.DemoInterval = parsed
		}
		if len(file.CORS.Origins) > 0 {
			cfg.CORSOrigins = file.CORS.Origins
		}
		if strings.TrimSpace(file.Database.Host) != "" {
			cfg.Database.Host = strings.TrimSpace(file.Database.Host)
		}
		if file.Database.Port != nil {
			cfg.Database.Port = *file.Database.Port
		}
		if strings.TrimSpace(file.Database.Name) != "" {
			cfg.Database.Name = strings.TrimSpace(file.Database.Name)
		}
		if strings.TrimSpace(file.Database.User) != "" {
			cfg.Database.User = strings.TrimSpace(file.Database.User)
		}
		if strings.TrimSpace(file.Database.SSLMode) != "" {
			cfg.Database.SSLMode = strings.TrimSpace(file.Database.SSLMode)
		}
		if file.Database.MaxOpenConns != nil {
			cfg.Database.MaxOpenConns = *file.Database.MaxOpenConns
		}
		if file.Database.MaxIdleConns != nil {
			cfg.Database.MaxIdleConns = *file.Database.MaxIdleConns
		}
		if strings.TrimSpace(file.Database.ConnMaxLifetime) != "" {
			parsed, err := time.ParseDuration(file.Database.ConnMaxLifetime)
			if err != nil {
				return Config{}, fmt.Errorf("invalid database.conn_max_lifetime: %w", err)
			}
			cfg.Database.ConnMaxLifetime = parsed
		}
		if strings.TrimSpace(file.Database.ConnMaxIdleTime) != "" {
			parsed, err := time.ParseDuration(file.Database.ConnMaxIdleTime)
			if err != nil {
				return Config{}, fmt.Errorf("invalid database.conn_max_idle_time: %w", err)
			}
			cfg.Database.ConnMaxIdleTime = parsed
		}
		if file.Database.PartitionMonthsAhead != nil {
			cfg.Database.PartitionMonthsAhead = *file.Database.PartitionMonthsAhead
		}
	} else if !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("read config file %s: %w", configFile, err)
	}

	cfg.Addr = envOr("APP_ADDR", cfg.Addr)
	cfg.JWTSecret = envOr("JWT_SECRET", cfg.JWTSecret)
	cfg.LogLevel = envOr("LOG_LEVEL", cfg.LogLevel)
	if value := os.Getenv("CORS_ORIGINS"); value != "" {
		cfg.CORSOrigins = splitCSV(value)
	}

	var err error
	if cfg.AccessTTL, err = envDuration("ACCESS_TOKEN_TTL", cfg.AccessTTL); err != nil {
		return Config{}, err
	}
	if cfg.RefreshTTL, err = envDuration("REFRESH_TOKEN_TTL", cfg.RefreshTTL); err != nil {
		return Config{}, err
	}
	if cfg.DemoInterval, err = envDuration("DEMO_INTERVAL", cfg.DemoInterval); err != nil {
		return Config{}, err
	}
	if cfg.SeedDemo, err = envBool("SEED_DEMO", cfg.SeedDemo); err != nil {
		return Config{}, err
	}
	if cfg.DemoCollector, err = envBool("DEMO_COLLECTOR", cfg.DemoCollector); err != nil {
		return Config{}, err
	}

	cfg.Database.Host = envOr("DB_HOST", cfg.Database.Host)
	cfg.Database.Name = envOr("DB_NAME", cfg.Database.Name)
	cfg.Database.User = envOr("DB_USER", cfg.Database.User)
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.SSLMode = envOr("DB_SSLMODE", cfg.Database.SSLMode)
	cfg.Database.DSNOverride = envOr("DB_DSN", cfg.Database.DSNOverride)
	if cfg.Database.Port, err = envInt("DB_PORT", cfg.Database.Port); err != nil {
		return Config{}, err
	}
	if cfg.Database.MaxOpenConns, err = envInt("DB_MAX_OPEN_CONNS", cfg.Database.MaxOpenConns); err != nil {
		return Config{}, err
	}
	if cfg.Database.MaxIdleConns, err = envInt("DB_MAX_IDLE_CONNS", cfg.Database.MaxIdleConns); err != nil {
		return Config{}, err
	}
	if cfg.Database.PartitionMonthsAhead, err = envInt("DB_PARTITION_MONTHS_AHEAD", cfg.Database.PartitionMonthsAhead); err != nil {
		return Config{}, err
	}
	if cfg.Database.ConnMaxLifetime, err = envDuration("DB_CONN_MAX_LIFETIME", cfg.Database.ConnMaxLifetime); err != nil {
		return Config{}, err
	}
	if cfg.Database.ConnMaxIdleTime, err = envDuration("DB_CONN_MAX_IDLE_TIME", cfg.Database.ConnMaxIdleTime); err != nil {
		return Config{}, err
	}

	if cfg.Database.Port <= 0 || cfg.Database.Port > 65535 {
		return Config{}, fmt.Errorf("DB_PORT must be between 1 and 65535")
	}
	if cfg.Database.Name == "" {
		return Config{}, fmt.Errorf("DB_NAME is required")
	}
	if cfg.Database.MaxOpenConns <= 0 {
		return Config{}, fmt.Errorf("DB_MAX_OPEN_CONNS must be greater than 0")
	}
	if cfg.Database.MaxIdleConns < 0 || cfg.Database.MaxIdleConns > cfg.Database.MaxOpenConns {
		return Config{}, fmt.Errorf("DB_MAX_IDLE_CONNS must be between 0 and DB_MAX_OPEN_CONNS")
	}
	if cfg.Database.PartitionMonthsAhead < 0 {
		return Config{}, fmt.Errorf("DB_PARTITION_MONTHS_AHEAD must not be negative")
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = devJWTSecret
		fmt.Fprintln(os.Stderr, "WARNING: JWT_SECRET is not set, using a development-only secret")
	}
	if len(cfg.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("JWT_SECRET must be at least 32 bytes")
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return parsed, nil
}

func envBool(key string, fallback bool) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid %s: %w", key, err)
	}
	return parsed, nil
}

func envInt(key string, fallback int) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return parsed, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}
