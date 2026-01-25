package role

import (
	"king-starter/internal/app"
)

// RegisterAutoMigrate 统一在这里自动迁移数据库表结构, 按需启用
func RegisterAutoMigrate(app *app.App) {
	app.Db.AutoMigrate(
		&CoreRole{},
		&CoreUserRole{},
	)
}

// RegisterRoutes 提供 Access 模块的路由注册方法
func RegisterRoutes(app *app.App) {
	var roleRepo = NewRoleRepo(app.Db.DB)
	var roleHandler = NewRoleHandler(roleRepo)

	e := app.Server.Engine()

	// 角色路由
	roleGroup := e.Group("/api/core/roles")
	{
		roleGroup.POST("", roleHandler.CreateRole)
		roleGroup.GET("", roleHandler.ListRoles)
		roleGroup.GET("/:id", roleHandler.GetRoleDetail)
		roleGroup.PUT("/:id", roleHandler.UpdateRole)
		roleGroup.DELETE("/:id", roleHandler.DeleteRole)
	}

	// 用户角色绑定路由
	userRoleGroup := e.Group("/api/core/user-roles")
	{
		userRoleGroup.PUT("/users/:user_id/roles", roleHandler.AssignRolesToUser)
		userRoleGroup.GET("/users/:user_id/roles", roleHandler.GetUserRoles)
		userRoleGroup.GET("/roles/:role_id/users", roleHandler.GetRoleUsers)
		userRoleGroup.DELETE("/users/:user_id/roles/:role_id", roleHandler.RemoveUserRole)
	}

}
