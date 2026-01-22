package identity

import (
	"time"
)

// LoginLog 登录日志模型
type LoginLog struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Username  string    `gorm:"type:varchar(50)" json:"username"`
	IP        string    `gorm:"type:varchar(50)" json:"ip"`
	UserAgent string    `gorm:"type:varchar(255)" json:"user_agent"`
	Status    int       `gorm:"type:tinyint;default:0" json:"status"` // 0: 失败, 1: 成功
	Message   string    `gorm:"type:varchar(255)" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "core_login_logs"
}

// RefreshToken 刷新令牌模型
type RefreshToken struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (RefreshToken) TableName() string {
	return "core_refresh_tokens"
}
