package auth_2fa

import (
	"time"
)

// TwoFAConfig 2FA 配置模型
type TwoFAConfig struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);uniqueIndex" json:"user_id"`
	Secret    string    `gorm:"type:varchar(255)" json:"secret"`
	Status    int       `gorm:"type:tinyint;default:0" json:"status"` // 0: 禁用, 1: 启用
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (TwoFAConfig) TableName() string {
	return "core_user_twofa"
}
