package database

import (
	"context"
	"fmt"
	"time"

	"king-starter/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DB 封装 gorm.DB，提供统一数据库接口
type DB struct {
	*gorm.DB
	log *logger.Logger
}

// New 创建 Gorm 实例，使用提供的配置
func New(cfg *DatabaseConfig, log *logger.Logger) (*DB, error) {
	dbLog := log.Named("database")
	// 根据驱动类型创建dialector
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "postgresql", "postgres", "pg", "pgsql":
		dialector = postgres.Open(cfg.DSN)
	case "sqlite3", "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	default:
		log.Error(fmt.Sprintf("[database] unsupported driver: %s", cfg.Driver))
		return nil, fmt.Errorf("[database] instance invalid driver: %s (supported: mysql, postgresql, postgres, pg, pgsql, sqlite3, sqlite)", cfg.Driver)
	}

	// 配置命名策略
	ns := schema.NamingStrategy{
		SingularTable: true,            // 单数表名
		TablePrefix:   cfg.TablePrefix, // 表名前缀
	}
	//创建GORM配置
	gormConfig := &gorm.Config{
		Logger: newGormLogger(dbLog, cfg.LogLevel, time.Duration(cfg.SlowThreshold)*time.Millisecond),
		//Logger:                                   glog.Default,
		SkipDefaultTransaction:                   true, // 跳过默认每条的事务，提高性能
		PrepareStmt:                              true, // 预编译语句，提高性能
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束，避免迁移时的循环依赖问题
		NamingStrategy:                           ns,   // 配置命名策略
	}

	// 创建Gorm实例, 连接数据库
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Error(fmt.Sprintf("[database] failed to open database: %v", err))
		return nil, err
	}

	// 获取底层sql.DB并配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(fmt.Sprintf("[database] failed to get sql.DB: %v", err))
		return nil, err
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
		DB:  db,
		log: dbLog,
	}, nil
}

// HealthCheck 健康检查
func (d *DB) HealthCheck(ctx context.Context) error {
	start := time.Now()
	sqlDB, err := d.DB.DB()
	if err != nil {
		d.log.Error("get sql db failed", logger.Error(err))
		return err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		d.log.Error("[database] ping failed", logger.Error(err))
		return err
	}
	d.log.Debug("[database] ping ok",
		logger.Duration("cost", time.Since(start)),
	)
	return nil
}

// Close 关闭数据库连接
func (d *DB) Close() {
	if d.DB == nil {
		return
	}
	var err error
	sqlDB, err := d.DB.DB()
	if err != nil {
		d.log.Error("[database] get sql db failed")
	}
	err = sqlDB.Close()
	if err != nil {
		d.log.Error("[database] close failed")
	}
	d.log.Info("[database] closed")
}
