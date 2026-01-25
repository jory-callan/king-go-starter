package auth_password

import (
	"context"

	"gorm.io/gorm"
)

// Repository 密码认证仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建密码认证仓库实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateRefreshToken 创建刷新令牌
func (r *Repository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// GetRefreshTokenByToken 根据令牌获取刷新令牌
func (r *Repository) GetRefreshTokenByToken(ctx context.Context, token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken 删除刷新令牌
func (r *Repository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&RefreshToken{}).Error
}

// CreateLoginLog 创建登录日志
func (r *Repository) CreateLoginLog(ctx context.Context, log *LoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
