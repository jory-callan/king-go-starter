package user

import (
	"king-starter/internal/app"
	"king-starter/pkg/logx"
)

func RegisterRoutes(app *app.App) {
	var handler = NewHandler(app)

	e := app.Server.Engine()
	group := e.Group("/api/v1/core/users")
	{
		group.POST("", handler.Create)
		group.GET("", handler.List)
		group.GET("/:id", handler.GetByID)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}

	logx.Info("Registered user router")
}
