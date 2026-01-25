package auth_oauth2

import (
	"context"

	"gorm.io/gorm"
	"king-starter/internal/router/core/auth/auth_password"
)

// Repository OAuth2 认证仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 OAuth2 认证仓库实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetClientByClientID 根据客户端ID获取客户端
func (r *Repository) GetClientByClientID(ctx context.Context, clientID string) (*OAuthClient, error) {
	var client OAuthClient
	err := r.db.WithContext(ctx).Where("client_id = ?", clientID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetClientByClientIDAndSecret 根据客户端ID和密钥获取客户端
func (r *Repository) GetClientByClientIDAndSecret(ctx context.Context, clientID, clientSecret string) (*OAuthClient, error) {
	var client OAuthClient
	err := r.db.WithContext(ctx).Where("client_id = ? AND client_secret = ?", clientID, clientSecret).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// CreateOAuthCode 创建 OAuth 授权码
func (r *Repository) CreateOAuthCode(ctx context.Context, code *OAuthCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

// GetOAuthCodeByCode 根据授权码获取记录
func (r *Repository) GetOAuthCodeByCode(ctx context.Context, code string) (*OAuthCode, error) {
	var authCode OAuthCode
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&authCode).Error
	if err != nil {
		return nil, err
	}
	return &authCode, nil
}

// DeleteOAuthCode 删除 OAuth 授权码
func (r *Repository) DeleteOAuthCode(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Where("code = ?", code).Delete(&OAuthCode{}).Error
}

// CreateOAuthToken 创建 OAuth 令牌
func (r *Repository) CreateOAuthToken(ctx context.Context, token *OAuthToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetOAuthTokenByRefreshToken 根据刷新令牌获取令牌
func (r *Repository) GetOAuthTokenByRefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	var token OAuthToken
	err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetOAuthTokenByAccessToken 根据访问令牌获取令牌
func (r *Repository) GetOAuthTokenByAccessToken(ctx context.Context, accessToken string) (*OAuthToken, error) {
	var token OAuthToken
	err := r.db.WithContext(ctx).Where("access_token = ?", accessToken).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteOAuthToken 删除 OAuth 令牌
func (r *Repository) DeleteOAuthToken(ctx context.Context, accessToken string) error {
	return r.db.WithContext(ctx).Where("access_token = ?", accessToken).Delete(&OAuthToken{}).Error
}

// CreateLoginLog 创建登录日志
func (r *Repository) CreateLoginLog(ctx context.Context, log *auth_password.LoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
