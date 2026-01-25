package router

import (
	"king-starter/internal/app"
	"king-starter/internal/router/core/auth"
	"king-starter/internal/router/core/permission"
	"king-starter/internal/router/core/role"
	"king-starter/internal/router/core/user"
	"king-starter/internal/router/hello"
)

// RegisterAutoMigrate 统一在这里自动迁移数据库表结构, 按需启用
func RegisterAutoMigrate(app *app.App) {
	user.RegisterAutoMigrate(app)
	role.RegisterAutoMigrate(app)
	permission.RegisterAutoMigrate(app)
}

func RegisterAll(app *app.App) {
	// 按需启用自动迁移数据库表结构
	RegisterAutoMigrate(app)

	// hello 测试模块
	hello.RegisterRoutes(app)

	// core 模块
	user.RegisterRoutes(app)
	role.RegisterRoutes(app)
	permission.RegisterRoutes(app)
	// 认证模块
	// identity.RegisterRoutes(app)
	auth.RegisterAuthRoutes(app)
}
