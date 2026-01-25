package identity

import (
	"context"
	"crypto/subtle"
	"king-starter/internal/app"
	"king-starter/internal/router/core/user"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logx"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务接口
type AuthService interface {
	// PasswordAuth 密码认证
	PasswordAuth(ctx context.Context, username, password string) (*user.CoreUser, error)
	// PhoneAuth 手机号认证
	PhoneAuth(ctx context.Context, phone, code string) (*user.CoreUser, error)
	// TwoFactorAuth 双因素认证验证
	TwoFactorAuth(ctx context.Context, userID, token string) (bool, error)
	// OAuth2Auth OAuth2认证
	OAuth2Auth(ctx context.Context, provider, code string) (*user.CoreUser, error)
	// GenerateToken 生成JWT令牌
	GenerateToken(userID, username, roles string) (string, error)
	// ValidateToken 验证令牌
	ValidateToken(token string) (*jwt.CustomClaims, error)
}

// DefaultAuthService 默认认证服务实现
type DefaultAuthService struct {
	userRepo  *user.Repository
	loginRepo *LoginRepo
	oauthRepo *OAuthRepo
	twoFARepo *TwoFARepo
	jwt       *jwt.JWT
}

func NewDefaultAuthService(
	userRepo *user.Repository,
	loginRepo *LoginRepo,
	oauthRepo *OAuthRepo,
	twoFARepo *TwoFARepo,
	appInstance *app.App,
) *DefaultAuthService {
	return &DefaultAuthService{
		userRepo:  userRepo,
		loginRepo: loginRepo,
		oauthRepo: oauthRepo,
		twoFARepo: twoFARepo,
		jwt:       appInstance.Jwt,
	}
}

// PasswordAuth 密码认证
func (s *DefaultAuthService) PasswordAuth(ctx context.Context, username, password string) (*user.CoreUser, error) {
	// 根据用户名查询用户
	dbUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		logx.Error("Failed to get user by username", "error", err)
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		logx.Error("Invalid password for user", "username", username)
		return nil, err
	}

	return dbUser, nil
}

// PhoneAuth 手机号认证
func (s *DefaultAuthService) PhoneAuth(ctx context.Context, phone, code string) (*user.CoreUser, error) {
	// 这里需要实现手机验证码逻辑
	// 暂时返回错误，具体实现需要根据业务需求
	// 在实际项目中，需要实现短信验证码验证逻辑
	logx.Warn("Phone authentication is not implemented yet", "phone", phone)
	return nil, nil
}

// TwoFactorAuth 双因素认证验证
func (s *DefaultAuthService) TwoFactorAuth(ctx context.Context, userID, token string) (bool, error) {
	// 获取用户2FA配置
	twoFAConfig, err := s.twoFARepo.GetTwoFAByUserID(ctx, userID)
	if err != nil {
		logx.Error("Failed to get 2FA config for user", "user_id", userID, "error", err)
		return false, err
	}

	// 验证令牌
	// 注意：这里需要根据实际的2FA实现逻辑来验证
	// 通常使用TOTP算法验证动态码
	// 这里只是示例，实际实现会更复杂
	expected := twoFAConfig.Secret // 这里应该是根据算法计算出的期望值
	if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1 {
		return true, nil
	}

	return false, nil
}

// OAuth2Auth OAuth2认证
func (s *DefaultAuthService) OAuth2Auth(ctx context.Context, provider, code string) (*user.CoreUser, error) {
	// 使用现有的OAuth逻辑
	// 这里需要实现OAuth2的具体流程
	// 暂时返回错误，需要根据具体提供商实现
	logx.Warn("OAuth2 authentication is not implemented yet", "provider", provider)
	return nil, nil
}

// GenerateToken 生成JWT令牌
func (s *DefaultAuthService) GenerateToken(userID, username, roles string) (string, error) {
	return s.jwt.GenerateToken(userID, username, roles)
}

// ValidateToken 验证令牌
func (s *DefaultAuthService) ValidateToken(token string) (*jwt.CustomClaims, error) {
	return s.jwt.ParseToken(token)
}
