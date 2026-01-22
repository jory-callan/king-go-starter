package identity

import (
	"king-starter/internal/response"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// LoginHandler 登录相关的处理器
type LoginHandler struct {
	loginRepo *LoginRepo
}

// NewLoginHandler 创建登录处理器实例
func NewLoginHandler(loginRepo *LoginRepo) *LoginHandler {
	return &LoginHandler{
		loginRepo: loginRepo,
	}
}

// Claims 自定义 JWT Claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login 用户登录
func (h *LoginHandler) Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 模拟用户验证（实际应该从数据库查询并验证密码）
	if req.Username != "admin" || req.Password != "123456" {
		// 记录登录失败日志
		loginLog := &CoreLoginLog{
			ID:        uuid.New().String(),
			UserID:    "",
			Username:  req.Username,
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
			Status:    0,
			Message:   "用户名或密码错误",
		}
		h.loginRepo.CreateLoginLog(c.Request().Context(), loginLog)

		return response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
	}

	// 生成 JWT Token
	claims := Claims{
		UserID:   "1",
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "king-starter",
			Subject:   "user-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成令牌失败")
	}

	// 生成刷新令牌
	refreshToken := &CoreRefreshToken{
		ID:        uuid.New().String(),
		UserID:    "1",
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := h.loginRepo.CreateRefreshToken(c.Request().Context(), refreshToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成刷新令牌失败")
	}

	// 记录登录成功日志
	loginLog := &CoreLoginLog{
		ID:        uuid.New().String(),
		UserID:    "1",
		Username:  req.Username,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		Status:    1,
		Message:   "登录成功",
	}
	h.loginRepo.CreateLoginLog(c.Request().Context(), loginLog)

	return response.Success[any](c, map[string]interface{}{
		"access_token":  tokenString,
		"refresh_token": refreshToken.Token,
		"expires_at":    claims.ExpiresAt,
		"user": map[string]interface{}{
			"id":       "1",
			"username": req.Username,
			"name":     "管理员",
		},
	})
}

// Logout 用户登出
func (h *LoginHandler) Logout(c echo.Context) error {
	var req LogoutReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 清除刷新令牌
	if req.RefreshToken != "" {
		if err := h.loginRepo.DeleteRefreshToken(c.Request().Context(), req.RefreshToken); err != nil {
			// 即使删除失败也继续执行
		}
	}

	return response.SuccessWithMsg[any](c, "登出成功", nil)
}

// RefreshToken 刷新令牌
func (h *LoginHandler) RefreshToken(c echo.Context) error {
	var req RefreshTokenReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 获取刷新令牌
	refreshToken, err := h.loginRepo.GetRefreshTokenByToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusUnauthorized, "刷新令牌不存在")
		}
		return response.Error(c, http.StatusInternalServerError, "查询刷新令牌失败")
	}

	// 检查刷新令牌是否过期
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return response.Error(c, http.StatusUnauthorized, "刷新令牌已过期")
	}

	// 生成新的访问令牌
	claims := Claims{
		UserID:   refreshToken.UserID,
		Username: "admin", // 实际应该从数据库查询
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "king-starter",
			Subject:   "user-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成令牌失败")
	}

	// 生成新的刷新令牌
	newRefreshToken := &CoreRefreshToken{
		ID:        uuid.New().String(),
		UserID:    refreshToken.UserID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := h.loginRepo.CreateRefreshToken(c.Request().Context(), newRefreshToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成刷新令牌失败")
	}

	// 删除旧的刷新令牌
	if err := h.loginRepo.DeleteRefreshToken(c.Request().Context(), req.RefreshToken); err != nil {
		// 即使删除失败也继续执行
	}

	return response.Success[any](c, map[string]interface{}{
		"access_token":  tokenString,
		"refresh_token": newRefreshToken.Token,
		"expires_at":    claims.ExpiresAt,
	})
}
