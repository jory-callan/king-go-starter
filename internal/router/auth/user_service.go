package auth

import (
	"context"
	"errors"
	"king-starter/pkg/database/gormutil"
	"king-starter/pkg/goutils/idutil"
	"time"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	*gormutil.BaseService[User]
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		BaseService: gormutil.NewBaseService[User](db),
	}
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, user *User, operatorID string) error {
	// 生成UUID v7并去除连字符
	user.ID = idutil.ShortUUIDv7()
	user.CreatedBy = operatorID
	user.UpdatedBy = operatorID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.BaseService.Create(ctx, user)
}

// Update 更新用户
func (s *UserService) Update(ctx context.Context, user *User, operatorID string) error {
	// 验证用户是否存在
	existing, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	// 更新操作人信息
	user.UpdatedBy = operatorID
	user.UpdatedAt = time.Now()

	// 保留原有创建信息
	user.CreatedAt = existing.CreatedAt
	user.CreatedBy = existing.CreatedBy

	return s.BaseService.Update(ctx, user)
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, id string, operatorID string) error {
	// 软删除，更新删除人信息
	return s.GetDB(ctx).Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": operatorID,
	}).Error
}

// ListByCondition 条件查询用户列表
func (s *UserService) ListByCondition(ctx context.Context, conds map[string]interface{}, page, pageSize int) ([]User, int64, error) {
	// 构建查询条件
	query := s.GetDB(ctx).Model(&User{}).Where("deleted_at IS NULL")

	// 添加条件过滤
	if username, ok := conds["username"].(string); ok && username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if email, ok := conds["email"].(string); ok && email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if phone, ok := conds["phone"].(string); ok && phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
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
	var users []User
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := s.GetDB(ctx).Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (s *UserService) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.GetDB(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
