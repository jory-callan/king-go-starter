package access

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

	role := &Role{
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

	return response.SuccessWithMsg(c, "创建成功", nil)
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

	return response.Success(c, map[string]interface{}{
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

	return response.SuccessWithMsg(c, "更新成功", nil)
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

	return response.SuccessWithMsg(c, "删除成功", nil)
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

	return response.SuccessWithMsg(c, "菜单分配成功", nil)
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

	return response.SuccessWithMsg(c, "权限分配成功", nil)
}

// ListRoles 获取角色列表
func (h *RoleHandler) ListRoles(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	code := c.QueryParam("code")
	name := c.QueryParam("name")
	statusStr := c.QueryParam("status")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	db := h.roleRepo.GetDB(c.Request().Context())
	query := db.Model(&Role{})

	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	var roles []*Role
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	return response.Success(c, map[string]interface{}{
		"list":  roles,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
