package identity

import "king-starter/internal/app"

// RegisterRoutes 提供 Identity 模块的路由注册方法
func RegisterRoutes(app *app.App) {
	var loginRepo = NewLoginRepo(app.Db.DB)
	var oauthRepo = NewOAuthRepo(app.Db.DB)
	var twoFARepo = NewTwoFARepo(app.Db.DB)

	var loginHandler = NewLoginHandler(loginRepo)
	var oauthHandler = NewOAuthHandler(oauthRepo)
	var twoFAHandler = NewTwoFAHandler(twoFARepo)

	app.Db.AutoMigrate(
		&CoreLoginLog{},
		&CoreOAuthClient{},
		&CoreOAuthToken{},
		&CoreTwoFA{},
		&CoreTwoFALog{},
	)

	e := app.Server.Engine()

	// 登录路由
	loginGroup := e.Group("/api/core/auth")
	{
		loginGroup.POST("/login", loginHandler.Login)
		loginGroup.POST("/logout", loginHandler.Logout)
		loginGroup.POST("/refresh", loginHandler.RefreshToken)
	}

	// OAuth 路由
	oauthGroup := e.Group("/api/core/oauth")
	{
		oauthGroup.GET("/authorize", oauthHandler.Authorize)
		oauthGroup.POST("/token", oauthHandler.GetToken)
		oauthGroup.GET("/userinfo", oauthHandler.GetUserInfo)
	}

	// 2FA 路由
	twoFAGroup := e.Group("/api/core/2fa")
	{
		twoFAGroup.POST("/enable", twoFAHandler.EnableTwoFA)
		twoFAGroup.POST("/verify", twoFAHandler.VerifyTwoFA)
		twoFAGroup.POST("/disable", twoFAHandler.DisableTwoFA)
	}
}
