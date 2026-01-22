package access

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoleHandler struct {
	roleRepo       *RoleRepo
	menuRepo       *MenuRepo
	permissionRepo *PermissionRepo
}

func NewRoleHandler(roleRepo *RoleRepo, menuRepo *MenuRepo, permissionRepo *PermissionRepo) *RoleHandler {
	return &RoleHandler{
		roleRepo:       roleRepo,
		menuRepo:       menuRepo,
		permissionRepo: permissionRepo,
	}
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c echo.Context) error {
	var req CreateRoleReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从 Context 获取

	role := &CoreRole{
		Code:      req.Code,
		Name:      req.Name,
		Status:    req.Status,
		Remark:    req.Remark,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if err := h.roleRepo.Create(c.Request().Context(), role); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建失败")
	}

	return response.SuccessWithMsg[any](c, "创建成功", nil)
}

// GetRoleDetail 获取角色详情
func (h *RoleHandler) GetRoleDetail(c echo.Context) error {
	id := c.Param("id")
	role, err := h.roleRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "角色不存在")
	}

	menus, err := h.roleRepo.GetRoleMenus(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取角色菜单失败")
	}

	perms, err := h.roleRepo.GetRolePermissions(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取角色权限失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"role":           role,
		"menu_ids":       menus,
		"permission_ids": perms,
	})
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c echo.Context) error {
	id := c.Param("id")
	var req UpdateRoleReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从 Context 获取

	role, err := h.roleRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "角色不存在")
	}

	role.Name = req.Name
	role.Status = req.Status
	role.Remark = req.Remark
	role.UpdatedBy = operatorID

	if err := h.roleRepo.Update(c.Request().Context(), role); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	return response.SuccessWithMsg[any](c, "更新成功", nil)
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")
	operatorID := "system-admin" // TODO: 从 Context 获取

	role, err := h.roleRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "角色不存在")
	}

	// 先更新删除人
	role.DeletedBy = operatorID
	if err := h.roleRepo.Update(c.Request().Context(), role); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	// 再删除角色
	if err := h.roleRepo.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "删除失败")
	}

	return response.SuccessWithMsg[any](c, "删除成功", nil)
}

// UpdateRoleMenus 更新角色菜单
func (h *RoleHandler) UpdateRoleMenus(c echo.Context) error {
	id := c.Param("id")
	var req UpdateRoleMenusReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	if err := h.roleRepo.AssignMenus(c.Request().Context(), id, req.MenuIDs); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	return response.SuccessWithMsg[any](c, "菜单分配成功", nil)
}

// UpdateRolePermissions 更新角色权限
func (h *RoleHandler) UpdateRolePermissions(c echo.Context) error {
	id := c.Param("id")
	var req UpdateRolePermissionsReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	if err := h.roleRepo.AssignPermissions(c.Request().Context(), id, req.PermissionIDs); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	return response.SuccessWithMsg[any](c, "权限分配成功", nil)
}

// ListRoles 获取角色列表
func (h *RoleHandler) ListRoles(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 确保 NeedCount 为 true 以返回总数
	pq.NeedCount = true

	// 获取查询参数
	code := c.QueryParam("code")
	name := c.QueryParam("name")
	statusStr := c.QueryParam("status")

	// 创建筛选条件的 scope 函数
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if code != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("code LIKE ?", "%"+code+"%")
		})
	}
	if name != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("name LIKE ?", "%"+name+"%")
		})
	}
	if statusStr != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			status, _ := strconv.Atoi(statusStr)
			return db.Where("status = ?", status)
		})
	}

	// 使用 BaseRepo 的分页方法
	result, err := h.roleRepo.PaginationWithScopes(c.Request().Context(), &pq, scopes...)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	return response.SuccessPage[CoreRole](c, *result)
}
