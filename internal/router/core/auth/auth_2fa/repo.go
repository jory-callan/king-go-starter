package auth_2fa

import (
	"context"

	"gorm.io/gorm"
	"king-starter/internal/router/core/auth/auth_password"
)

// Repository 2FA 认证仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 2FA 认证仓库实例
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateTwoFA 创建 2FA 配置
func (r *Repository) CreateTwoFA(ctx context.Context, config *TwoFAConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// GetTwoFAByUserID 根据用户ID获取 2FA 配置
func (r *Repository) GetTwoFAByUserID(ctx context.Context, userID string) (*TwoFAConfig, error) {
	var config TwoFAConfig
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateTwoFA 更新 2FA 配置
func (r *Repository) UpdateTwoFA(ctx context.Context, config *TwoFAConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// CreateLoginLog 创建登录日志
func (r *Repository) CreateLoginLog(ctx context.Context, log *auth_password.LoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}
