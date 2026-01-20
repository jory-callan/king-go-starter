package auth

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService *RoleService
}

// NewRoleHandler 创建角色处理器实例
func NewRoleHandler(roleService *RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=50"`
	Code   string `json:"code" validate:"required,min=2,max=50"`
	Status int    `json:"status" validate:"omitempty,oneof=0 1"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	ID     string `json:"id" validate:"required"`
	Name   string `json:"name" validate:"omitempty,min=2,max=50"`
	Code   string `json:"code" validate:"omitempty,min=2,max=50"`
	Status int    `json:"status" validate:"omitempty,oneof=0 1"`
}

// ListRolesRequest 查询角色列表请求
type ListRolesRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Name     string `query:"name" validate:"omitempty"`
	Code     string `query:"code" validate:"omitempty"`
	Status   int    `query:"status" validate:"omitempty,oneof=0 1"`
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c echo.Context) error {
	var req CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	role := &Role{
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
	}

	if role.Status == 0 {
		role.Status = 1 // 默认启用
	}

	err := h.roleService.Create(c.Request().Context(), role, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建角色失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "创建角色成功", role)
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c echo.Context) error {
	var req UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	// 获取现有角色信息
	existingRole, err := h.roleService.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "角色不存在")
	}

	// 更新角色信息
	if req.Name != "" {
		existingRole.Name = req.Name
	}
	if req.Code != "" {
		existingRole.Code = req.Code
	}
	if req.Status != 0 {
		existingRole.Status = req.Status
	}

	err = h.roleService.Update(c.Request().Context(), existingRole, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新角色失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "更新角色成功", existingRole)
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "角色ID不能为空")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	err := h.roleService.Delete(c.Request().Context(), id, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "删除角色失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "删除角色成功", nil)
}

// GetRoleByID 根据ID获取角色
func (h *RoleHandler) GetRoleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "角色ID不能为空")
	}

	role, err := h.roleService.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "角色不存在")
	}

	return response.Success(c, role)
}

// ListRoles 查询角色列表
func (h *RoleHandler) ListRoles(c echo.Context) error {
	// 获取查询参数
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	name := c.QueryParam("name")
	code := c.QueryParam("code")
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
	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		conds["status"] = status
	}

	roles, total, err := h.roleService.ListByCondition(c.Request().Context(), conds, page, pageSize)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询角色列表失败: "+err.Error())
	}

	return response.Success(c, map[string]interface{}{
		"list":      roles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
