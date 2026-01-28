package user

import (
	"time"

	"gorm.io/gorm"
)

type CoreUser struct {
	ID        string         `gorm:"type:varchar(32);primaryKey;comment:用户ID(UUID v7)" json:"id"`
	Username  string         `gorm:"type:varchar(50);not null;comment:用户名" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null;comment:密码哈希" json:"-"` // 序列化时忽略
	Nickname  string         `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Status    int            `gorm:"type:tinyint;default:1;comment:状态(1:正常 0:禁用)" json:"status"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null;comment:邮箱" json:"email"`
	Phone     string         `gorm:"type:varchar(20);uniqueIndex;not null;comment:手机号" json:"phone"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`
	DeletedBy string         `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by,omitempty"`
}

func (u *CoreUser) TableName() string {
	return "core_user"
}
