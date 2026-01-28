package user

import (
	"errors"
	"fmt"
	"king-starter/internal/response"
	"king-starter/pkg/goutils/echoutil"
	"king-starter/pkg/goutils/idutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Create 创建用户
func (h *Handler) Create(c echo.Context) error {
	var req CreateUserReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 检查用户名是否存在
	exist, _ := h.repo.GetByUsername(c.Request().Context(), req.Username)
	if exist != nil {
		return response.Error(c, http.StatusInternalServerError, "用户名已存在")
	}

	// 密码加密
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, fmt.Sprintf("密码加密失败: %v", err))
	}

	id := idutil.ShortUUIDv7()
	operatorID := echoutil.GetUserID(c)
	if operatorID == "" {
		operatorID = id
	}
	user := &CoreUser{
		ID:        id,
		Username:  req.Username,
		Password:  string(hashedBytes),
		Nickname:  req.Nickname,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    1, // 默认启用
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	if err := h.repo.Create(c.Request().Context(), user); err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success[any](c, user)
}

// GetByID 获取用户详情
func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "ID 不能为空")
	}

	user, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "用户不存在")
	}

	// 隐藏密码哈希
	user.Password = ""
	return response.Success[any](c, user)
}

// List 用户列表
func (h *Handler) List(c echo.Context) error {
	var pq response.PageQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 确保 NeedCount 为 true 以返回总数
	pq.NeedCount = true

	// 使用 BaseRepo 的分页方法
	result, err := h.repo.Pagination(c.Request().Context(), &pq, nil)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	// 隐藏密码哈希
	for _, user := range result.Items {
		user.Password = ""
	}

	return response.SuccessPage[CoreUser](c, *result)
}

// Update 更新用户
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := echoutil.GetUserID(c)
	if operatorID == "" {
		operatorID = id
	}

	// 检查是否存在
	_, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, http.StatusInternalServerError, "用户不存在")
		}
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	updates := map[string]interface{}{
		"nickname":   req.Nickname,
		"email":      req.Email,
		"phone":      req.Phone,
		"updated_by": operatorID,
	}

	if err := h.repo.GetDB(c.Request().Context()).Model(&CoreUser{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMsg[any](c, "更新成功", nil)
}

// Delete 删除用户
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	operatorID := echoutil.GetUserID(c)
	if operatorID == "" {
		operatorID = id
	}

	if err := h.repo.Delete(c.Request().Context(), id, operatorID); err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMsg[any](c, "删除成功", nil)
}
