package auth

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permService *PermissionService
}

// NewPermissionHandler 创建权限处理器实例
func NewPermissionHandler(permService *PermissionService) *PermissionHandler {
	return &PermissionHandler{
		permService: permService,
	}
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Code     string `json:"code" validate:"required,min=2,max=100"`
	Resource string `json:"resource" validate:"required,min=2,max=255"`
	Method   string `json:"method" validate:"required,min=3,max=10"`
	Status   int    `json:"status" validate:"omitempty,oneof=0 1"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"omitempty,min=2,max=50"`
	Code     string `json:"code" validate:"omitempty,min=2,max=100"`
	Resource string `json:"resource" validate:"omitempty,min=2,max=255"`
	Method   string `json:"method" validate:"omitempty,min=3,max=10"`
	Status   int    `json:"status" validate:"omitempty,oneof=0 1"`
}

// ListPermissionsRequest 查询权限列表请求
type ListPermissionsRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Name     string `query:"name" validate:"omitempty"`
	Code     string `query:"code" validate:"omitempty"`
	Resource string `query:"resource" validate:"omitempty"`
	Method   string `query:"method" validate:"omitempty"`
	Status   int    `query:"status" validate:"omitempty,oneof=0 1"`
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c echo.Context) error {
	var req CreatePermissionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	permission := &Permission{
		Name:     req.Name,
		Code:     req.Code,
		Resource: req.Resource,
		Method:   req.Method,
		Status:   req.Status,
	}

	if permission.Status == 0 {
		permission.Status = 1 // 默认启用
	}

	err := h.permService.Create(c.Request().Context(), permission, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建权限失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "创建权限成功", permission)
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c echo.Context) error {
	var req UpdatePermissionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	// 获取现有权限信息
	existingPerm, err := h.permService.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "权限不存在")
	}

	// 更新权限信息
	if req.Name != "" {
		existingPerm.Name = req.Name
	}
	if req.Code != "" {
		existingPerm.Code = req.Code
	}
	if req.Resource != "" {
		existingPerm.Resource = req.Resource
	}
	if req.Method != "" {
		existingPerm.Method = req.Method
	}
	if req.Status != 0 {
		existingPerm.Status = req.Status
	}

	err = h.permService.Update(c.Request().Context(), existingPerm, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新权限失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "更新权限成功", existingPerm)
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "权限ID不能为空")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	err := h.permService.Delete(c.Request().Context(), id, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "删除权限失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "删除权限成功", nil)
}

// GetPermissionByID 根据ID获取权限
func (h *PermissionHandler) GetPermissionByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "权限ID不能为空")
	}

	permission, err := h.permService.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "权限不存在")
	}

	return response.Success(c, permission)
}

// ListPermissions 查询权限列表
func (h *PermissionHandler) ListPermissions(c echo.Context) error {
	// 获取查询参数
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	name := c.QueryParam("name")
	code := c.QueryParam("code")
	resource := c.QueryParam("resource")
	method := c.QueryParam("method")
	statusStr := c.QueryParam("status")

	// 设置默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询条件
	conds := make(map[string]interface{})
	if name != "" {
		conds["name"] = name
	}
	if code != "" {
		conds["code"] = code
	}
	if resource != "" {
		conds["resource"] = resource
	}
	if method != "" {
		conds["method"] = method
	}
	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		conds["status"] = status
	}

	permissions, total, err := h.permService.ListByCondition(c.Request().Context(), conds, page, pageSize)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询权限列表失败: "+err.Error())
	}

	return response.Success(c, map[string]interface{}{
		"list":      permissions,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
