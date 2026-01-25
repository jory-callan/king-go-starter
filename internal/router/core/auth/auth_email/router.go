package auth_email

import (
	"king-starter/internal/app"
)

// RegisterRoutes 注册邮箱认证路由
func RegisterRoutes(app *app.App) {
	handler := NewEmailHandler()

	e := app.Server.Engine()

	// 邮箱认证路由组
	authGroup := e.Group("/api/core/auth")
	{
		authGroup.POST("/email/send-code", handler.SendVerificationCode) // 发送验证码
		authGroup.POST("/email/verify", handler.VerifyEmail)             // 验证邮箱
	}
}
