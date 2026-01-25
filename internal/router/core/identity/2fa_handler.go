package identity

import (
	"king-starter/internal/response"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

// TwoFAHandler 2FA 相关的处理器
type TwoFAHandler struct {
	twoFARepo *TwoFARepo
}

// NewTwoFAHandler 创建 2FA 处理器实例
func NewTwoFAHandler(twoFARepo *TwoFARepo) *TwoFAHandler {
	return &TwoFAHandler{
		twoFARepo: twoFARepo,
	}
}

// EnableTwoFA 启用 2FA
func (h *TwoFAHandler) EnableTwoFA(c echo.Context) error {
	var req EnableTwoFAReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 检查用户是否已有 2FA 配置
	existingTwoFA, err := h.twoFARepo.GetTwoFAByUserID(c.Request().Context(), req.UserID)
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
		log := &CoreUserTwoFALog{
			ID:        uuid.New().String(),
			UserID:    req.UserID,
			Status:    0,
			Message:   "验证码错误",
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}
		h.twoFARepo.CreateTwoFALog(c.Request().Context(), log)

		return response.Error(c, http.StatusBadRequest, "验证码错误")
	}

	// 保存 2FA 配置
	if existingTwoFA == nil {
		twoFA := &CoreUserTwoFA{
			ID:     uuid.New().String(),
			UserID: req.UserID,
			Secret: secret,
			Status: 1,
		}
		if err := h.twoFARepo.CreateTwoFA(c.Request().Context(), twoFA); err != nil {
			return response.Error(c, http.StatusInternalServerError, "保存 2FA 配置失败")
		}
	} else {
		existingTwoFA.Status = 1
		if err := h.twoFARepo.UpdateTwoFA(c.Request().Context(), existingTwoFA); err != nil {
			return response.Error(c, http.StatusInternalServerError, "更新 2FA 配置失败")
		}
	}

	// 记录验证成功日志
	log := &CoreUserTwoFALog{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Status:    1,
		Message:   "启用 2FA 成功",
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	h.twoFARepo.CreateTwoFALog(c.Request().Context(), log)

	return response.SuccessWithMsg[any](c, "启用 2FA 成功", nil)
}

// VerifyTwoFA 验证 2FA
func (h *TwoFAHandler) VerifyTwoFA(c echo.Context) error {
	var req VerifyTwoFAReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 获取用户 2FA 配置
	twoFA, err := h.twoFARepo.GetTwoFAByUserID(c.Request().Context(), req.UserID)
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
		log := &CoreUserTwoFALog{
			ID:        uuid.New().String(),
			UserID:    req.UserID,
			Status:    0,
			Message:   "验证码错误",
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}
		h.twoFARepo.CreateTwoFALog(c.Request().Context(), log)

		return response.Error(c, http.StatusBadRequest, "验证码错误")
	}

	// 记录验证成功日志
	log := &CoreUserTwoFALog{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Status:    1,
		Message:   "验证成功",
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	h.twoFARepo.CreateTwoFALog(c.Request().Context(), log)

	return response.SuccessWithMsg[any](c, "验证成功", nil)
}

// DisableTwoFA 禁用 2FA
func (h *TwoFAHandler) DisableTwoFA(c echo.Context) error {
	var req DisableTwoFAReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 获取用户 2FA 配置
	twoFA, err := h.twoFARepo.GetTwoFAByUserID(c.Request().Context(), req.UserID)
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
		log := &CoreUserTwoFALog{
			ID:        uuid.New().String(),
			UserID:    req.UserID,
			Status:    0,
			Message:   "验证码错误",
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
		}
		h.twoFARepo.CreateTwoFALog(c.Request().Context(), log)

		return response.Error(c, http.StatusBadRequest, "验证码错误")
	}

	// 禁用 2FA
	twoFA.Status = 0
	if err := h.twoFARepo.UpdateTwoFA(c.Request().Context(), twoFA); err != nil {
		return response.Error(c, http.StatusInternalServerError, "更新 2FA 配置失败")
	}

	// 记录验证成功日志
	log := &CoreUserTwoFALog{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Status:    1,
		Message:   "禁用 2FA 成功",
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
	}
	h.twoFARepo.CreateTwoFALog(c.Request().Context(), log)

	return response.SuccessWithMsg[any](c, "禁用 2FA 成功", nil)
}
