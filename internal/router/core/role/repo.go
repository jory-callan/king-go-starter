package role

import (
	"context"

	"king-starter/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

type RoleRepo struct {
	*gormutil.BaseRepo[CoreRole]
	userRoleRepo *gormutil.BaseRepo[CoreUserRole]
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		BaseRepo:     gormutil.NewBaseRepo[CoreRole](db),
		userRoleRepo: gormutil.NewBaseRepo[CoreUserRole](db),
	}
}

// AssignRolesToUser 为用户分配多个角色
func (r *RoleRepo) AssignRolesToUser(ctx context.Context, userID string, roleIDs []string, operatorID string) error {
	db := r.userRoleRepo.DB.WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		// 先删除用户原有的角色
		if err := tx.Where("user_id = ?", userID).Delete(&CoreUserRole{}).Error; err != nil {
			return err
		}

		// 添加新角色
		if len(roleIDs) > 0 {
			var userRoles []CoreUserRole
			for _, roleID := range roleIDs {
				userRoles = append(userRoles, CoreUserRole{
					UserID:    userID,
					RoleID:    roleID,
					CreatedBy: operatorID,
				})
			}
			return tx.Create(&userRoles).Error
		}
		return nil
	})
}

// GetUserRoles 获取用户的角色ID列表
func (r *RoleRepo) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var roleIDs []string
	err := r.userRoleRepo.DB.WithContext(ctx).
		Model(&CoreUserRole{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Pluck("role_id", &roleIDs).Error
	return roleIDs, err
}

// GetUserRolesWithDetails 获取用户的角色详细信息
func (r *RoleRepo) GetUserRolesWithDetails(ctx context.Context, userID string) ([]CoreRole, error) {
	var roles []CoreRole
	err := r.GetDB(ctx).
		Joins("JOIN core_user_roles ON core_role.id = core_user_roles.role_id").
		Where("core_user_roles.user_id = ? AND core_user_roles.deleted_at IS NULL", userID).
		Find(&roles).Error
	return roles, err
}

// GetRoleUsers 获取角色下的用户ID列表
func (r *RoleRepo) GetRoleUsers(ctx context.Context, roleID string) ([]string, error) {
	var userIDs []string
	err := r.userRoleRepo.DB.WithContext(ctx).
		Model(&CoreUserRole{}).
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// RemoveUserRole 解绑用户与角色
func (r *RoleRepo) RemoveUserRole(ctx context.Context, userID, roleID string) error {
	return r.userRoleRepo.DB.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&CoreUserRole{}).Error
}

// BatchRemoveUserRoles 批量解绑用户角色
func (r *RoleRepo) BatchRemoveUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	if len(roleIDs) == 0 {
		return nil
	}
	return r.userRoleRepo.DB.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&CoreUserRole{}).Error
}
