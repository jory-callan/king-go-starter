package auth

import (
	"time"
)

// User 用户模型
type User struct {
	ID        string    `gorm:"type:varchar(32);primaryKey;comment:主键ID (UUID v7)" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:邮箱" json:"email"`
	Phone     string    `gorm:"type:varchar(20);comment:手机号" json:"phone"`
	Status    int       `gorm:"type:tinyint(4);default:1;index;comment:状态 1:正常 0:禁用" json:"status"`
	CreatedAt time.Time `gorm:"index;autoCreateTime;comment:创建时间" json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(32);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time `gorm:"index;autoUpdateTime;comment:更新时间" json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(32);comment:更新人ID" json:"updated_by"`
	DeletedAt time.Time `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string    `gorm:"type:varchar(32);comment:删除人ID" json:"deleted_by"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
