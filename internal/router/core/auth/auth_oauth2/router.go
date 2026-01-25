package auth_oauth2

import (
	"king-starter/internal/app"
)

// RegisterRoutes 注册 OAuth2 认证路由
func RegisterRoutes(app *app.App) {
	repo := NewRepository(app.Db.DB)
	handler := NewOAuthHandler(repo)

	e := app.Server.Engine()

	// OAuth2 认证路由组
	authGroup := e.Group("/api/core/auth")
	{
		authGroup.GET("/oauth/authorize", handler.Authorize)
		authGroup.POST("/oauth/token", handler.GetToken)
		authGroup.GET("/oauth/userinfo", handler.GetUserInfo)
	}
}
