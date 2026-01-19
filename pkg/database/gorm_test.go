package database

import (
	"king-starter/pkg/logx"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew
func TestNewDB(t *testing.T) {
	logxCfg := logx.DefaultLoggerConfig()
	logx.NewSlog(&logxCfg)

	// 使用SQLite内存数据库测试
	cfg := DatabaseConfig{
		Driver:          "sqlite3",
		DSN:             ":memory:",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * 60 * 1000000000,
		ConnMaxIdleTime: 5 * 60 * 1000000000,
		LogLevel:        "info",
	}

	db, err := New(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, db.DB)

	// 测试连接
	sqlDB, err := db.DB.DB()
	assert.NoError(t, err)
	assert.NoError(t, sqlDB.Ping())
	assert.NoError(t, sqlDB.Close())
}
