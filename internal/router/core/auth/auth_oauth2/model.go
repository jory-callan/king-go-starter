package auth_oauth2

import (
	"time"
)

// OAuthClient OAuth 客户端模型
type OAuthClient struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ClientID     string    `gorm:"type:varchar(100);uniqueIndex" json:"client_id"`
	ClientSecret string    `gorm:"type:varchar(255)" json:"client_secret"`
	Name         string    `gorm:"type:varchar(100)" json:"name"`
	RedirectURI  string    `gorm:"type:varchar(255)" json:"redirect_uri"`
	Status       int       `gorm:"type:tinyint;default:1" json:"status"` // 1: 启用, 0: 禁用
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (OAuthClient) TableName() string {
	return "core_user_oauth_clients"
}

// OAuthCode OAuth 授权码模型
type OAuthCode struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ClientID    string    `gorm:"type:varchar(100);index" json:"client_id"`
	UserID      string    `gorm:"type:varchar(36);index" json:"user_id"`
	Code        string    `gorm:"type:varchar(255);uniqueIndex" json:"code"`
	RedirectURI string    `gorm:"type:varchar(255)" json:"redirect_uri"`
	Scope       string    `gorm:"type:varchar(255)" json:"scope"`
	ExpiresAt   time.Time `gorm:"index" json:"expires_at"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (OAuthCode) TableName() string {
	return "core_user_oauth_codes"
}

// OAuthToken OAuth 令牌模型
type OAuthToken struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ClientID     string    `gorm:"type:varchar(100);index" json:"client_id"`
	UserID       string    `gorm:"type:varchar(36);index" json:"user_id"`
	AccessToken  string    `gorm:"type:varchar(255);uniqueIndex" json:"access_token"`
	RefreshToken string    `gorm:"type:varchar(255);uniqueIndex" json:"refresh_token"`
	Scope        string    `gorm:"type:varchar(255)" json:"scope"`
	TokenType    string    `gorm:"type:varchar(50);default:'Bearer'" json:"token_type"`
	ExpiresAt    time.Time `gorm:"index" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (OAuthToken) TableName() string {
	return "core_user_oauth_tokens"
}
