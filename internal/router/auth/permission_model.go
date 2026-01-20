package auth

import (
	"time"
)

// Permission 权限模型
type Permission struct {
	ID        string    `gorm:"type:varchar(32);primaryKey;comment:主键ID (UUID v7)" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null;comment:权限名称" json:"name"`
	Code      string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:权限编码" json:"code"`
	Resource  string    `gorm:"type:varchar(255);not null;comment:资源路径" json:"resource"`
	Method    string    `gorm:"type:varchar(10);not null;comment:HTTP方法" json:"method"`
	Status    int       `gorm:"type:tinyint(4);default:1;index;comment:状态 1:正常 0:禁用" json:"status"`
	CreatedAt time.Time `gorm:"index;autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time `gorm:"index;autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt time.Time `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string    `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by"`
}

// TableName 设置表名
func (Permission) TableName() string {
	return "permissions"
}
