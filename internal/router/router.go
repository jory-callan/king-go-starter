package router

import (
	"king-starter/internal/app"
	"king-starter/internal/router/core/access" // 主要的access路由
	permission_access "king-starter/internal/router/core/access/permission"
	role_access "king-starter/internal/router/core/access/role"
	"king-starter/internal/router/core/identity"
	"king-starter/internal/router/core/user"
	"king-starter/internal/router/hello"
)

// RegisterAutoMigrate 统一在这里自动迁移数据库表结构, 按需启用
func RegisterAutoMigrate(app *app.App) {
	app.Db.AutoMigrate(
		// core 模块
		// core/user
		&user.CoreUser{},
		// core/access
		&role_access.CoreRole{},
		&permission_access.CorePermission{}, // 权限表现在包含菜单功能
		&role_access.CoreRoleMenu{},         // 角色菜单关联表 (兼容旧版)
		&role_access.CoreRolePermission{},   // 角色权限关联表
		// core/identity
		&identity.CoreLoginLog{},
		&identity.CoreOAuthClient{},
		&identity.CoreOAuthToken{},
		&identity.CoreTwoFA{},
		&identity.CoreTwoFALog{},
	)
}

func RegisterAll(app *app.App) {
	// 按需启用自动迁移数据库表结构
	RegisterAutoMigrate(app)

	// hello 测试模块
	hello.RegisterRoutes(app)

	// user 模块
	user.RegisterRoutes(app)
	// core 模块
	access.RegisterRoutes(app)
	identity.RegisterRoutes(app)
}
