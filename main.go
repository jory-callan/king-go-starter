package main

import (
	"king-starter/config"
	"king-starter/internal/app"
	"king-starter/pkg/logger"

	"github.com/gookit/goutil/dump"
)

func main() {
	// 加载配置
	cfg := config.Load()
	dump.P(cfg)
	// 初始化日志
	logger := logger.New(cfg.Logger)
	defer logger.Sync()
	// 初始化应用核心
	core := app.New(cfg)
	// 初始化HTTP服务器
	// httpServer := http.NewWithConfig(cfg.Http, logger)
	// // 注册路由
	// router.RegisterRoutes(httpServer.Echo())
	// httpServer.Start()
	// httpServer.WaitForSignal()
	core.Shutdown()
}
