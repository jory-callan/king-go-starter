package main

import (
	"king-starter/config"
	"king-starter/internal/app"
	"king-starter/pkg/http"
	"king-starter/pkg/logger"

	"github.com/gookit/goutil/dump"
)

func main() {
	// 加载配置
	cfg := config.Load()
	dump.P(cfg)
	// 初始化日志
	log := logger.New(cfg.Logger)
	defer log.Close()
	// 初始化应用核心
	core := app.New(cfg)
	// 初始化HTTP服务器
	httpServer := http.New(cfg.Http, log)
	// // 注册路由
	// router.RegisterRoutes(httpServer.Echo())
	_ = httpServer.Start()
	// httpServer.WaitForSignal()
	log.Info("服务卸载")
	core.Shutdown()
}
