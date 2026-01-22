package access

import "king-starter/internal/app"

// RegisterRoutes 提供 Access 模块的路由注册方法
func RegisterRoutes(app *app.App) {
	var roleRepo = NewRoleRepo(app.Db.DB)
	var menuRepo = NewMenuRepo(app.Db.DB)
	var permissionRepo = NewPermissionRepo(app.Db.DB)

	var roleHandler = NewRoleHandler(roleRepo, menuRepo, permissionRepo)
	var menuHandler = NewMenuHandler(menuRepo)
	var permHandler = NewPermissionHandler(permissionRepo)

	app.Db.AutoMigrate(
		&CoreRole{},
		&CoreMenu{},
		&CorePermission{},
		&CoreRoleMenu{},
		&CoreRolePermission{},
	)

	e := app.Server.Engine()

	// 角色路由
	roleGroup := e.Group("/api/core/roles")
	{
		roleGroup.POST("", roleHandler.CreateRole)
		roleGroup.GET("", roleHandler.ListRoles)
		roleGroup.GET("/:id", roleHandler.GetRoleDetail)
		roleGroup.PUT("/:id", roleHandler.UpdateRole)
		roleGroup.DELETE("/:id", roleHandler.DeleteRole)
		roleGroup.PUT("/:id/menus", roleHandler.UpdateRoleMenus)
		roleGroup.PUT("/:id/permissions", roleHandler.UpdateRolePermissions)
	}

	// 菜单路由
	menuGroup := e.Group("/api/core/menus")
	{
		menuGroup.POST("", menuHandler.CreateMenu)
		menuGroup.GET("", menuHandler.ListMenus)
		menuGroup.GET("/:id", menuHandler.GetMenuDetail)
		menuGroup.PUT("/:id", menuHandler.UpdateMenu)
		menuGroup.DELETE("/:id", menuHandler.DeleteMenu)
	}

	// 权限路由
	permGroup := e.Group("/api/core/permissions")
	{
		permGroup.POST("", permHandler.CreatePermission)
		permGroup.GET("", permHandler.ListPermissions)
		permGroup.GET("/:id", permHandler.GetPermissionDetail)
		permGroup.PUT("/:id", permHandler.UpdatePermission)
		permGroup.DELETE("/:id", permHandler.DeletePermission)
	}
}
