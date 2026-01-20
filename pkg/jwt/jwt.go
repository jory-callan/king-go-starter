package jwt

import (
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

type JWT struct {
	Secret []byte
	Issuer string
	Expire int
}

// New 创建 JWT 实例
func New(secret []byte, issuer string, expire int) *JWT {
	return &JWT{
		Secret: secret,
		Issuer: issuer,
		Expire: expire,
	}
}

// NewWithConfig 使用配置文件进行创建 JWT 实例
func NewWithConfig(cfg *JwtConfig) (*JWT, error) {
	secret := []byte(cfg.Secret)
	issuer := cfg.Issuer
	expire := cfg.Expire

	return New(secret, issuer, expire), cfg.Validate()
}

// GenerateToken 生成 JWT 令牌
func (j *JWT) GenerateToken(userID, username, roles string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.Expire))),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}

// ParseToken 解析并验证 JWT 令牌
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析 JWT 令牌，验证签名和claims，返回 CustomClaims 结构体
	// 如果令牌无效或claims不匹配，返回错误
	// keyFunc 用于提供密钥，这里使用预定义的密钥
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// RefreshToken 刷新 JWT 令牌（延长有效期）
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// 解析旧令牌，验证签名和claims
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	// 更新过期时间
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(j.Expire)))
	// 生成新令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名并返回新令牌
	return token.SignedString(j.Secret)
}
