package router

import (
	"king-starter/internal/app"
	"king-starter/internal/router/core/auth"
	"king-starter/internal/router/core/permission"
	"king-starter/internal/router/core/role"
	"king-starter/internal/router/core/user"
	"king-starter/internal/router/hello"
)

var prefix = "/api/v1"

// RegisterAutoMigrate 统一在这里自动迁移数据库表结构, 按需启用
func RegisterAutoMigrate(app *app.App) {
	user.RegisterAutoMigrate(app)
	role.RegisterAutoMigrate(app)
	permission.RegisterAutoMigrate(app)
	auth.RegisterAutoMigrate(app)
}

func RegisterAll(app *app.App) {
	// 按需启用自动迁移数据库表结构
	RegisterAutoMigrate(app)

	// hello 测试模块
	hello.RegisterRoutes(app, prefix)

	// core 模块
	user.RegisterRoutes(app, prefix)
	role.RegisterRoutes(app, prefix)
	permission.RegisterRoutes(app, prefix)
	// 认证模块
	// identity.RegisterRoutes(app)
	auth.RegisterAuthRoutes(app)
}
