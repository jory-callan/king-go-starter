package auth

import (
	"king-starter/internal/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"omitempty"`
	Status   int    `json:"status" validate:"omitempty,oneof=0 1"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       string `json:"id" validate:"required"`
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty"`
	Status   int    `json:"status" validate:"omitempty,oneof=0 1"`
}

// ListUsersRequest 查询用户列表请求
type ListUsersRequest struct {
	Page     int    `query:"page" validate:"omitempty,min=1"`
	PageSize int    `query:"page_size" validate:"omitempty,min=1,max=100"`
	Username string `query:"username" validate:"omitempty"`
	Email    string `query:"email" validate:"omitempty,email"`
	Phone    string `query:"phone" validate:"omitempty"`
	Status   int    `query:"status" validate:"omitempty,oneof=0 1"`
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	user := &User{
		Username: req.Username,
		Password: NewAuthService(h.userService.GetDB(c.Request().Context()), nil).hashPassword(req.Password),
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	}

	if user.Status == 0 {
		user.Status = 1 // 默认启用
	}

	err := h.userService.Create(c.Request().Context(), user, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建用户失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "创建用户成功", user)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c echo.Context) error {
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	// 获取现有用户信息
	existingUser, err := h.userService.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "用户不存在")
	}

	// 更新用户信息
	if req.Username != "" {
		existingUser.Username = req.Username
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Phone != "" {
		existingUser.Phone = req.Phone
	}
	if req.Status != 0 {
		existingUser.Status = req.Status
	}

	err = h.userService.Update(c.Request().Context(), existingUser, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新用户失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "更新用户成功", existingUser)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "用户ID不能为空")
	}

	// 从上下文获取当前用户ID
	operatorID := c.Get("user_id").(string)

	err := h.userService.Delete(c.Request().Context(), id, operatorID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "删除用户失败: "+err.Error())
	}

	return response.SuccessWithMsg(c, "删除用户成功", nil)
}

// GetUserByID 根据ID获取用户
func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.Error(c, http.StatusBadRequest, "用户ID不能为空")
	}

	user, err := h.userService.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, "用户不存在")
	}

	return response.Success(c, user)
}

// ListUsers 查询用户列表
func (h *UserHandler) ListUsers(c echo.Context) error {
	// 获取查询参数
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	username := c.QueryParam("username")
	email := c.QueryParam("email")
	phone := c.QueryParam("phone")
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
	if username != "" {
		conds["username"] = username
	}
	if email != "" {
		conds["email"] = email
	}
	if phone != "" {
		conds["phone"] = phone
	}
	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		conds["status"] = status
	}

	users, total, err := h.userService.ListByCondition(c.Request().Context(), conds, page, pageSize)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "查询用户列表失败: "+err.Error())
	}

	return response.Success(c, map[string]interface{}{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
