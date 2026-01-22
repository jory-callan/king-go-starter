package access

// CreatePermissionReq 创建权限请求
type CreatePermissionReq struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Remark string `json:"remark"`
}

// UpdatePermissionReq 更新权限请求
type UpdatePermissionReq struct {
	Name   string `json:"name" binding:"required"`
	Type   string `json:"type" binding:"required"`
	Remark string `json:"remark"`
}

// PermissionListReq 权限列表请求
type PermissionListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	Code     string `json:"code" form:"code"`
	Name     string `json:"name" form:"name"`
	Type     string `json:"type" form:"type"`
}
