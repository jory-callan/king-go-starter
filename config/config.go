package config

import (
	"fmt"

	"king-starter/pkg/database"
	"king-starter/pkg/http"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logx"
)

type Config struct {
	Logger   *logx.LoggerConfig
	Http     *http.HttpConfig
	Database struct {
		Default *database.DatabaseConfig
	}
	Jwt *jwt.JwtConfig
}

// DefaultConfig 返回默认的日志配置
func DefaultConfig() Config {
	c := Config{}
	defaultLoggerConfig := logx.DefaultLoggerConfig()
	defaultHttpConfig := http.DefaultHttpConfig()
	defaultDatabaseConfig := database.DefaultDatabaseConfig()
	defaultJwtConfig := jwt.DefaultJwtConfig()
	c.Logger = &defaultLoggerConfig
	c.Http = &defaultHttpConfig
	c.Database.Default = &defaultDatabaseConfig
	c.Jwt = &defaultJwtConfig
	return c
}

func (c *Config) Validate() error {
	return nil
}

// Summary 返回配置摘要（用于日志）
func (c *Config) Summary() string {
	return fmt.Sprintf("[config] 摘要如下: %v", c)
}
