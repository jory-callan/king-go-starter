package role

import (
	"context"
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

type RoleRepo struct {
	*gormutil.BaseRepo[CoreRole]
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{BaseRepo: gormutil.NewBaseRepo[CoreRole](db)}
}

// AssignMenus 分配菜单给角色（为了兼容性保留，实际操作权限表）
func (r *RoleRepo) AssignMenus(ctx context.Context, roleID string, menuIDs []string) error {
	db := r.GetDB(ctx)

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除旧的菜单权限关联 (type = 'menu')
		if err := tx.Where("role_id = ? AND permission_id IN (SELECT id FROM core_permission WHERE type = 'menu')", roleID).Delete(&CoreRolePermission{}).Error; err != nil {
			return err
		}

		// 2. 批量插入新关联
		if len(menuIDs) > 0 {
			var relations []CoreRolePermission
			for _, menuID := range menuIDs {
				relations = append(relations, CoreRolePermission{
					RoleID:       roleID,
					PermissionID: menuID,
				})
			}
			if err := tx.Create(&relations).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// AssignPermissions 分配权限给角色
func (r *RoleRepo) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
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

// GetRoleMenus 获取角色拥有的菜单ID列表
func (r *RoleRepo) GetRoleMenus(ctx context.Context, roleID string) ([]string, error) {
	var ids []string
	err := r.GetDB(ctx).
		Model(&CoreRolePermission{}).
		Joins("JOIN core_permission ON core_role_permission.permission_id = core_permission.id").
		Where("core_role_permission.role_id = ? AND core_permission.type = 'menu'", roleID).
		Pluck("core_role_permission.permission_id", &ids).
		Error
	return ids, err
}

// GetRolePermissions 获取角色拥有的权限ID列表
func (r *RoleRepo) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	var ids []string
	err := r.GetDB(ctx).Model(&CoreRolePermission{}).Where("role_id = ?", roleID).Pluck("permission_id", &ids).Error
	return ids, err
}
