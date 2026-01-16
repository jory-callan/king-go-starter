package router

import (
	"king-starter/internal/app"
	"king-starter/internal/router/hello"
	"king-starter/internal/router/user"
)

// Router 不自己初始化 Core，只注册自己到 Core 提供的运行环境
type Router interface {
	Name() string
	Register(app *app.App)
}

func RegisterAll(app *app.App) {
	// 注册模块
	hello.New().Register(app)

	// 用户路由
	user.New(app).Register(app)
}
