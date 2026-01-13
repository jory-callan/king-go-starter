package database

import (
	"context"
	"fmt"
	"king-starter/config"
	"king-starter/pkg/logger"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Gorm 封装 gorm.DB，提供统一数据库接口
type DB struct {
	*gorm.DB
	logger *logger.Logger
	once   sync.Once
}

// New 创建 Gorm 实例，使用提供的配置
func New(cfg *config.DatabaseConfig, log *logger.Logger) *DB {
	// 根据驱动类型创建dialector
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "postgresql", "postgres":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite3", "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	default:
		panic(fmt.Sprintf("[database] unsupported driver: %s", cfg.Driver))
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		Logger:                 newGormLogger(log, cfg.LogLevel, time.Duration(cfg.SlowThreshold)*time.Millisecond),
		SkipDefaultTransaction: true,
		PrepareStmt:            true, // 预编译语句，提高性能
	}

	// 创建Gorm实例, 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		panic(fmt.Sprintf("[database] failed to open database: %v", err))
	}

	// 获取底层sql.DB并配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("[database] failed to get sql.DB: %v", err))
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Millisecond)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTime) * time.Millisecond)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("[database] failed to ping database: %v", err))
	}

	return &DB{
		DB:     db,
		logger: log,
	}
}

// Ping 健康检查
func (d *DB) Ping(ctx context.Context) error {
	start := time.Now()

	sqlDB, err := d.DB.DB()
	if err != nil {
		d.logger.Error("get sql db failed", zap.Error(err))
		return fmt.Errorf("get underlying db failed: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		d.logger.Error("database ping failed", zap.Error(err))
		return fmt.Errorf("ping failed: %w", err)
	}

	d.logger.Debug("database ping ok",
		zap.Duration("cost", time.Since(start)),
	)

	return nil
}

// NewWithDefaultConfig 使用默认配置创建 Gorm 实例
func NewWithDefaultConfig(log *logger.Logger) *DB {
	cfg := config.DefaultDatabaseConfig()

	return New(&cfg, log)
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	if d.DB == nil {
		return nil
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	d.logger.Info("database closing")

	return sqlDB.Close()
}
