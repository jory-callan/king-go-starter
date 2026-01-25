package permission

import (
	"fmt"
	"king-starter/internal/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PermissionHandler struct {
	repo *PermissionRepo
}

func NewPermissionHandler(repo *PermissionRepo) *PermissionHandler {
	return &PermissionHandler{repo: repo}
}

// CreatePermission 创建权限
func (h *PermissionHandler) CreatePermission(c echo.Context) error {
	var req CreatePermissionReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从 Context 获取

	permission := &CorePermission{
		Code:      req.Code,
		Name:      req.Name,
		Type:      req.Type,
		ParentID:  req.ParentID,
		Path:      req.Path,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if err := h.repo.Create(c.Request().Context(), permission); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建失败")
	}

	return response.SuccessWithMsg[any](c, "创建成功", nil)
}

// GetPermissionDetail 获取权限详情
func (h *PermissionHandler) GetPermissionDetail(c echo.Context) error {
	id := c.Param("id")
	permission, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "权限不存在")
	}

	return response.Success[any](c, permission)
}

// UpdatePermission 更新权限
func (h *PermissionHandler) UpdatePermission(c echo.Context) error {
	id := c.Param("id")
	var req UpdatePermissionReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从 Context 获取

	permission, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "权限不存在")
	}

	permission.Name = req.Name
	permission.Type = req.Type
	permission.ParentID = req.ParentID
	permission.Path = req.Path
	permission.Icon = req.Icon
	permission.Sort = req.Sort
	permission.Status = req.Status
	permission.Remark = req.Remark
	permission.UpdatedBy = operatorID

	if err := h.repo.Update(c.Request().Context(), permission); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	return response.SuccessWithMsg[any](c, "更新成功", nil)
}

// DeletePermission 删除权限
func (h *PermissionHandler) DeletePermission(c echo.Context) error {
	id := c.Param("id")
	operatorID := "system-admin" // TODO: 从 Context 获取

	permission, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "权限不存在")
	}

	// 先更新删除人
	permission.DeletedBy = operatorID
	if err := h.repo.Update(c.Request().Context(), permission); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	// 再删除权限
	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "删除失败")
	}

	return response.SuccessWithMsg[any](c, "删除成功", nil)
}

// ListPermissions 获取权限列表
func (h *PermissionHandler) ListPermissions(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 确保 NeedCount 为 true 以返回总数
	pq.NeedCount = true

	// 获取查询参数
	code := c.QueryParam("code")
	name := c.QueryParam("name")
	type_ := c.QueryParam("type")
	parentID := c.QueryParam("parent_id")
	statusStr := c.QueryParam("status")

	var status *int
	if statusStr != "" {
		statusVal := 0
		fmt.Sscanf(statusStr, "%d", &statusVal)
		status = &statusVal
	}

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
	if type_ != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("type = ?", type_)
		})
	}
	if parentID != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id = ?", parentID)
		})
	}
	if status != nil {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", *status)
		})
	}

	// 使用 BaseRepo 的分页方法
	result, err := h.repo.PaginationWithScopes(c.Request().Context(), &pq, scopes...)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	return response.SuccessPage[CorePermission](c, *result)
}

// GetPermissionTree 获取权限树结构
func (h *PermissionHandler) GetPermissionTree(c echo.Context) error {
	parentID := c.QueryParam("parent_id")
	if parentID == "" {
		parentID = "0"
	}

	permissions, err := h.repo.GetPermissionTree(c.Request().Context(), parentID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取权限树失败")
	}

	return response.Success[[]CorePermission](c, permissions)
}

// AssignRolePermissions 为角色分配权限
func (h *PermissionHandler) AssignRolePermissions(c echo.Context) error {
	roleID := c.Param("role_id")
	var req AssignRolePermissionsReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	if err := h.repo.AssignRolePermissions(c.Request().Context(), roleID, req.PermissionIDs); err != nil {
		return response.Error(c, http.StatusInternalServerError, "分配权限失败")
	}

	return response.SuccessWithMsg[any](c, "权限分配成功", nil)
}

// GetRolePermissions 获取角色拥有的权限
func (h *PermissionHandler) GetRolePermissions(c echo.Context) error {
	roleID := c.Param("role_id")
	permissions, err := h.repo.GetRolePermissions(c.Request().Context(), roleID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取角色权限失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"permission_ids": permissions,
	})
}

// GetRolePermissionsWithDetails 获取角色拥有的权限详细信息
func (h *PermissionHandler) GetRolePermissionsWithDetails(c echo.Context) error {
	roleID := c.Param("role_id")
	permissions, err := h.repo.GetRolePermissionsWithDetails(c.Request().Context(), roleID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取角色权限详情失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"role_id":      roleID,
		"permissions":  permissions,
	})
}

// GetUserAllPermissions 获取用户的所有权限
func (h *PermissionHandler) GetUserAllPermissions(c echo.Context) error {
	userID := c.Param("user_id")
	permissions, err := h.repo.GetUserAllPermissions(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取用户权限失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"user_id":     userID,
		"permissions": permissions,
	})
}

// RemoveRolePermissions 移除角色的部分或全部权限
func (h *PermissionHandler) RemoveRolePermissions(c echo.Context) error {
	roleID := c.Param("role_id")
	var req AssignRolePermissionsReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	if err := h.repo.RemoveRolePermissions(c.Request().Context(), roleID, req.PermissionIDs); err != nil {
		return response.Error(c, http.StatusInternalServerError, "移除角色权限失败")
	}

	msg := "角色权限移除成功"
	if len(req.PermissionIDs) == 0 {
		msg = "角色所有权限已清空"
	}

	return response.SuccessWithMsg[any](c, msg, nil)
}

// GetRolePermissionTree 获取角色权限树结构
func (h *PermissionHandler) GetRolePermissionTree(c echo.Context) error {
	roleID := c.Param("role_id")
	permissions, err := h.repo.GetRolePermissionTree(c.Request().Context(), roleID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "获取角色权限树失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"role_id": roleID,
		"tree":    permissions,
	})
}
