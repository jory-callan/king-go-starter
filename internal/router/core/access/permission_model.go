package access

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限码 (api:xxx, menu:xxx)
type Permission struct {
	ID        string         `gorm:"type:varchar(32);primaryKey;comment:权限ID" json:"id"`
	Code      string         `gorm:"type:varchar(100);uniqueIndex;not null;comment:权限码" json:"code"`
	Name      string         `gorm:"type:varchar(50);not null;comment:权限名称" json:"name"`
	Type      string         `gorm:"type:varchar(20);not null;comment:类型(menu/api)" json:"type"`
	Remark    string         `gorm:"type:varchar(255);comment:备注" json:"remark"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`
	DeletedBy string         `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by,omitempty"`
}

func (Permission) TableName() string {
	return "sys_permission"
}
