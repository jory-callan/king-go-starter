package permission

import (
	"time"

	"gorm.io/gorm"
)

// CorePermission 权限码 (api:xxx, menu:xxx)，支持通配符匹配
type CorePermission struct {
	ID        string         `gorm:"type:varchar(32);primaryKey;comment:权限ID" json:"id"`
	Code      string         `gorm:"type:varchar(100);uniqueIndex;not null;comment:权限码(支持通配符*)" json:"code"`
	Name      string         `gorm:"type:varchar(50);not null;comment:权限名称" json:"name"`
	Type      string         `gorm:"type:varchar(20);not null;comment:类型(menu/api)" json:"type"`
	ParentID  string         `gorm:"type:varchar(32);default:0;comment:父级权限ID" json:"parent_id"` // 支持菜单层级
	Path      string         `gorm:"type:varchar(200);comment:路由路径" json:"path"`                 // 菜单路径
	Icon      string         `gorm:"type:varchar(50);comment:图标" json:"icon"`                    // 菜单图标
	Sort      int            `gorm:"type:int;default:0;comment:排序" json:"sort"`                  // 排序
	Status    int            `gorm:"type:tinyint;default:1;comment:状态" json:"status"`            // 状态
	Remark    string         `gorm:"type:varchar(255);comment:备注" json:"remark"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`
	DeletedBy string         `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by,omitempty"`
	Children  []CorePermission `gorm:"-" json:"children,omitempty"` // 子权限，不存储到数据库
}

func (CorePermission) TableName() string {
	return "core_permission"
}

// CoreRolePermission 角色权限关联表
type CoreRolePermission struct {
	RoleID       string `gorm:"type:varchar(32);primaryKey;comment:角色ID" json:"role_id"`
	PermissionID string `gorm:"type:varchar(32);primaryKey;comment:权限ID" json:"permission_id"`
}

func (CoreRolePermission) TableName() string {
	return "core_role_permission"
}
