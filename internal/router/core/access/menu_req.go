package access

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	ParentID string `json:"parent_id"`
	Name     string `json:"name" binding:"required"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	ParentID string `json:"parent_id"`
	Name     string `json:"name" binding:"required"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort"`
	Status   int    `json:"status"`
}

// MenuListReq 菜单列表请求
type MenuListReq struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	ParentID string `json:"parent_id" form:"parent_id"`
	Name     string `json:"name" form:"name"`
	Status   int    `json:"status" form:"status"`
}
