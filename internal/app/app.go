package app

import (
	"king-starter/config"
	"king-starter/pkg/database"
	"king-starter/pkg/logger"
)

// Core 应用核心，持有所有共享依赖
type Core struct {
	// 基础设施
	Logger *logger.Logger
	Config *config.Config
	DB     *database.DB
}

// 1. 定义一个包级私有变量，用来持有全局唯一的 Core 实例
var core *Core

// New 初始化 Core 实例，并赋值给全局变量 core
// 注意：通常建议返回 *Core 而不是无返回，或者只初始化一次。这里为了兼容你的逻辑。
func New(cfg *config.Config) *Core {
	log := logger.New(cfg.Logger)
	log.Info("logger initialized")
	db := database.NewInstanceManager(cfg.Database, log)
	defaultDB := db.Get("default")

	// 2. 实例化并赋值给全局变量
	core = &Core{
		Logger: log,
		Config: cfg,
		DB:     defaultDB,
	}

	return core
}

// Shutdown 资源清理
func (c *Core) Shutdown() {
	if c == nil {
		return
	}
	// 关闭Redis连接
	// c.Redis.Close()
	// 关闭数据库连接
	c.DB.Close()
	// 关闭日志
	c.Logger.Sync()
}

func mustCore() *Core {
	if core == nil {
		panic("g core not initialized, call g.NewWithConfig() first")
	}
	return core
}

// 提供全局访问方法
func DB() *database.DB       { return mustCore().DB }
func Logger() *logger.Logger { return mustCore().Logger }
func Config() *config.Config { return mustCore().Config }

// func Redis() *redis.Client   { return mustCore().Redis }
// func JWT() *jwt.JWT          { return mustCore().JWT }
