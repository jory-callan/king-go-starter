package app

import (
	"errors"
	gohttp "net/http"
	"time"

	"king-starter/config"
	"king-starter/pkg/database"
	"king-starter/pkg/http"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logx"
)

// 全局唯一的 App 实例
var globalApp *App

// App 应用核心，持有所有共享依赖
// App 只提供能力（HTTP Server / MQ Client / Cron Scheduler / DB / Redis）
type App struct {
	// 配置文件实例
	Config *config.Config
	// 数据库实例
	Db *database.DB
	// JWT 实例
	Jwt *jwt.JWT
	// Http 服务实例
	Server *http.Server
}

// New 初始化 App 实例
func New(cfg *config.Config) *App {
	// 初始化 log
	Must(logx.NewZap(cfg.Logger))
	logx.Info("logger initialized")

	// 初始化 database.default
	databaseConfig := cfg.Database.Default
	defaultDB := Must(database.New(databaseConfig))
	logx.Info("database default initialized")

	// 初始化 JWT
	jwtIns := Must(jwt.NewWithConfig(cfg.Jwt))
	logx.Info("jwt initialized")

	// 初始化 HTTP 服务
	server := Must(http.New(cfg.Http))

	globalApp = &App{
		Config: cfg,
		Db:     defaultDB,
		Jwt:    jwtIns,
		Server: server,
	}
	logx.Info("globalApp initialized")
	return globalApp
}

// Start 启动
func (c *App) Start() {
	err := c.Server.Start()
	if err != nil && !errors.Is(err, gohttp.ErrServerClosed) {
		msg := "server start failed. Error msg is: %s" + err.Error()
		logx.Info(msg)
		panic(err)
	}
}

// Shutdown 资源清理
func (c *App) Shutdown() {
	if c == nil {
		return
	}
	time.Sleep(3 * time.Second)
	// 关闭数据库连接
	c.Db.Close()
	// 关闭日志
	logx.Close()
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func MustCore() *App {
	if globalApp == nil {
		panic("app not initialized, init it first")
	}
	return globalApp
}

// 提供全局访问方法

func DB() *database.DB       { return MustCore().Db }
func Config() *config.Config { return MustCore().Config }
func JWT() *jwt.JWT          { return MustCore().Jwt }
func Server() *http.Server   { return MustCore().Server }
