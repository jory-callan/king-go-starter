package database

import (
	"king-starter/config"
	"king-starter/pkg/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew
func TestNewDB(t *testing.T) {
	// 使用默认配置创建logger
	log := logger.NewWithDefaultConfig()
	assert.NotNil(t, log)

	// 使用SQLite内存数据库测试
	cfg := config.DatabaseConfig{
		Driver:          "sqlite3",
		DSN:             ":memory:",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * 60 * 1000000000,
		ConnMaxIdleTime: 5 * 60 * 1000000000,
		LogLevel:        "info",
	}

	db := New(&cfg, log)
	assert.NotNil(t, db)
	assert.NotNil(t, db.DB)

	// 测试连接
	sqlDB, err := db.DB.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
	assert.NoError(t, sqlDB.Close())
}

// TestNewManager
func TestNewManager(t *testing.T) {
	defaultCfg := config.DatabaseConfig{
		Driver:          "sqlite3",
		DSN:             ":memory:",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * 60 * 1000000000,
	}
	cfg := map[string]*config.DatabaseConfig{
		"default": &defaultCfg,
	}

	log := logger.NewWithDefaultConfig()
	assert.NotNil(t, log)

	manager := NewInstanceManager(cfg, log)
	assert.NotNil(t, manager)
	assert.NotEmpty(t, manager.instances)
}
