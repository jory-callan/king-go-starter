package auth_2fa

import (
	"net/http"

	"king-starter/internal/response"
	"king-starter/internal/router/core/auth/auth_password"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

// TwoFAHandler 2FA 认证处理器
type TwoFAHandler struct {
	repo *Repository
}

// NewTwoFAHandler 创建 2FA 认证处理器实例
func NewTwoFAHandler(repo *Repository) *TwoFAHandler {
	return &TwoFAHandler{
		repo: repo,
	}
}

// EnableTwoFA 启用 2FA
func (h *TwoFAHandler) EnableTwoFA(c echo.Context) error {
	var req EnableTwoFAReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 检查用户是否已有 2FA 配置
	existingTwoFA, err := h.repo.GetTwoFAByUserID(c.Request().Context(), req.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return response.Error(c, http.StatusInternalServerError, "查询 2FA 配置失败")
	}

	var secret string
	if existingTwoFA == nil {
		// 生成新的 TOTP 密钥
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "King Starter",
			AccountName: "user@example.com", // 实际应该使用用户邮箱
		})
		if err != nil {
			return response.Error(c, http.StatusInternalServerError, "生成密钥失败")
		}
		secret = key.Secret()
	} else {
		secret = existingTwoFA.Secret
	}

	// 验证 TOTP 码
	valid := totp.Validate(req.Code, secret)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "验证失败")
	}
	if !valid {
		// 记录验证失败日志
		log := &auth_password.LoginLog{
			ID:        uuid.New().String(),
			UserID:    req.UserID,
			Username:  "", // 暂时为空，后续可以优化
			AuthType:  AuthType2FA,
			LoginType: LoginTypeFailed,
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
			Message:   "验证码错误",
		}
		h.repo.CreateLoginLog(c.Request().Context(), log)

		return response.Error(c, http.StatusBadRequest, "验证码错误")
	}

	// 保存 2FA 配置
	if existingTwoFA == nil {
		twoFA := &TwoFAConfig{
			ID:     uuid.New().String(),
			UserID: req.UserID,
			Secret: secret,
			Status: 1,
		}
		if err := h.repo.CreateTwoFA(c.Request().Context(), twoFA); err != nil {
			return response.Error(c, http.StatusInternalServerError, "保存 2FA 配置失败")
		}
	} else {
		existingTwoFA.Status = 1
		if err := h.repo.UpdateTwoFA(c.Request().Context(), existingTwoFA); err != nil {
			return response.Error(c, http.StatusInternalServerError, "更新 2FA 配置失败")
		}
	}

	// 记录验证成功日志
	log := &auth_password.LoginLog{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Username:  "", // 暂时为空，后续可以优化
		AuthType:  AuthType2FA,
		LoginType: LoginTypeSuccess,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		Message:   "启用 2FA 成功",
	}
	h.repo.CreateLoginLog(c.Request().Context(), log)

	return response.SuccessWithMsg[any](c, "启用 2FA 成功", nil)
}

// VerifyTwoFA 验证 2FA
func (h *TwoFAHandler) VerifyTwoFA(c echo.Context) error {
	var req VerifyTwoFAReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 获取用户 2FA 配置
	twoFA, err := h.repo.GetTwoFAByUserID(c.Request().Context(), req.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "2FA 配置不存在")
		}
		return response.Error(c, http.StatusInternalServerError, "查询 2FA 配置失败")
	}

	// 检查 2FA 是否已启用
	if twoFA.Status != 1 {
		return response.Error(c, http.StatusBadRequest, "2FA 未启用")
	}

	// 验证 TOTP 码
	valid := totp.Validate(req.Code, twoFA.Secret)
	if !valid {
		// 记录验证失败日志
		log := &auth_password.LoginLog{
			ID:        uuid.New().String(),
			UserID:    req.UserID,
			Username:  "", // 暂时为空，后续可以优化
			AuthType:  AuthType2FA,
			LoginType: LoginTypeFailed,
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
			Message:   "验证码错误",
		}
		h.repo.CreateLoginLog(c.Request().Context(), log)

		return response.Error(c, http.StatusBadRequest, "验证码错误")
	}

	// 记录验证成功日志
	log := &auth_password.LoginLog{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Username:  "", // 暂时为空，后续可以优化
		AuthType:  AuthType2FA,
		LoginType: LoginTypeSuccess,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		Message:   "验证成功",
	}
	h.repo.CreateLoginLog(c.Request().Context(), log)

	return response.SuccessWithMsg[any](c, "验证成功", nil)
}

// DisableTwoFA 禁用 2FA
func (h *TwoFAHandler) DisableTwoFA(c echo.Context) error {
	var req DisableTwoFAReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 获取用户 2FA 配置
	twoFA, err := h.repo.GetTwoFAByUserID(c.Request().Context(), req.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "2FA 配置不存在")
		}
		return response.Error(c, http.StatusInternalServerError, "查询 2FA 配置失败")
	}

	// 验证 TOTP 码
	valid := totp.Validate(req.Code, twoFA.Secret)
	if !valid {
		// 记录验证失败日志
		log := &auth_password.LoginLog{
			ID:        uuid.New().String(),
			UserID:    req.UserID,
			Username:  "", // 暂时为空，后续可以优化
			AuthType:  AuthType2FA,
			LoginType: LoginTypeFailed,
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
			Message:   "验证码错误",
		}
		h.repo.CreateLoginLog(c.Request().Context(), log)

		return response.Error(c, http.StatusBadRequest, "验证码错误")
	}

	// 禁用 2FA
	twoFA.Status = 0
	if err := h.repo.UpdateTwoFA(c.Request().Context(), twoFA); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新 2FA 配置失败")
	}

	// 记录验证成功日志
	log := &auth_password.LoginLog{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Username:  "", // 暂时为空，后续可以优化
		AuthType:  AuthType2FA,
		LoginType: LoginTypeSuccess,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		Message:   "禁用 2FA 成功",
	}
	h.repo.CreateLoginLog(c.Request().Context(), log)

	return response.SuccessWithMsg[any](c, "禁用 2FA 成功", nil)
}
