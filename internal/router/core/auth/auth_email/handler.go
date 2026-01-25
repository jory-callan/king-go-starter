package auth_email

import (
	"net/http"

	"king-starter/internal/response"

	"github.com/labstack/echo/v4"
)

// EmailAuthReq 邮箱认证请求参数
type EmailAuthReq struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

// EmailHandler 邮箱认证处理器
type EmailHandler struct {
	// repo *Repository
}

// NewEmailHandler 创建邮箱认证处理器实例
func NewEmailHandler() *EmailHandler {
	return &EmailHandler{}
}

// SendVerificationCode 发送验证码
func (h *EmailHandler) SendVerificationCode(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return response.Error(c, http.StatusBadRequest, "邮箱不能为空")
	}

	// 这里应该是发送邮件验证码的逻辑
	// 为了示例，我们只返回成功

	return response.SuccessWithMsg[any](c, "验证码已发送", nil)
}

// VerifyEmail 验证邮箱
func (h *EmailHandler) VerifyEmail(c echo.Context) error {
	var req EmailAuthReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 这里应该是验证邮箱验证码的逻辑
	// 为了示例，我们假设验证成功

	// 记录邮箱认证成功日志
	// log := &core.LoginLog{
	// 	ID:        uuid.New().String(),
	// 	UserID:    "1", // 示例用户ID
	// 	Username:  req.Email,
	// 	AuthType:  core.AuthTypeEmail,
	// 	LoginType: core.LoginTypeSuccess,
	// 	IP:        c.RealIP(),
	// 	UserAgent: c.Request().UserAgent(),
	// 	Message:   "邮箱验证成功",
	// }
	// h.repo.CreateLoginLog(c.Request().Context(), log) // 暂时注释掉，因为repo未初始化

	return response.SuccessWithMsg[any](c, "邮箱验证成功", nil)
}
