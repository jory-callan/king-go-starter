package user

import (
	"errors"
	"fmt"
	"king-starter/internal/app"
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	repo *Repository
}

func NewHandler(app *app.App) *Handler {
	return &Handler{repo: NewRepository(app.Db.DB)}
}

// Create 创建用户
func (h *Handler) Create(c echo.Context) error {
	var req CreateUserReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// TODO: 加入参数校验逻辑，如使用 go-playground/validate

	// 模拟从 Context 中获取当前操作人 ID (后续中间件实现)
	operatorID := "system-admin"

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

	user := &User{
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

	return response.Success(c, user)
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
	return response.Success(c, user)
}

// List 用户列表
func (h *Handler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var users []*User
	var total int64

	db := h.repo.GetDB(c.Request().Context())

	// 计算总数
	if err := db.Model(&User{}).Count(&total).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询失败")
	}

	return response.Success(c, map[string]interface{}{
		"list":  users,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

// Update 更新用户
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	operatorID := "system-admin" // TODO: 从中间件获取

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

	if err := h.repo.GetDB(c.Request().Context()).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMsg(c, "更新成功", nil)
}

// Delete 删除用户
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	operatorID := "system-admin" // TODO: 从中间件获取

	// GORM 的软删除默认只更新 deleted_at，我们这里手动处理 deleted_by
	if err := h.repo.GetDB(c.Request().Context()).Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": gorm.Expr("NOW()"), // 或者使用 time.Now()
		"deleted_by": operatorID,
	}).Error; err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMsg(c, "删除成功", nil)
}
