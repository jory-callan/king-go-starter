package auth_password

import (
	"net/http"
	"time"

	"king-starter/internal/response"
	"king-starter/internal/router/core/user"
	"king-starter/pkg/goutils/idutil"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// LoginHandler 密码登录处理器
type LoginHandler struct {
	repo     *Repository
	userRepo *user.Repository
}

// NewLoginHandler 创建密码登录处理器实例
func NewLoginHandler(repo *Repository, userRepo *user.Repository) *LoginHandler {
	return &LoginHandler{
		repo:     repo,
		userRepo: userRepo,
	}
}

// Claims 自定义 JWT Claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (h *LoginHandler) Register(c echo.Context) error {
	var req RegisterReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}
	// 检查用户邮箱是否已存在
	var existingUser user.CoreUser
	if err := h.userRepo.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return response.Error(c, http.StatusBadRequest, "邮箱已被注册")
	}
	// 检查用户手机号是否已存在
	if err := h.userRepo.DB.Where("phone = ?", req.Phone).First(&existingUser).Error; err == nil {
		return response.Error(c, http.StatusBadRequest, "手机号已被注册")
	}
	// 创建新用户
	newUser := &user.CoreUser{
		ID:       idutil.ShortUUIDv7(),
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
	}
	err := h.userRepo.Create(c.Request().Context(), newUser)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "注册用户失败")
	}

	return response.SuccessWithMsg(c, "注册成功", *newUser)
}

// Login 用户登录
func (h *LoginHandler) Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 模拟用户验证（实际应该从数据库查询并验证密码）
	// 在实际应用中，应从用户服务获取用户信息并验证密码
	userID := "1" // 示例用户ID
	username := req.Username

	// 这里应该实际验证用户凭据
	// 示例：验证密码是否正确
	if req.Username != "admin" || req.Password != "123456" {
		// 记录登录失败日志
		loginLog := &CoreLoginLog{
			ID:        uuid.New().String(),
			UserID:    "", // 登录失败时没有用户ID
			Username:  req.Username,
			AuthType:  AuthTypePassword,
			LoginType: LoginTypeFailed,
			IP:        c.RealIP(),
			UserAgent: c.Request().UserAgent(),
			Message:   "用户名或密码错误",
		}
		h.repo.CreateLoginLog(c.Request().Context(), loginLog)

		return response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
	}

	// 生成 JWT Token
	claims := Claims{
		UserID:   userID,
		Username: username,
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
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := h.repo.CreateRefreshToken(c.Request().Context(), refreshToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成刷新令牌失败")
	}

	// 记录登录成功日志
	loginLog := &CoreLoginLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Username:  username,
		AuthType:  AuthTypePassword,
		LoginType: LoginTypeSuccess,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		Message:   "登录成功",
	}
	h.repo.CreateLoginLog(c.Request().Context(), loginLog)

	return response.Success[any](c, map[string]interface{}{
		"access_token":  tokenString,
		"refresh_token": refreshToken.Token,
		"expires_at":    claims.ExpiresAt,
		"user": map[string]interface{}{
			"id":       userID,
			"username": username,
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
		if err := h.repo.DeleteRefreshToken(c.Request().Context(), req.RefreshToken); err != nil {
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
	refreshToken, err := h.repo.GetRefreshTokenByToken(c.Request().Context(), req.RefreshToken)
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

	if err := h.repo.CreateRefreshToken(c.Request().Context(), newRefreshToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成刷新令牌失败")
	}

	// 删除旧的刷新令牌
	if err := h.repo.DeleteRefreshToken(c.Request().Context(), req.RefreshToken); err != nil {
		// 即使删除失败也继续执行
	}

	return response.Success[any](c, map[string]interface{}{
		"access_token":  tokenString,
		"refresh_token": newRefreshToken.Token,
		"expires_at":    claims.ExpiresAt,
	})
}
