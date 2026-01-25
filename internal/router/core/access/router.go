package access

import (
	"king-starter/internal/app"
	perm "king-starter/internal/router/core/access/permission"
	role_mod "king-starter/internal/router/core/access/role"
)

// RegisterRoutes 提供 Access 模块的路由注册方法
func RegisterRoutes(app *app.App) {
	var roleRepo = role_mod.NewRoleRepo(app.Db.DB)
	var permissionRepo = perm.NewPermissionRepo(app.Db.DB)

	var roleHandler = role_mod.NewRoleHandler(roleRepo, permissionRepo)
	var permHandler = perm.NewPermissionHandler(permissionRepo)

	e := app.Server.Engine()

	// 角色路由
	roleGroup := e.Group("/api/core/roles")
	{
		roleGroup.POST("", roleHandler.CreateRole)
		roleGroup.GET("", roleHandler.ListRoles)
		roleGroup.GET("/:id", roleHandler.GetRoleDetail)
		roleGroup.PUT("/:id", roleHandler.UpdateRole)
		roleGroup.DELETE("/:id", roleHandler.DeleteRole)
		roleGroup.PUT("/:id/permissions", roleHandler.UpdateRolePermissions) // 移除了UpdateRoleMenus
	}

	// 权限路由（现在包括菜单功能）
	permGroup := e.Group("/api/core/permissions")
	{
		permGroup.POST("", permHandler.CreatePermission)
		permGroup.GET("", permHandler.ListPermissions)
		permGroup.GET("/:id", permHandler.GetPermissionDetail)
		permGroup.PUT("/:id", permHandler.UpdatePermission)
		permGroup.DELETE("/:id", permHandler.DeletePermission)
		permGroup.GET("/tree", permHandler.GetPermissionTree) // 新增权限树接口
	}
}
