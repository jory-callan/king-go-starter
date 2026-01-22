package identity

import (
	"time"
)

// TwoFA 2FA 配置模型
type TwoFA struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);uniqueIndex" json:"user_id"`
	Secret    string    `gorm:"type:varchar(255)" json:"secret"`
	Status    int       `gorm:"type:tinyint;default:0" json:"status"` // 0: 禁用, 1: 启用
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (TwoFA) TableName() string {
	return "core_twofa"
}

// TwoFALog 2FA 验证日志模型
type TwoFALog struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Status    int       `gorm:"type:tinyint;default:0" json:"status"` // 0: 失败, 1: 成功
	Message   string    `gorm:"type:varchar(255)" json:"message"`
	IP        string    `gorm:"type:varchar(50)" json:"ip"`
	UserAgent string    `gorm:"type:varchar(255)" json:"user_agent"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (TwoFALog) TableName() string {
	return "core_twofa_logs"
}
