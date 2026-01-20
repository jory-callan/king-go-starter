package auth

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"king-starter/pkg/jwt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db          *gorm.DB
	jwt         *jwt.JWT
	userService *UserService
	roleService *RoleService
	permService *PermissionService
}

// NewAuthService 创建认证服务实例
func NewAuthService(db *gorm.DB, jwt *jwt.JWT) *AuthService {
	return &AuthService{
		db:          db,
		jwt:         jwt,
		userService: NewUserService(db),
		roleService: NewRoleService(db),
		permService: NewPermissionService(db),
	}
}

// hashPassword 密码哈希
func (s *AuthService) hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, username, password string) (string, *User, error) {
	// 查询用户
	user, err := s.userService.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("用户不存在")
	}

	// 验证用户状态
	if user.Status != 1 {
		return "", nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if user.Password != s.hashPassword(password) {
		return "", nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := s.jwt.GenerateToken(user.ID, user.Username, strconv.Itoa(s.jwt.Expire))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// VerifyPermission 验证用户是否具有指定权限码
func (s *AuthService) VerifyPermission(ctx context.Context, userID, permissionCode string) (bool, error) {
	// 通过SQL JOIN查询用户是否具有指定权限
	var count int64
	err := s.db.Table("users u").
		Joins("JOIN user_roles ur ON u.id = ur.user_id").
		Joins("JOIN role_permissions rp ON ur.role_id = rp.role_id").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Where("u.id = ? AND p.code = ? AND u.status = 1 AND p.status = 1", userID, permissionCode).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetUserPermissions 获取用户所有权限码
func (s *AuthService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	var permissions []string
	err := s.db.Table("users u").
		Select("DISTINCT p.code").
		Joins("JOIN user_roles ur ON u.id = ur.user_id").
		Joins("JOIN role_permissions rp ON ur.role_id = rp.role_id").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Where("u.id = ? AND u.status = 1 AND p.status = 1", userID).
		Pluck("p.code", &permissions).Error

	return permissions, err
}

// GetUserRoles 获取用户所有角色
func (s *AuthService) GetUserRoles(ctx context.Context, userID string) ([]Role, error) {
	var roles []Role
	err := s.db.Table("roles r").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.status = 1", userID).
		Find(&roles).Error

	return roles, err
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// 查询用户
	user, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if user.Password != s.hashPassword(oldPassword) {
		return errors.New("旧密码错误")
	}

	// 更新密码
	user.Password = s.hashPassword(newPassword)
	user.UpdatedAt = time.Now()
	user.UpdatedBy = userID

	return s.userService.Update(ctx, user, userID)
}
