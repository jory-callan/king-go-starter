package app

import (
	"king-starter/config"
	"king-starter/pkg/database"
	"king-starter/pkg/http"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logger"
	gohttp "net/http"
)

// 全局唯一的 App 实例
var globalApp *App

// App 应用核心，持有所有共享依赖
// App 只提供能力（HTTP Server / MQ Client / Cron Scheduler / DB / Redis）
type App struct {
	// 配置文件实例
	Config *config.Config
	// 日志实例
	Log *logger.Logger
	// 数据库实例
	Db *database.DB
	// 缓存实例
	//Cache *cache.Cache
	// JWT 实例
	Jwt *jwt.JWT
	// Http 服务实例
	Server *http.Server
}

// New 初始化 App 实例
func New(cfg *config.Config) *App {
	// 初始化 log
	log := Must(logger.New(cfg.Logger))
	log.Info("logger initialized")

	// 初始化 database.default
	databaseConfig := cfg.Database["default"]
	defaultDB := Must(database.New(databaseConfig, log))
	log.Info("database default initialized")

	// 初始化 JWT
	jwtIns := Must(jwt.NewWithConfig(cfg.Jwt))
	log.Info("jwt initialized")

	// 初始化 HTTP 服务
	server := Must(http.New(cfg.Http, log))

	globalApp = &App{
		Log:    log.Named("app"),
		Config: cfg,
		Db:     defaultDB,
		Jwt:    jwtIns,
		Server: server,
	}
	log.Info("globalApp initialized")
	return globalApp
}

// Start 启动
func (c *App) Start() {
	err := c.Server.Start()
	if err != nil && err != gohttp.ErrServerClosed {
		c.Log.Panic("server start failed. Error msg is: %s" + err.Error())
	}
}

// Shutdown 资源清理
func (c *App) Shutdown() {
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

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func MustCore() *App {
	if globalApp == nil {
		panic("g globalApp not initialized, call g.NewWithConfig() first")
	}
	return globalApp
}

// 提供全局访问方法

func DB() *database.DB       { return MustCore().Db }
func Logger() *logger.Logger { return MustCore().Log }
func Config() *config.Config { return MustCore().Config }
func JWT() *jwt.JWT          { return MustCore().Jwt }
func Server() *http.Server   { return MustCore().Server }
