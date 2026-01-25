package permission

// CreatePermissionReq 创建权限请求
type CreatePermissionReq struct {
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	ParentID string `json:"parent_id"` // 父级权限ID，用于菜单层级
	Path     string `json:"path"`      // 路由路径
	Icon     string `json:"icon"`      // 图标
	Sort     int    `json:"sort"`      // 排序
	Status   int    `json:"status"`    // 状态
	Remark   string `json:"remark"`
}

// UpdatePermissionReq 更新权限请求
type UpdatePermissionReq struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	ParentID string `json:"parent_id"` // 父级权限ID，用于菜单层级
	Path     string `json:"path"`      // 路由路径
	Icon     string `json:"icon"`      // 图标
	Sort     int    `json:"sort"`      // 排序
	Status   int    `json:"status"`    // 状态
	Remark   string `json:"remark"`
}

// PermissionListReq 权限列表请求
type PermissionListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Code     string `json:"code" form:"code"`
	Name     string `json:"name" form:"name"`
	Type     string `json:"type" form:"type"`
	ParentID string `json:"parent_id" form:"parent_id"`
	Status   int    `json:"status" form:"status"`
}

// AssignRolePermissionsReq 分配角色权限请求
type AssignRolePermissionsReq struct {
	PermissionIDs []string `json:"permission_ids"`
}
