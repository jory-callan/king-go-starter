package auth

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID        string    `gorm:"type:varchar(32);primaryKey;comment:主键ID (UUID v7)" json:"id"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色名称" json:"name"`
	Code      string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色编码" json:"code"`
	Status    int       `gorm:"type:tinyint(4);default:1;index;comment:状态 1:正常 0:禁用" json:"status"`
	CreatedAt time.Time `gorm:"index;autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time `gorm:"index;autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt time.Time `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string    `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by"`
}

// TableName 设置表名
func (Role) TableName() string {
	return "roles"
}
