package identity

import (
	"context"
	"king-starter/internal/router/core/access"
	"king-starter/internal/router/core/user"
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

type RegisterRepo struct {
	userRepo  *user.Repository
	roleRepo  *access.RoleRepo
	userRoleRepo *gormutil.BaseRepo[access.CoreUserRole]
}

func NewRegisterRepo(db *gorm.DB) *RegisterRepo {
	return &RegisterRepo{
		userRepo:  user.NewRepository(db),
		roleRepo:  access.NewRoleRepo(db),
		userRoleRepo: gormutil.NewBaseRepo[access.CoreUserRole](db),
	}
}

// AssignRoleToUser 分配角色给用户
func (r *RegisterRepo) AssignRoleToUser(ctx context.Context, userID, roleID string, operatorID string) error {
	userRole := &access.CoreUserRole{
		UserID: userID,
		RoleID: roleID,
		CreatedBy: operatorID,
	}
	return r.userRoleRepo.Create(ctx, userRole)
}

// GetUserRole 获取用户角色
func (r *RegisterRepo) GetUserRole(ctx context.Context, userID string) ([]*access.CoreUserRole, error) {
	var userRoles []*access.CoreUserRole
	err := r.userRoleRepo.DB.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).Find(&userRoles).Error
	return userRoles, err
}

// RemoveUserRole 移除用户角色
func (r *RegisterRepo) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	return r.userRoleRepo.DB.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&access.CoreUserRole{}).Error
}