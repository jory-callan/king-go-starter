package access

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

	menu := &Menu{
		ParentID: req.ParentID,
		Name:     req.Name,
		Path:     req.Path,
		Icon:     req.Icon,
		Sort:     req.Sort,
		Status:   req.Status,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if err := h.repo.Create(c.Request().Context(), menu); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建失败")
	}

	return response.SuccessWithMsg(c, "创建成功", nil)
}

// GetMenuDetail 获取菜单详情
func (h *MenuHandler) GetMenuDetail(c echo.Context) error {
	id := c.Param("id")
	menu, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "菜单不存在")
	}

	return response.Success(c, menu)
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

	return response.SuccessWithMsg(c, "更新成功", nil)
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

	return response.SuccessWithMsg(c, "删除成功", nil)
}

// ListMenus 获取菜单列表
func (h *MenuHandler) ListMenus(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	parentID := c.QueryParam("parent_id")
	name := c.QueryParam("name")
	statusStr := c.QueryParam("status")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	db := h.repo.GetDB(c.Request().Context())
	query := db.Model(&Menu{})

	if parentID != "" {
		query = query.Where("parent_id = ?", parentID)
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

	var menus []*Menu
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&menus).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	return response.Success(c, map[string]interface{}{
		"list":  menus,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}
