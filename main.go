package main

import (
	"king-starter/config"
	"king-starter/internal/app"
	"king-starter/internal/router"

	"github.com/gookit/goutil/dump"
)

func main() {
	// 加载配置
	cfg := config.Load()
	dump.P(cfg)
	// 初始化应用核心
	core := app.New(cfg)
	router.RegisterAll(core)
	core.Start()
	core.Shutdown()
}
