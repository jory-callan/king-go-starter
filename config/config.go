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
	return nil
}

// Summary 返回配置摘要（用于日志）
func (c *Config) Summary() string {
	return fmt.Sprintf("[config] 摘要如下: %v", c)
}
