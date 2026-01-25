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

// RoleListReq 角色列表请求
type RoleListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Code     string `json:"code" form:"code"`
	Name     string `json:"name" form:"name"`
	Status   int    `json:"status" form:"status"`
}

// AssignUserRoleReq 分配用户角色请求
type AssignUserRoleReq struct {
	RoleIDs []string `json:"role_ids" binding:"required"`
}

// GetUserRolesResp 获取用户角色响应
type GetUserRolesResp struct {
	UserID string     `json:"user_id"`
	Roles  []CoreRole `json:"roles"`
}
