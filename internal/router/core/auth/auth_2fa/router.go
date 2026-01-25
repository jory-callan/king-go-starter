package auth_2fa

import (
	"king-starter/internal/app"
)

// RegisterRoutes 注册 2FA 认证路由
func RegisterRoutes(app *app.App) {
	repo := NewRepository(app.Db.DB)
	handler := NewTwoFAHandler(repo)

	e := app.Server.Engine()

	// 2FA 认证路由组
	authGroup := e.Group("/api/core/auth")
	{
		authGroup.POST("/2fa/verify", handler.VerifyTwoFA)   // 2FA验证
		authGroup.POST("/2fa/enable", handler.EnableTwoFA)   // 启用2FA
		authGroup.POST("/2fa/disable", handler.DisableTwoFA) // 禁用2FA
	}
}