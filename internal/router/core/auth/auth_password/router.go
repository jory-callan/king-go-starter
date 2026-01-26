package auth_password

import (
	"king-starter/internal/app"
	"king-starter/internal/router/core/user"
)

func RegisterAutoMigrate(app *app.App) {
	app.Db.AutoMigrate(
		&CoreLoginLog{},
		&CoreRefreshToken{},
	)
}

// RegisterRoutes 注册密码认证路由
func RegisterRoutes(app *app.App) {
	repo := NewRepository(app.Db.DB)
	handler := NewLoginHandler(repo, user.NewRepository(app.Db.DB))

	e := app.Server.Engine()

	// 密码认证路由组
	authGroup := e.Group("/api/core/auth")
	{
		authGroup.POST("/register", handler.Register) // 注册
		authGroup.POST("/login", handler.Login)       // 密码登录
		authGroup.POST("/logout", handler.Logout)
		authGroup.POST("/refresh", handler.RefreshToken)
	}
}
