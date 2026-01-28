package identity

import (
	"context"

	"king-starter/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

// TwoFARepo 2FA 相关的数据访问层
type TwoFARepo struct {
	*gormutil.BaseRepo[CoreUserTwoFA]
}

// NewTwoFARepo 创建 2FA 数据访问层实例
func NewTwoFARepo(db *gorm.DB) *TwoFARepo {
	return &TwoFARepo{BaseRepo: gormutil.NewBaseRepo[CoreUserTwoFA](db)}
}

// GetTwoFAByUserID 根据用户 ID 获取 2FA 配置
func (r *TwoFARepo) GetTwoFAByUserID(ctx context.Context, userID string) (*CoreUserTwoFA, error) {
	var twoFA CoreUserTwoFA
	err := r.GetDB(ctx).Where("user_id = ?", userID).First(&twoFA).Error
	if err != nil {
		return nil, err
	}
	return &twoFA, nil
}

// CreateTwoFA 创建 2FA 配置
func (r *TwoFARepo) CreateTwoFA(ctx context.Context, twoFA *CoreUserTwoFA) error {
	return r.Create(ctx, twoFA)
}

// UpdateTwoFA 更新 2FA 配置
func (r *TwoFARepo) UpdateTwoFA(ctx context.Context, twoFA *CoreUserTwoFA) error {
	return r.Update(ctx, twoFA)
}

// DeleteTwoFA 删除 2FA 配置
func (r *TwoFARepo) DeleteTwoFA(ctx context.Context, userID string) error {
	return r.GetDB(ctx).Where("user_id = ?", userID).Delete(&CoreUserTwoFA{}).Error
}

// CreateTwoFALog 创建 2FA 验证日志
func (r *TwoFARepo) CreateTwoFALog(ctx context.Context, log *CoreUserTwoFALog) error {
	return r.GetDB(ctx).Create(log).Error
}
