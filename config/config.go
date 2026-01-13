package config

import (
	"fmt"
	"king-starter/pkg/database"
	"king-starter/pkg/http"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logger"
)

type Config struct {
	Logger   *logger.LoggerConfig
	Http     *http.HttpConfig
	Database map[string]*database.DatabaseConfig
	Jwt      *jwt.JwtConfig
}

// DefaultConfig 返回默认的日志配置
func DefaultConfig() Config {
	defaultLoggerConfig := logger.DefaultLoggerConfig()
	defaultHttpConfig := http.DefaultHttpConfig()
	defaultDatabaseConfig := database.DefaultDatabaseConfig()
	return Config{
		Logger:   &defaultLoggerConfig,
		Http:     &defaultHttpConfig,
		Database: map[string]*database.DatabaseConfig{"default": &defaultDatabaseConfig},
	}
}

func (c *Config) Validate() error {
	var errList []error

	// 验证HTTP配置
	if err := c.Http.Validate(); err != nil {
		errList = append(errList, fmt.Errorf("[config] http validation error: %w", err))
	}

	for name, dbConfig := range c.Database {
		if err := dbConfig.Validate(); err != nil {
			errList = append(errList, fmt.Errorf("[config] database %s validation error: %w", name, err))
		}
	}

	if err := c.Logger.Validate(); err != nil {
		errList = append(errList, fmt.Errorf("[config] logger validation error: %w", err))
	}

	if len(errList) > 0 {
		return fmt.Errorf("[config] validation errors: %v", errList)
	}
	return nil
}

// Summary 返回配置摘要（用于日志）
func (c *Config) Summary() string {
	return fmt.Sprintf("[config] 摘要如下: %v", c)
}
