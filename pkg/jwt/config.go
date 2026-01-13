package jwt

import (
	"fmt"
	"time"
)

// JwtConfig 简化配置
type JwtConfig struct {
	Secret string `yaml:"secret"` // 必填
	Expire int    `yaml:"expire"` // 过期时间（秒）
	Issuer string `yaml:"issuer"` // 签发者
}

// Validate 简单验证
func (c *JwtConfig) Validate() error {
	if c.Secret == "" {
		return fmt.Errorf("jwt secret is required")
	}
	return nil
}

func DefaultJwtConfig() JwtConfig {
	return JwtConfig{
		Secret: "your-32-char-secret-key-here",
		Expire: int(time.Hour * 24), // 默认24小时
		Issuer: "sre-server",
	}
}

/*
jwt:
  secret: "your-32-char-secret-key-here"
  expires: 86400    # 24小时过期
  issuer: "myapp"
*/
