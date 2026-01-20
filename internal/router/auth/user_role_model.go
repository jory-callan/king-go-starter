package auth

import (
	"time"
)

// UserRole 用户角色关联表
type UserRole struct {
	ID        string    `gorm:"type:varchar(32);primaryKey;comment:主键ID (UUID v7)" json:"id"`
	UserID    string    `gorm:"type:varchar(32);index;not null;comment:用户ID" json:"user_id"`
	RoleID    string    `gorm:"type:varchar(32);index;not null;comment:角色ID" json:"role_id"`
	CreatedAt time.Time `gorm:"index;autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time `gorm:"index;autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt time.Time `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string    `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by"`
}

// TableName 设置表名
func (UserRole) TableName() string {
	return "user_roles"
}
