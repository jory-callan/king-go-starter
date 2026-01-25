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
