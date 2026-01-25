package role

import (
	"time"

	"gorm.io/gorm"
)

// CoreRole 角色
type CoreRole struct {
	ID        string         `gorm:"type:varchar(32);primaryKey;comment:角色ID" json:"id"`
	Code      string         `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色编码" json:"code"`
	Name      string         `gorm:"type:varchar(50);not null;comment:角色名称" json:"name"`
	Status    int            `gorm:"type:tinyint;default:1;comment:状态" json:"status"`
	Remark    string         `gorm:"type:varchar(255);comment:备注" json:"remark"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`
	DeletedBy string         `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by,omitempty"`
}

func (CoreRole) TableName() string {
	return "core_role"
}

// CoreUserRole 用户角色关联表
type CoreUserRole struct {
	UserID    string         `gorm:"type:varchar(32);primaryKey;comment:用户ID" json:"user_id"`
	RoleID    string         `gorm:"type:varchar(32);primaryKey;comment:角色ID" json:"role_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`
}

func (CoreUserRole) TableName() string {
	return "core_user_roles"
}

// CoreRoleMenu 角色菜单关联表 (手动维护)
type CoreRoleMenu struct {
	RoleID string `gorm:"type:varchar(32);primaryKey;comment:角色ID" json:"role_id"`
	MenuID string `gorm:"type:varchar(32);primaryKey;comment:菜单ID" json:"menu_id"`
}

func (CoreRoleMenu) TableName() string {
	return "core_role_menu"
}
