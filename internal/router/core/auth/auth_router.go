package auth

import (
	"king-starter/internal/app"
	"king-starter/internal/router/core/auth/auth_2fa"
	"king-starter/internal/router/core/auth/auth_email"
	"king-starter/internal/router/core/auth/auth_oauth2"
	"king-starter/internal/router/core/auth/auth_password"
)

// RegisterAuthRoutes 注册所有认证相关路由
func RegisterAuthRoutes(app *app.App) {
	// 注册密码认证路由
	auth_password.RegisterRoutes(app)

	// 注册邮箱认证路由
	auth_email.RegisterRoutes(app)

	// 注册 2FA 认证路由
	auth_2fa.RegisterRoutes(app)

	// 注册 OAuth2 认证路由
	auth_oauth2.RegisterRoutes(app)
}
