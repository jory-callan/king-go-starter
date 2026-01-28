package identity

import (
	"context"

	"king-starter/internal/router/core/role"
	"king-starter/internal/router/core/user"
	"king-starter/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

type RegisterRepo struct {
	userRepo     *user.Repository
	roleRepo     *role.RoleRepo
	userRoleRepo *gormutil.BaseRepo[role.CoreUserRole]
}

func NewRegisterRepo(db *gorm.DB) *RegisterRepo {
	return &RegisterRepo{
		userRepo:     user.NewRepository(db),
		roleRepo:     role.NewRoleRepo(db),
		userRoleRepo: gormutil.NewBaseRepo[role.CoreUserRole](db),
	}
}

// AssignRoleToUser 分配角色给用户
func (r *RegisterRepo) AssignRoleToUser(ctx context.Context, userID, roleID string, operatorID string) error {
	userRole := &role.CoreUserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedBy: operatorID,
	}
	return r.userRoleRepo.Create(ctx, userRole)
}

// GetUserRole 获取用户角色
func (r *RegisterRepo) GetUserRole(ctx context.Context, userID string) ([]*role.CoreUserRole, error) {
	var userRoles []*role.CoreUserRole
	err := r.userRoleRepo.DB.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).Find(&userRoles).Error
	return userRoles, err
}

// RemoveUserRole 移除用户角色
func (r *RegisterRepo) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	return r.userRoleRepo.DB.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&role.CoreUserRole{}).Error
}
