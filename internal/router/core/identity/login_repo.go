package identity

import (
	"context"
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

// LoginRepo 登录相关的数据访问层
type LoginRepo struct {
	*gormutil.BaseRepo[LoginLog]
}

// NewLoginRepo 创建登录数据访问层实例
func NewLoginRepo(db *gorm.DB) *LoginRepo {
	return &LoginRepo{BaseRepo: gormutil.NewBaseRepo[LoginLog](db)}
}

// CreateLoginLog 创建登录日志
func (r *LoginRepo) CreateLoginLog(ctx context.Context, log *LoginLog) error {
	return r.Create(ctx, log)
}

// CreateRefreshToken 创建刷新令牌
func (r *LoginRepo) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.GetDB(ctx).Create(token).Error
}

// GetRefreshTokenByToken 根据令牌获取刷新令牌
func (r *LoginRepo) GetRefreshTokenByToken(ctx context.Context, token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	err := r.GetDB(ctx).Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken 删除刷新令牌
func (r *LoginRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.GetDB(ctx).Where("token = ?", token).Delete(&RefreshToken{}).Error
}

// DeleteExpiredRefreshTokens 删除过期的刷新令牌
func (r *LoginRepo) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return r.GetDB(ctx).Where("expires_at < NOW()").Delete(&RefreshToken{}).Error
}
