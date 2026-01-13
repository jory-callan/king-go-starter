package jwt

import (
	"king-starter/config"
	"king-starter/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义声明（按需扩展）
type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username,omitempty"`
	Roles    string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// JWT 极简封装
type JWT struct {
	config *config.JwtConfig
	logger *logger.Logger
}

// New 创建实例
func New(cfg config.JwtConfig, logger *logger.Logger) *JWT {
	return &JWT{config: &cfg, logger: logger}
}

// NewWithDefault 创建实例（默认日志）
func NewWithDefault(logger *logger.Logger) *JWT {
	return New(config.DefaultJwtConfig(), logger)
}

// Generate 生成Token（简化版）
func (j *JWT) Generate(userID, username string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(j.config.Expires))),
			Issuer:    j.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.Secret))
}

// Parse 解析Token（简化版）
func (j *JWT) Parse(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

// 可选：常用快捷方法
func (j *JWT) GenerateWithData(userID string, data map[string]interface{}) (string, error) {
	// 如果需要支持更多数据，可以扩展
	return j.Generate(userID, "")
}

func (j *JWT) IsValid(tokenString string) bool {
	_, err := j.Parse(tokenString)
	return err == nil
}
