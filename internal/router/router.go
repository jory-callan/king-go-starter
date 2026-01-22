package router

import (
	"king-starter/internal/app"
	"king-starter/internal/router/core/access"
	"king-starter/internal/router/core/identity"
	"king-starter/internal/router/hello"
)

func RegisterAll(app *app.App) {
	// hello 模块
	hello.Register(app)
	
	// core 模块
	access.RegisterRoutes(app)
	identity.RegisterRoutes(app)
}
