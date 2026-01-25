package auth_password

import (
	"king-starter/internal/app"
)

// RegisterRoutes 注册密码认证路由
func RegisterRoutes(app *app.App) {
	repo := NewRepository(app.Db.DB)
	handler := NewLoginHandler(repo)

	e := app.Server.Engine()

	// 密码认证路由组
	authGroup := e.Group("/api/core/auth")
	{
		authGroup.POST("/login", handler.Login) // 密码登录
		authGroup.POST("/logout", handler.Logout)
		authGroup.POST("/refresh", handler.RefreshToken)
	}
}
