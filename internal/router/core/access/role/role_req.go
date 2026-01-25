package role

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	Name   string `json:"name" binding:"required"`
	Status int    `json:"status"`
	Remark string `json:"remark"`
}

// UpdateRoleMenusReq 更新角色菜单请求
type UpdateRoleMenusReq struct {
	MenuIDs []string `json:"menu_ids"`
}

// UpdateRolePermissionsReq 更新角色权限请求
type UpdateRolePermissionsReq struct {
	PermissionIDs []string `json:"permission_ids"`
}

// RoleListReq 角色列表请求
type RoleListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Code     string `json:"code" form:"code"`
	Name     string `json:"name" form:"name"`
	Status   int    `json:"status" form:"status"`
}
