package router

import (
	"king-starter/internal/app"
	"king-starter/internal/router/core/access"
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
		&access.CoreRole{},
		&access.CoreMenu{},
		&access.CorePermission{},
		&access.CoreRoleMenu{},
		&access.CoreRolePermission{},
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
