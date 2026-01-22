package access

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MenuHandler struct {
	repo *MenuRepo
}

func NewMenuHandler(repo *MenuRepo) *MenuHandler {
	return &MenuHandler{repo: repo}
}

// CreateMenu 创建菜单
func (h *MenuHandler) CreateMenu(c echo.Context) error {
	var req CreateMenuReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从 Context 获取

	menu := &CoreMenu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Path:      req.Path,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Status:    req.Status,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if err := h.repo.Create(c.Request().Context(), menu); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建失败")
	}

	return response.SuccessWithMsg[any](c, "创建成功", nil)
}

// GetMenuDetail 获取菜单详情
func (h *MenuHandler) GetMenuDetail(c echo.Context) error {
	id := c.Param("id")
	menu, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "菜单不存在")
	}

	return response.Success[any](c, menu)
}

// UpdateMenu 更新菜单
func (h *MenuHandler) UpdateMenu(c echo.Context) error {
	id := c.Param("id")
	var req UpdateMenuReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从 Context 获取

	menu, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "菜单不存在")
	}

	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.Path = req.Path
	menu.Icon = req.Icon
	menu.Sort = req.Sort
	menu.Status = req.Status
	menu.UpdatedBy = operatorID

	if err := h.repo.Update(c.Request().Context(), menu); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	return response.SuccessWithMsg[any](c, "更新成功", nil)
}

// DeleteMenu 删除菜单
func (h *MenuHandler) DeleteMenu(c echo.Context) error {
	id := c.Param("id")
	operatorID := "system-admin" // TODO: 从 Context 获取

	menu, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "菜单不存在")
	}

	// 先更新删除人
	menu.DeletedBy = operatorID
	if err := h.repo.Update(c.Request().Context(), menu); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新失败")
	}

	// 再删除菜单
	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, "删除失败")
	}

	return response.SuccessWithMsg[any](c, "删除成功", nil)
}

// ListMenus 获取菜单列表
func (h *MenuHandler) ListMenus(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 确保 NeedCount 为 true 以返回总数
	pq.NeedCount = true

	// 获取查询参数
	parentID := c.QueryParam("parent_id")
	name := c.QueryParam("name")
	statusStr := c.QueryParam("status")

	// 创建筛选条件的 scope 函数
	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	if parentID != "" {
		scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_id = ?", parentID)
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
	result, err := h.repo.PaginationWithScopes(c.Request().Context(), &pq, scopes...)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	return response.SuccessPage[CoreMenu](c, *result)
}
