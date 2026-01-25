package identity

import (
	"king-starter/internal/response"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// OAuthHandler OAuth 相关的处理器
type OAuthHandler struct {
	oauthRepo *OAuthRepo
}

// NewOAuthHandler 创建 OAuth 处理器实例
func NewOAuthHandler(oauthRepo *OAuthRepo) *OAuthHandler {
	return &OAuthHandler{
		oauthRepo: oauthRepo,
	}
}

// Authorize OAuth 授权
func (h *OAuthHandler) Authorize(c echo.Context) error {
	var req OAuthAuthorizeReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 验证客户端
	client, err := h.oauthRepo.GetClientByClientID(c.Request().Context(), req.ClientID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusUnauthorized, "客户端不存在")
		}
		return response.Error(c, http.StatusInternalServerError, "查询客户端失败")
	}

	// 验证重定向 URI
	if req.RedirectURI != client.RedirectURI {
		return response.Error(c, http.StatusBadRequest, "重定向 URI 错误")
	}

	// 生成授权码
	authCode := &CoreUserOAuthCode{
		ID:          uuid.New().String(),
		ClientID:    client.ClientID,
		UserID:      "1", // 模拟用户 ID
		Code:        uuid.New().String(),
		RedirectURI: req.RedirectURI,
		Scope:       req.Scope,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
	}

	if err := h.oauthRepo.CreateOAuthCode(c.Request().Context(), authCode); err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成授权码失败")
	}

	// 重定向到客户端，携带授权码
	return c.Redirect(http.StatusFound, req.RedirectURI+"?code="+authCode.Code+"&state="+req.State)
}

// GetToken OAuth 获取令牌
func (h *OAuthHandler) GetToken(c echo.Context) error {
	var req OAuthTokenReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 验证客户端
	client, err := h.oauthRepo.GetClientByClientID(c.Request().Context(), req.ClientID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusUnauthorized, "客户端不存在")
		}
		return response.Error(c, http.StatusInternalServerError, "查询客户端失败")
	}

	// 验证客户端密钥
	if req.ClientSecret != client.ClientSecret {
		return response.Error(c, http.StatusUnauthorized, "客户端密钥错误")
	}

	var userID string

	switch req.GrantType {
	case "authorization_code":
		// 验证授权码
		authCode, err := h.oauthRepo.GetOAuthCodeByCode(c.Request().Context(), req.Code)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return response.Error(c, http.StatusUnauthorized, "授权码不存在")
			}
			return response.Error(c, http.StatusInternalServerError, "查询授权码失败")
		}

		// 验证授权码是否过期
		if authCode.ExpiresAt.Before(time.Now()) {
			return response.Error(c, http.StatusUnauthorized, "授权码已过期")
		}

		// 验证客户端 ID
		if authCode.ClientID != req.ClientID {
			return response.Error(c, http.StatusUnauthorized, "客户端 ID 错误")
		}

		userID = authCode.UserID

		// 删除授权码（一次性使用）
		if err := h.oauthRepo.DeleteOAuthCode(c.Request().Context(), req.Code); err != nil {
			// 即使删除失败也继续执行
		}

	case "refresh_token":
		// 验证刷新令牌
		token, err := h.oauthRepo.GetOAuthTokenByRefreshToken(c.Request().Context(), req.RefreshToken)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return response.Error(c, http.StatusUnauthorized, "刷新令牌不存在")
			}
			return response.Error(c, http.StatusInternalServerError, "查询刷新令牌失败")
		}

		// 验证客户端 ID
		if token.ClientID != req.ClientID {
			return response.Error(c, http.StatusUnauthorized, "客户端 ID 错误")
		}

		userID = token.UserID

	case "password":
		// 模拟密码验证
		if req.Username != "admin" || req.Password != "123456" {
			return response.Error(c, http.StatusUnauthorized, "用户名或密码错误")
		}

		userID = "1"

	case "client_credentials":
		// 客户端凭证模式，直接使用客户端信息
		userID = "1"

	default:
		return response.Error(c, http.StatusBadRequest, "不支持的授权类型")
	}

	// 生成访问令牌和刷新令牌
	oauthToken := &CoreUserOAuthToken{
		ID:           uuid.New().String(),
		ClientID:     req.ClientID,
		UserID:       userID,
		AccessToken:  uuid.New().String(),
		RefreshToken: uuid.New().String(),
		Scope:        req.Scope,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	if err := h.oauthRepo.CreateOAuthToken(c.Request().Context(), oauthToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "生成令牌失败")
	}

	return response.Success[any](c, map[string]interface{}{
		"access_token":  oauthToken.AccessToken,
		"refresh_token": oauthToken.RefreshToken,
		"token_type":    oauthToken.TokenType,
		"expires_in":    86400, // 24 小时
		"scope":         oauthToken.Scope,
	})
}

// GetUserInfo OAuth 获取用户信息
func (h *OAuthHandler) GetUserInfo(c echo.Context) error {
	// 从请求头获取令牌
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return response.Error(c, http.StatusUnauthorized, "缺少授权头")
	}

	// 提取令牌
	accessToken := authHeader[7:] // 移除 "Bearer " 前缀

	// 验证令牌
	token, err := h.oauthRepo.GetOAuthTokenByAccessToken(c.Request().Context(), accessToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusUnauthorized, "令牌不存在")
		}
		return response.Error(c, http.StatusInternalServerError, "查询令牌失败")
	}

	// 验证令牌是否过期
	if token.ExpiresAt.Before(time.Now()) {
		return response.Error(c, http.StatusUnauthorized, "令牌已过期")
	}

	// 模拟用户信息
	userInfo := map[string]interface{}{
		"sub":                token.UserID,
		"name":               "管理员",
		"email":              "admin@example.com",
		"preferred_username": "admin",
	}

	return response.Success[any](c, userInfo)
}
