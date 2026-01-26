package auth_oauth2

import (
	"net/http"
	"time"

	"king-starter/internal/response"
	auth_password2 "king-starter/internal/router/core/auth/auth_password"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// OAuthAuthorizeReq OAuth 授权请求参数
type OAuthAuthorizeReq struct {
	ClientID     string `query:"client_id" validate:"required"`
	RedirectURI  string `query:"redirect_uri" validate:"required"`
	ResponseType string `query:"response_type" validate:"required,oneof=code token"`
	Scope        string `query:"scope"`
	State        string `query:"state"`
}

// OAuthTokenReq OAuth 获取令牌请求参数
type OAuthTokenReq struct {
	ClientID     string `json:"client_id" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	GrantType    string `json:"grant_type" validate:"required,oneof=authorization_code refresh_token password client_credentials"`
	Code         string `json:"code,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RedirectURI  string `json:"redirect_uri,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// OAuthUserInfoReq OAuth 获取用户信息请求参数
type OAuthUserInfoReq struct {
	AccessToken string `header:"Authorization" validate:"required"`
}

// OAuthHandler OAuth2 认证处理器
type OAuthHandler struct {
	repo *Repository
}

// NewOAuthHandler 创建 OAuth2 认证处理器实例
func NewOAuthHandler(repo *Repository) *OAuthHandler {
	return &OAuthHandler{
		repo: repo,
	}
}

// Authorize 授权接口
func (h *OAuthHandler) Authorize(c echo.Context) error {
	var req OAuthAuthorizeReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 验证客户端是否存在
	client, err := h.repo.GetClientByClientID(c.Request().Context(), req.ClientID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusBadRequest, "无效的客户端ID")
		}
		return response.Error(c, http.StatusInternalServerError, "查询客户端失败")
	}

	// 检查客户端状态
	if client.Status != 1 {
		return response.Error(c, http.StatusBadRequest, "客户端已被禁用")
	}

	// 验证回调地址
	if client.RedirectURI != req.RedirectURI {
		return response.Error(c, http.StatusBadRequest, "回调地址不匹配")
	}

	// 这里应该有用户授权页面逻辑
	// 为了示例，我们直接返回授权码
	code := uuid.New().String()

	authCode := &OAuthCode{
		ID:          uuid.New().String(),
		ClientID:    client.ClientID,
		UserID:      "1", // 示例用户ID
		Code:        code,
		RedirectURI: req.RedirectURI,
		Scope:       req.Scope,
		ExpiresAt:   time.Now().Add(10 * time.Minute), // 10分钟过期
	}

	if err := h.repo.CreateOAuthCode(c.Request().Context(), authCode); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建授权码失败")
	}

	// 记录 OAuth2 授权成功日志
	log := &auth_password2.CoreLoginLog{
		ID:        uuid.New().String(),
		UserID:    "1",
		Username:  "admin", // 示例用户名
		AuthType:  auth_password2.AuthTypeOAuth2,
		LoginType: auth_password2.LoginTypeSuccess,
		IP:        c.RealIP(),
		UserAgent: c.Request().UserAgent(),
		Message:   "OAuth2 授权成功",
	}
	h.repo.CreateLoginLog(c.Request().Context(), log)

	// 重定向到回调地址
	redirectURL := req.RedirectURI + "?code=" + code
	if req.State != "" {
		redirectURL += "&state=" + req.State
	}

	return c.Redirect(http.StatusFound, redirectURL)
}

// GetToken 获取令牌接口
func (h *OAuthHandler) GetToken(c echo.Context) error {
	var req OAuthTokenReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	// 验证客户端
	client, err := h.repo.GetClientByClientIDAndSecret(c.Request().Context(), req.ClientID, req.ClientSecret)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusBadRequest, "无效的客户端凭证")
		}
		return response.Error(c, http.StatusInternalServerError, "验证客户端失败")
	}

	// 检查客户端状态
	if client.Status != 1 {
		return response.Error(c, http.StatusBadRequest, "客户端已被禁用")
	}

	switch req.GrantType {
	case "authorization_code":
		// 使用授权码换取令牌
		return h.handleAuthorizationCode(c, req, client)
	case "refresh_token":
		// 使用刷新令牌
		return h.handleRefreshToken(c, req, client)
	default:
		return response.Error(c, http.StatusBadRequest, "不支持的授权类型")
	}
}

