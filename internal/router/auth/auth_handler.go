package auth

import (
	"king-starter/internal/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// Login 用户登录
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	token, user, err := h.authService.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return response.Error(c, http.StatusUnauthorized, err.Error())
	}

	return response.Success(c, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
		},
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	var req ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 从上下文获取用户ID
	userID := c.Get("user_id").(string)

	err := h.authService.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.SuccessWithMsg(c, "密码修改成功", nil)
}

// GetUserPermissions 获取用户权限列表
func (h *AuthHandler) GetUserPermissions(c echo.Context) error {
	// 从上下文获取用户ID
	userID := c.Get("user_id").(string)

	permissions, err := h.authService.GetUserPermissions(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, permissions)
}

// GetUserRoles 获取用户角色列表
func (h *AuthHandler) GetUserRoles(c echo.Context) error {
	// 从上下文获取用户ID
	userID := c.Get("user_id").(string)

	roles, err := h.authService.GetUserRoles(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err.Error())
	}

	return response.Success(c, roles)
}
