package user

import (
	"king-starter/internal/app"
	"king-starter/internal/middleware"
)

type Module struct {
	handler *Handler
}

func New(app *app.App) *Module {
	db := app.Db.DB
	jwt := app.Jwt
	jwtExpire := app.Config.Jwt.Expire
	return &Module{
		handler: NewHandler(NewService(db, jwt, jwtExpire)),
	}
}

func (m *Module) Name() string {
	return "user"
}

func (m *Module) Register(core *app.App) {
	// AutoMigrate
	err := core.Db.Debug().AutoMigrate(
		&User{},
	)
	if err != nil {
		panic(err)
	}

	e := core.Server.Engine()
	jwtMiddleware := middleware.JWT(core.Jwt)
	rateLimitMiddleware := middleware.RateLimit(core.Log)

	public := e.Group("/api/v1/user", rateLimitMiddleware)
	{
		public.POST("/register", m.handler.Register)
		public.POST("/login", m.handler.Login)
		public.POST("/reset-code", m.handler.GenerateResetCode)
		public.POST("/reset-password", m.handler.ResetPassword)
	}

	protected := e.Group("/api/v1/user", rateLimitMiddleware, jwtMiddleware)
	{
		protected.GET("/profile", m.handler.GetProfile)
	}
}
