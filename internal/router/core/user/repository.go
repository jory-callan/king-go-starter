package user

import (
	"context"

	"king-starter/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

type Repository struct {
	*gormutil.BaseRepo[CoreUser]
}

// NewRepository 创建 User Repo
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		BaseRepo: gormutil.NewBaseRepo[CoreUser](db),
	}
}

// GetByUsername 根据用户名查询用户（用于登录校验等）
func (r *Repository) GetByUsername(ctx context.Context, username string) (*CoreUser, error) {
	var user CoreUser
	err := r.GetDB(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePassword 更新密码
func (r *Repository) UpdatePassword(ctx context.Context, userID, newHash string) error {
	return r.GetDB(ctx).Model(&CoreUser{}).Where("id = ?", userID).Update("password", newHash).Error
}

// UpdateStatus 更新状态
func (r *Repository) UpdateStatus(ctx context.Context, userID string, status int) error {
	return r.GetDB(ctx).Model(&CoreUser{}).Where("id = ?", userID).Update("status", status).Error
}
