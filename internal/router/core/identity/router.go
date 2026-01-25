package identity

import "king-starter/internal/app"

func RegisterAutoMigrate(app *app.App) {
	// 自动迁移数据库表结构
	app.Db.DB.AutoMigrate(
		&CoreLoginLog{},
		&CoreRefreshToken{},

		&CoreUserOAuthClient{},
		&CoreUserOAuthCode{},
		&CoreUserOAuthToken{},

		&CoreUserTwoFA{},
		&CoreUserTwoFALog{},
	)
}

// RegisterRoutes 提供 Identity 模块的路由注册方法
func RegisterRoutes(app *app.App) {
	var loginRepo = NewLoginRepo(app.Db.DB)
	var oauthRepo = NewOAuthRepo(app.Db.DB)
	var twoFARepo = NewTwoFARepo(app.Db.DB)

	var loginHandler = NewLoginHandler(loginRepo)
	var oauthHandler = NewOAuthHandler(oauthRepo)
	var twoFAHandler = NewTwoFAHandler(twoFARepo)
	var registerHandler = NewRegisterHandler(app)

	e := app.Server.Engine()

	// 统一认证路由
	authGroup := e.Group("/api/core/auth")
	{
		authGroup.POST("/login", loginHandler.Login) // 密码登录
		authGroup.POST("/logout", loginHandler.Logout)
		authGroup.POST("/refresh", loginHandler.RefreshToken)

		// 手机认证路由
		authGroup.POST("/login/phone", loginHandler.Login) // 手机号登录（需要扩展）

		// OAuth2 认证路由
		authGroup.GET("/oauth/authorize", oauthHandler.Authorize)
		authGroup.POST("/oauth/token", oauthHandler.GetToken)
		authGroup.GET("/oauth/userinfo", oauthHandler.GetUserInfo)

		// 2FA 认证路由
		authGroup.POST("/2fa/verify", twoFAHandler.VerifyTwoFA)   // 2FA验证
		authGroup.POST("/2fa/enable", twoFAHandler.EnableTwoFA)   // 启用2FA
		authGroup.POST("/2fa/disable", twoFAHandler.DisableTwoFA) // 禁用2FA
	}

	// 注册路由
	registerGroup := e.Group("/api/core/register")
	{
		registerGroup.POST("", registerHandler.Register)
	}

	// 重置密码路由
	resetPasswordGroup := e.Group("/api/core/password")
	{
		resetPasswordGroup.PUT("/reset", registerHandler.ResetPassword)
	}
}
