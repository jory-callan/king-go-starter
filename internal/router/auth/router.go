package auth

import (
	"king-starter/internal/app"
	"king-starter/internal/middleware"
	"king-starter/pkg/logx"
)

// Module 认证模块
type Module struct {
	handler     *AuthHandler
	userHandler *UserHandler
	roleHandler *RoleHandler
	permHandler *PermissionHandler
}

// New 创建认证模块实例
func New(core *app.App) *Module {
	// 初始化服务
	authService := NewAuthService(core.Db.DB, core.Jwt)
	userService := NewUserService(core.Db.DB)
	roleService := NewRoleService(core.Db.DB)
	permService := NewPermissionService(core.Db.DB)

	// 初始化处理器
	handler := NewAuthHandler(authService)
	userHandler := NewUserHandler(userService)
	roleHandler := NewRoleHandler(roleService)
	permHandler := NewPermissionHandler(permService)

	return &Module{
		handler:     handler,
		userHandler: userHandler,
		roleHandler: roleHandler,
		permHandler: permHandler,
	}
}

// Name 返回模块名称
func (m *Module) Name() string {
	return "auth"
}

// Register 注册模块
func (m *Module) Register(core *app.App) {
	// 自动迁移模型
	err := core.Db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
	)
	if err != nil {
		panic("auth module: failed to auto migrate models: " + err.Error())
	}

	logx.Info("auth module: models migrated successfully")

	// 获取Echo实例
	e := core.Server.Engine()

	// 注册中间件
	rateLimitMiddleware := middleware.RateLimit()
	jwtMiddleware := middleware.JWT(core.Jwt)

	// 公共路由组
	public := e.Group("/api/v1/auth", rateLimitMiddleware)
	{
		// 认证相关路由
		public.POST("/login", m.handler.Login)
	}

	// 受保护路由组（需要JWT认证）
	protected := e.Group("/api/v1/auth", rateLimitMiddleware, jwtMiddleware)
	{
		// 用户相关路由
		protected.POST("/change-password", m.handler.ChangePassword)
		protected.GET("/permissions", m.handler.GetUserPermissions)
		protected.GET("/roles", m.handler.GetUserRoles)

		// 用户管理路由
		protected.POST("/users", m.userHandler.CreateUser)
		protected.PUT("/users", m.userHandler.UpdateUser)
		protected.DELETE("/users/:id", m.userHandler.DeleteUser)
		protected.GET("/users/:id", m.userHandler.GetUserByID)
		protected.GET("/users", m.userHandler.ListUsers)

		// 角色管理路由
		protected.POST("/roles", m.roleHandler.CreateRole)
		protected.PUT("/roles", m.roleHandler.UpdateRole)
		protected.DELETE("/roles/:id", m.roleHandler.DeleteRole)
		protected.GET("/roles/:id", m.roleHandler.GetRoleByID)
		protected.GET("/roles", m.roleHandler.ListRoles)

		// 权限管理路由
		protected.POST("/permissions", m.permHandler.CreatePermission)
		protected.PUT("/permissions", m.permHandler.UpdatePermission)
		protected.DELETE("/permissions/:id", m.permHandler.DeletePermission)
		protected.GET("/permissions/:id", m.permHandler.GetPermissionByID)
		protected.GET("/permissions", m.permHandler.ListPermissions)
	}

	logx.Info("auth module: routes registered successfully")
}
