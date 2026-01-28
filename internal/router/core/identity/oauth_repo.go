package identity

import (
	"context"

	"king-starter/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

// OAuthRepo OAuth 相关的数据访问层
type OAuthRepo struct {
	*gormutil.BaseRepo[CoreUserOAuthClient]
}

// NewOAuthRepo 创建 OAuth 数据访问层实例
func NewOAuthRepo(db *gorm.DB) *OAuthRepo {
	return &OAuthRepo{BaseRepo: gormutil.NewBaseRepo[CoreUserOAuthClient](db)}
}

// GetClientByClientID 根据客户端 ID 获取 OAuth 客户端
func (r *OAuthRepo) GetClientByClientID(ctx context.Context, clientID string) (*CoreUserOAuthClient, error) {
	var client CoreUserOAuthClient
	err := r.GetDB(ctx).Where("client_id = ? AND status = 1", clientID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// CreateOAuthCode 创建 OAuth 授权码
func (r *OAuthRepo) CreateOAuthCode(ctx context.Context, code *CoreUserOAuthCode) error {
	return r.GetDB(ctx).Create(code).Error
}

// GetOAuthCodeByCode 根据授权码获取 OAuth 授权码
func (r *OAuthRepo) GetOAuthCodeByCode(ctx context.Context, code string) (*CoreUserOAuthCode, error) {
	var oauthCode CoreUserOAuthCode
	err := r.GetDB(ctx).Where("code = ?", code).First(&oauthCode).Error
	if err != nil {
		return nil, err
	}
	return &oauthCode, nil
}

// DeleteOAuthCode 删除 OAuth 授权码
func (r *OAuthRepo) DeleteOAuthCode(ctx context.Context, code string) error {
	return r.GetDB(ctx).Where("code = ?", code).Delete(&CoreUserOAuthCode{}).Error
}

// CreateOAuthToken 创建 OAuth 令牌
func (r *OAuthRepo) CreateOAuthToken(ctx context.Context, token *CoreUserOAuthToken) error {
	return r.GetDB(ctx).Create(token).Error
}

// GetOAuthTokenByAccessToken 根据访问令牌获取 OAuth 令牌
func (r *OAuthRepo) GetOAuthTokenByAccessToken(ctx context.Context, accessToken string) (*CoreUserOAuthToken, error) {
	var oauthToken CoreUserOAuthToken
	err := r.GetDB(ctx).Where("access_token = ?", accessToken).First(&oauthToken).Error
	if err != nil {
		return nil, err
	}
	return &oauthToken, nil
}

// GetOAuthTokenByRefreshToken 根据刷新令牌获取 OAuth 令牌
func (r *OAuthRepo) GetOAuthTokenByRefreshToken(ctx context.Context, refreshToken string) (*CoreUserOAuthToken, error) {
	var oauthToken CoreUserOAuthToken
	err := r.GetDB(ctx).Where("refresh_token = ?", refreshToken).First(&oauthToken).Error
	if err != nil {
		return nil, err
	}
	return &oauthToken, nil
}

// DeleteExpiredOAuthCodes 删除过期的 OAuth 授权码
func (r *OAuthRepo) DeleteExpiredOAuthCodes(ctx context.Context) error {
	return r.GetDB(ctx).Where("expires_at < NOW()").Delete(&CoreUserOAuthCode{}).Error
}

// DeleteExpiredOAuthTokens 删除过期的 OAuth 令牌
func (r *OAuthRepo) DeleteExpiredOAuthTokens(ctx context.Context) error {
	return r.GetDB(ctx).Where("expires_at < NOW()").Delete(&CoreUserOAuthToken{}).Error
}
