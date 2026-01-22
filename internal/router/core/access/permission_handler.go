package access

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

	permission := &Permission{
		Code:      req.Code,
		Name:      req.Name,
		Type:      req.Type,
		Remark:    req.Remark,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if err := h.repo.Create(c.Request().Context(), permission); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建失败")
	}

	return response.SuccessWithMsg(c, "创建成功", nil)
}

// GetPermissionDetail 获取权限详情
func (h *PermissionHandler) GetPermissionDetail(c echo.Context) error {
	id := c.Param("id")
	permission, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "权限不存在")
	}

	return response.Success(c, permission)
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
	permission.Remark = req.Remark
	permission.UpdatedBy = operatorID

	if err := h.repo.Update(c.Request().Context(), permission); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	return response.SuccessWithMsg(c, "更新成功", nil)
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

	return response.SuccessWithMsg(c, "删除成功", nil)
}

// ListPermissions 获取权限列表
func (h *PermissionHandler) ListPermissions(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	code := c.QueryParam("code")
	name := c.QueryParam("name")
	type_ := c.QueryParam("type")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	db := h.repo.GetDB(c.Request().Context())
	query := db.Model(&Permission{})

	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if type_ != "" {
		query = query.Where("type = ?", type_)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	var permissions []*Permission
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	return response.Success(c, map[string]interface{}{
		"list":  permissions,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
