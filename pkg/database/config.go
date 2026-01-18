package database

import (
	"fmt"
)

// DatabaseConfig 单个数据库实例配置
type DatabaseConfig struct {
	Driver          string // 驱动类型: mysql, postgresql, sqlite3
	DSN             string // 数据源连接字符串
	MaxOpenConns    int    // 最大打开连接数
	MaxIdleConns    int    // 最大空闲连接数
	ConnMaxLifetime int    // 连接最大存活时间(毫秒)
	ConnMaxIdleTime int    // 连接最大空闲时间(毫秒)
	LogLevel        string // 日志级别: silent, error, warn, info
	SlowThreshold   int    // 慢查询阈值（毫秒）
	TablePrefix     string // 表名前缀
}

// DefaultDatabaseConfig 返回默认的单实例配置
func DefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:          "sqlite3",
		DSN:             "./test.db",
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * 60 * 1000, // 30分钟
		ConnMaxIdleTime: 5 * 60 * 1000,  // 5分钟
		LogLevel:        "info",
		SlowThreshold:   1000, // 1秒
		TablePrefix:     "king_",
	}
}

// Validate 验证实例配置
func (c *DatabaseConfig) Validate() error {
	// 验证驱动类型
	switch c.Driver {
	case "mysql", "postgresql", "postgres", "pg", "pgsql", "sqlite3", "sqlite":
	default:
		return fmt.Errorf("[config] database instance invalid driver: %s (supported: mysql, postgresql, postgres, pg, pgsql, sqlite3, sqlite)", c.Driver)
	}
	// 验证DSN
	if c.DSN == "" {
		return fmt.Errorf("[config] database instance DSN is required")
	}
	// 验证连接数
	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("[config] database instance max_open_conns must be positive")
	}
	if c.MaxIdleConns < 0 {
		return fmt.Errorf("[config] database instance max_idle_conns cannot be negative")
	}
	if c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("[config] database instance max_idle_conns cannot exceed max_open_conns")
	}
	// 验证日志级别
	switch c.LogLevel {
	case "silent", "error", "warn", "info":
	default:
		return fmt.Errorf("[config] database instance invalid log level: %s (supported: silent, error, warn, info)", c.LogLevel)
	}
	return nil
}
