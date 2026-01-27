package permission

import (
	"king-starter/internal/app"
)

// RegisterAutoMigrate 统一在这里自动迁移数据库表结构, 按需启用
func RegisterAutoMigrate(app *app.App) {
	app.Db.AutoMigrate(
		&CorePermission{},
		&CoreRolePermission{},
	)
}

// RegisterRoutes 模块的路由注册方法
func RegisterRoutes(app *app.App, prefix string) {
	var permissionRepo = NewPermissionRepo(app.Db.DB)

	var permHandler = NewPermissionHandler(permissionRepo)

	e := app.Server.Engine()
	// 权限路由（现在包括菜单功能）
	permGroup := e.Group(prefix + "/core/permissions")
	{
		permGroup.POST("", permHandler.CreatePermission)
		permGroup.GET("", permHandler.ListPermissions)
		permGroup.GET("/:id", permHandler.GetPermissionDetail)
		permGroup.PUT("/:id", permHandler.UpdatePermission)
		permGroup.DELETE("/:id", permHandler.DeletePermission)
		permGroup.GET("/tree", permHandler.GetPermissionTree) // 新增权限树接口
	}

	// 角色权限路由
	rolePermGroup := e.Group(prefix + "/core/role-permissions")
	{
		rolePermGroup.PUT("/roles/:role_id/permissions", permHandler.AssignRolePermissions)
		rolePermGroup.GET("/roles/:role_id/permissions", permHandler.GetRolePermissions)
		rolePermGroup.GET("/roles/:role_id/permissions/detail", permHandler.GetRolePermissionsWithDetails)
		rolePermGroup.GET("/roles/:role_id/permissions/tree", permHandler.GetRolePermissionTree)
		rolePermGroup.DELETE("/roles/:role_id/permissions", permHandler.RemoveRolePermissions)
	}

	// 用户权限路由
	userPermGroup := e.Group(prefix + "/core/user-permissions")
	{
		userPermGroup.GET("/users/:user_id/permissions", permHandler.GetUserAllPermissions)
	}
}
