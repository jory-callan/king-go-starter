package auth

import (
	"context"
	"errors"
	"king-starter/pkg/database/gormutil"
	"king-starter/pkg/goutils/idutil"
	"time"

	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct {
	*gormutil.BaseService[Role]
}

// NewRoleService 创建角色服务实例
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		BaseService: gormutil.NewBaseService[Role](db),
	}
}

// Create 创建角色
func (s *RoleService) Create(ctx context.Context, role *Role, operatorID string) error {
	// 生成UUID v7并去除连字符
	role.ID = idutil.ShortUUIDv7()
	role.CreatedBy = operatorID
	role.UpdatedBy = operatorID
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()
	
	return s.BaseService.Create(ctx, role)
}

// Update 更新角色
func (s *RoleService) Update(ctx context.Context, role *Role, operatorID string) error {
	// 验证角色是否存在
	existing, err := s.GetByID(ctx, role.ID)
	if err != nil {
		return err
	}
	
	// 更新操作人信息
	role.UpdatedBy = operatorID
	role.UpdatedAt = time.Now()
	
	// 保留原有创建信息
	role.CreatedAt = existing.CreatedAt
	role.CreatedBy = existing.CreatedBy
	
	return s.BaseService.Update(ctx, role)
}

// Delete 删除角色
func (s *RoleService) Delete(ctx context.Context, id string, operatorID string) error {
	// 软删除，更新删除人信息
	return s.GetDB(ctx).Model(&Role{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": operatorID,
	}).Error
}

// ListByCondition 条件查询角色列表
func (s *RoleService) ListByCondition(ctx context.Context, conds map[string]interface{}, page, pageSize int) ([]Role, int64, error) {
	// 构建查询条件
	query := s.GetDB(ctx).Model(&Role{}).Where("deleted_at IS NULL")
	
	// 添加条件过滤
	if name, ok := conds["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code, ok := conds["code"].(string); ok && code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if status, ok := conds["status"].(int); ok {
		query = query.Where("status = ?", status)
	}
	
	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	var roles []Role
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	
	return roles, total, nil
}

// GetByCode 根据角色编码获取角色
func (s *RoleService) GetByCode(ctx context.Context, code string) (*Role, error) {
	var role Role
	err := s.GetDB(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}