// 处理授权码方式
func (h *OAuthHandler) handleAuthorizationCode(c echo.Context, req OAuthTokenReq, client *OAuthClient) error {
	// 验证授权码
	authCode, err := h.repo.GetOAuthCodeByCode(c.Request().Context(), req.Code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusBadRequest, "无效的授权码")
		}
		return response.Error(c, http.StatusInternalServerError, "查询授权码失败")
	}

	// 检查授权码是否过期
	if authCode.ExpiresAt.Before(time.Now()) {
		return response.Error(c, http.StatusBadRequest, "授权码已过期")
	}

	// 验证客户端ID
	if authCode.ClientID != req.ClientID {
		return response.Error(c, http.StatusBadRequest, "客户端ID不匹配")
	}

	// 验证回调地址（如果提供）
	if req.RedirectURI != "" && authCode.RedirectURI != req.RedirectURI {
		return response.Error(c, http.StatusBadRequest, "回调地址不匹配")
	}

	// 生成访问令牌和刷新令牌
	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()

	oauthToken := &OAuthToken{
		ID:           uuid.New().String(),
		ClientID:     client.ClientID,
		UserID:       authCode.UserID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Scope:        authCode.Scope,
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(2 * time.Hour), // 2小时过期
	}

	if err := h.repo.CreateOAuthToken(c.Request().Context(), oauthToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建令牌失败")
	}

	// 删除已使用的授权码
	if err := h.repo.DeleteOAuthCode(c.Request().Context(), req.Code); err != nil {
		// 记录错误但不中断流程
	}

	return response.Success[any](c, map[string]interface{}{
		"access_token":  oauthToken.AccessToken,
		"refresh_token": oauthToken.RefreshToken,
		"token_type":    oauthToken.TokenType,
		"expires_in":    7200, // 2小时
		"scope":         oauthToken.Scope,
	})
}

// 处理刷新令牌方式
func (h *OAuthHandler) handleRefreshToken(c echo.Context, req OAuthTokenReq, client *OAuthClient) error {
	// 获取刷新令牌
	token, err := h.repo.GetOAuthTokenByRefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusBadRequest, "无效的刷新令牌")
		}
		return response.Error(c, http.StatusInternalServerError, "查询刷新令牌失败")
	}

	// 检查刷新令牌是否过期
	if token.ExpiresAt.Before(time.Now()) {
		return response.Error(c, http.StatusBadRequest, "刷新令牌已过期")
	}

	// 检查客户端ID是否匹配
	if token.ClientID != req.ClientID {
		return response.Error(c, http.StatusBadRequest, "客户端ID不匹配")
	}

	// 生成新的访问令牌
	newAccessToken := uuid.New().String()
	newRefreshToken := uuid.New().String()

	newToken := &OAuthToken{
		ID:           uuid.New().String(),
		ClientID:     token.ClientID,
		UserID:       token.UserID,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		Scope:        token.Scope,
		TokenType:    token.TokenType,
		ExpiresAt:    time.Now().Add(2 * time.Hour),
	}

	if err := h.repo.CreateOAuthToken(c.Request().Context(), newToken); err != nil {
		return response.Error(c, http.StatusInternalServerError, "创建新令牌失败")
	}

	// 删除旧的令牌
	if err := h.repo.DeleteOAuthToken(c.Request().Context(), token.AccessToken); err != nil {
		// 记录错误但不中断流程
	}

	return response.Success[any](c, map[string]interface{}{
		"access_token":  newToken.AccessToken,
		"refresh_token": newToken.RefreshToken,
		"token_type":    newToken.TokenType,
		"expires_in":    7200, // 2小时
		"scope":         newToken.Scope,
	})
}

// GetUserInfo 获取用户信息接口
func (h *OAuthHandler) GetUserInfo(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return response.Error(c, http.StatusUnauthorized, "缺少授权头")
	}

	// 解析 Bearer token
	var accessToken string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken = authHeader[7:]
	} else {
		return response.Error(c, http.StatusUnauthorized, "无效的授权头格式")
	}

	// 验证访问令牌
	token, err := h.repo.GetOAuthTokenByAccessToken(c.Request().Context(), accessToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusUnauthorized, "无效的访问令牌")
		}
		return response.Error(c, http.StatusInternalServerError, "验证令牌失败")
	}

	// 检查令牌是否过期
	if token.ExpiresAt.Before(time.Now()) {
		return response.Error(c, http.StatusUnauthorized, "访问令牌已过期")
	}

	// 返回用户信息
	return response.Success[any](c, map[string]interface{}{
		"user_id":   token.UserID,
		"client_id": token.ClientID,
		"scope":     token.Scope,
	})
}
