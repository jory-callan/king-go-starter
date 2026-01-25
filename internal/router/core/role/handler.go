package role

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoleHandler struct {
	roleRepo *RoleRepo
}

func NewRoleHandler(roleRepo *RoleRepo) *RoleHandler {
	return &RoleHandler{
		roleRepo: roleRepo,
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

	return response.Success[any](c, map[string]interface{}{
		"role": role,
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

// AssignRolesToUser 为用户分配角色
func (h *RoleHandler) AssignRolesToUser(c echo.Context) error {
	userID := c.Param("user_id")
	var req AssignUserRoleReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从 Context 获取

	if err := h.roleRepo.AssignRolesToUser(c.Request().Context(), userID, req.RoleIDs, operatorID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "分配角色失败")
	}

	return response.SuccessWithMsg[any](c, "角色分配成功", nil)
}

// GetUserRoles 获取用户的角色
func (h *RoleHandler) GetUserRoles(c echo.Context) error {
	userID := c.Param("user_id")
	roles, err := h.roleRepo.GetUserRolesWithDetails(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取用户角色失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"user_id": userID,
		"roles":   roles,
	})
}

// GetRoleUsers 获取角色下的用户列表
func (h *RoleHandler) GetRoleUsers(c echo.Context) error {
	roleID := c.Param("role_id")
	userIDs, err := h.roleRepo.GetRoleUsers(c.Request().Context(), roleID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取角色用户失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"role_id": roleID,
		"user_ids": userIDs,
	})
}

// RemoveUserRole 解绑用户角色
func (h *RoleHandler) RemoveUserRole(c echo.Context) error {
	userID := c.Param("user_id")
	roleID := c.Param("role_id")

	if err := h.roleRepo.RemoveUserRole(c.Request().Context(), userID, roleID); err != nil {
		return response.Error(c, http.StatusInternalServerError, "解绑用户角色失败")
	}

	return response.SuccessWithMsg[any](c, "用户角色解绑成功", nil)
}
