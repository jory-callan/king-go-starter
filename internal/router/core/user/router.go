package user

import (
	"king-starter/internal/app"
	"king-starter/pkg/logx"
)

// RegisterAutoMigrate 统一在这里自动迁移数据库表结构, 按需启用
func RegisterAutoMigrate(app *app.App) {
	app.Db.AutoMigrate(
		&CoreUser{},
	)
}

func RegisterRoutes(app *app.App) {
	var repo = NewRepository(app.Db.DB)
	var handler = NewHandler(repo)

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
