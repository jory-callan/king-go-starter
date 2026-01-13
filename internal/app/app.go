package app

import (
	"king-starter/config"
	"king-starter/pkg/database"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logger"
)

// Core 应用核心，持有所有共享依赖
type Core struct {
	// 基础设施
	Log    *logger.Logger
	Config *config.Config
	Db     *database.DB
	Jwt    *jwt.JWT
}

// 1. 定义一个包级私有变量，用来持有全局唯一的 Core 实例
var core *Core

// New 初始化 Core 实例，并赋值给全局变量 core
// 注意：通常建议返回 *Core 而不是无返回，或者只初始化一次。这里为了兼容你的逻辑。
func New(cfg *config.Config) *Core {
	log := logger.New(cfg.Logger)
	log.Info("logger initialized")

	databaseConfig := cfg.Database["default"]
	defaultDB, _ := database.New(databaseConfig, log)
	log.Info("database default initialized")

	jwtConfig := cfg.Jwt
	jwtIns := jwt.NewWithConfig(jwtConfig)
	log.Info("jwt initialized")
	// 2. 实例化并赋值给全局变量
	core = &Core{
		Log:    log.With(logger.String("component", "core")),
		Config: cfg,
		Db:     defaultDB,
		Jwt:    jwtIns,
	}
	log.Info("core initialized")

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
	c.Db.Close()
	// 关闭日志
	c.Log.Close()
}

func mustCore() *Core {
	if core == nil {
		panic("g core not initialized, call g.NewWithConfig() first")
	}
	return core
}

// 提供全局访问方法

func DB() *database.DB       { return mustCore().Db }
func Logger() *logger.Logger { return mustCore().Log }
func Config() *config.Config { return mustCore().Config }
func JWT() *jwt.JWT          { return mustCore().Jwt }
