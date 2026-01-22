package access

import (
	"time"

	"gorm.io/gorm"
)

// CoreMenu 菜单资源
type CoreMenu struct {
	ID        string         `gorm:"type:varchar(32);primaryKey;comment:菜单ID" json:"id"`
	ParentID  string         `gorm:"type:varchar(32);default:0;comment:父级菜单ID" json:"parent_id"`
	Name      string         `gorm:"type:varchar(50);not null;comment:菜单名称" json:"name"`
	Path      string         `gorm:"type:varchar(200);comment:路由路径" json:"path"`
	Icon      string         `gorm:"type:varchar(50);comment:图标" json:"icon"`
	Sort      int            `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status    int            `gorm:"type:tinyint;default:1;comment:状态" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at,omitempty"`
	DeletedBy string         `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by,omitempty"`
}

func (CoreMenu) TableName() string {
	return "core_menu"
}
