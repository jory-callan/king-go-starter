package auth_password

import (
	"time"
)

// CoreLoginLog 登录日志模型
type CoreLoginLog struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID     string    `gorm:"type:varchar(36);index" json:"user_id"`
	Username   string    `gorm:"type:varchar(50)" json:"username"`
	AuthType   string    `gorm:"type:varchar(20);index" json:"auth_type"`            // 认证类型
	LoginType  string    `gorm:"type:varchar(10);index" json:"login_type"`           // 登录类型
	IP         string    `gorm:"type:varchar(50)" json:"ip"`                         // IP地址
	UserAgent  string    `gorm:"type:varchar(255)" json:"user_agent"`                // 用户代理
	DeviceInfo string    `gorm:"embedded;embeddedPrefix:device_" json:"device_info"` // 设备信息
	Country    string    `gorm:"type:varchar(50)" json:"country"`                    // 国家
	Province   string    `gorm:"type:varchar(50)" json:"province"`                   // 省份
	City       string    `gorm:"type:varchar(50)" json:"city"`                       // 城市
	Message    string    `gorm:"type:varchar(255)" json:"message"`                   // 描述信息
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (CoreLoginLog) TableName() string {
	return "core_login_logs"
}

// CoreRefreshToken 刷新令牌模型
type CoreRefreshToken struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (CoreRefreshToken) TableName() string {
	return "core_refresh_tokens"
}
