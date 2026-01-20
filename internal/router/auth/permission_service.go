package auth

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"king-starter/pkg/database/gormutil"
	"king-starter/pkg/goutils/idutil"
	"time"
)

// PermissionService 权限服务
type PermissionService struct {
	*gormutil.BaseService[Permission]
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{
		BaseService: gormutil.NewBaseService[Permission](db),
	}
}

// Create 创建权限
func (s *PermissionService) Create(ctx context.Context, permission *Permission, operatorID string) error {
	// 生成UUID v7并去除连字符
	permission.ID = idutil.ShortUUIDv7()
	permission.CreatedBy = operatorID
	permission.UpdatedBy = operatorID
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()
	
	return s.BaseService.Create(ctx, permission)
}

// Update 更新权限
func (s *PermissionService) Update(ctx context.Context, permission *Permission, operatorID string) error {
	// 验证权限是否存在
	existing, err := s.GetByID(ctx, permission.ID)
	if err != nil {
		return err
	}
	
	// 更新操作人信息
	permission.UpdatedBy = operatorID
	permission.UpdatedAt = time.Now()
	
	// 保留原有创建信息
	permission.CreatedAt = existing.CreatedAt
	permission.CreatedBy = existing.CreatedBy
	
	return s.BaseService.Update(ctx, permission)
}

// Delete 删除权限
func (s *PermissionService) Delete(ctx context.Context, id string, operatorID string) error {
	// 软删除，更新删除人信息
	return s.GetDB(ctx).Model(&Permission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": operatorID,
	}).Error
}

// ListByCondition 条件查询权限列表
func (s *PermissionService) ListByCondition(ctx context.Context, conds map[string]interface{}, page, pageSize int) ([]Permission, int64, error) {
	// 构建查询条件
	query := s.GetDB(ctx).Model(&Permission{}).Where("deleted_at IS NULL")
	
	// 添加条件过滤
	if name, ok := conds["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code, ok := conds["code"].(string); ok && code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if resource, ok := conds["resource"].(string); ok && resource != "" {
		query = query.Where("resource LIKE ?", "%"+resource+"%")
	}
	if method, ok := conds["method"].(string); ok && method != "" {
		query = query.Where("method = ?", method)
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
	var permissions []Permission
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&permissions).Error; err != nil {
		return nil, 0, err
	}
	
	return permissions, total, nil
}

// GetByCode 根据权限编码获取权限
func (s *PermissionService) GetByCode(ctx context.Context, code string) (*Permission, error) {
	var permission Permission
	err := s.GetDB(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}
