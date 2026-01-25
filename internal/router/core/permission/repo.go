package permission

import (
	"context"
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	*gormutil.BaseRepo[CorePermission]
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{BaseRepo: gormutil.NewBaseRepo[CorePermission](db)}
}

// GetPermissionTree 获取权限树结构
func (r *PermissionRepo) GetPermissionTree(ctx context.Context, parentID string) ([]CorePermission, error) {
	var permissions []CorePermission
	query := r.GetDB(ctx).Order("sort ASC")

	if parentID == "" || parentID == "0" {
		// 获取顶级权限
		err := query.Where("parent_id = ? OR parent_id = '' OR parent_id IS NULL", "0").Find(&permissions).Error
		return permissions, err
	} else {
		// 获取指定父级下的权限
		err := query.Where("parent_id = ?", parentID).Find(&permissions).Error
		return permissions, err
	}
}

// GetUserPermissions 获取用户拥有的权限
func (r *PermissionRepo) GetUserPermissions(ctx context.Context, userID string) ([]CorePermission, error) {
	var permissions []CorePermission
	err := r.GetDB(ctx).
		Joins("JOIN core_role_permission ON core_permission.id = core_role_permission.permission_id").
		Joins("JOIN core_user_roles ON core_role_permission.role_id = core_user_roles.role_id").
		Where("core_user_roles.user_id = ?", userID).
		Find(&permissions).Error
	return permissions, err
}

// AssignRolePermissions 分配权限给角色
func (r *PermissionRepo) AssignRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	db := r.GetDB(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&CoreRolePermission{}).Error; err != nil {
			return err
		}
		if len(permissionIDs) > 0 {
			var relations []CoreRolePermission
			for _, pid := range permissionIDs {
				relations = append(relations, CoreRolePermission{
					RoleID:       roleID,
					PermissionID: pid,
				})
			}
			return tx.Create(&relations).Error
		}
		return nil
	})
}

// GetRolePermissions 获取角色拥有的权限ID列表
func (r *PermissionRepo) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	var ids []string
	err := r.GetDB(ctx).Model(&CoreRolePermission{}).Where("role_id = ?", roleID).Pluck("permission_id", &ids).Error
	return ids, err
}

// GetRolePermissionsWithDetails 获取角色拥有的权限详细信息
func (r *PermissionRepo) GetRolePermissionsWithDetails(ctx context.Context, roleID string) ([]CorePermission, error) {
	var permissions []CorePermission
	err := r.GetDB(ctx).
		Joins("JOIN core_role_permission ON core_permission.id = core_role_permission.permission_id").
		Where("core_role_permission.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// GetUserAllPermissions 获取用户通过角色获得的所有权限详细信息
func (r *PermissionRepo) GetUserAllPermissions(ctx context.Context, userID string) ([]CorePermission, error) {
	var permissions []CorePermission
	err := r.GetDB(ctx).
		Joins("JOIN core_role_permission ON core_permission.id = core_role_permission.permission_id").
		Joins("JOIN core_user_roles ON core_role_permission.role_id = core_user_roles.role_id").
		Where("core_user_roles.user_id = ?", userID).
		Find(&permissions).Error
	return permissions, err
}

// RemoveRolePermissions 移除角色的部分或全部权限
func (r *PermissionRepo) RemoveRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	db := r.GetDB(ctx)
	
	if len(permissionIDs) == 0 {
		// 如果没有指定权限ID，则移除该角色的所有权限
		return db.Where("role_id = ?", roleID).Delete(&CoreRolePermission{}).Error
	}
	
	// 否则只移除指定的权限
	return db.Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).Delete(&CoreRolePermission{}).Error
}

// GetRolePermissionTree 获取角色拥有的权限树结构
func (r *PermissionRepo) GetRolePermissionTree(ctx context.Context, roleID string) ([]CorePermission, error) {
	permissions, err := r.GetRolePermissionsWithDetails(ctx, roleID)
	if err != nil {
		return nil, err
	}
	
	// 使用工具函数构建树形结构
	return BuildPermissionTree(permissions), nil
}
