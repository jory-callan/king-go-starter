package gormutil

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	// 主键 UUID V7 长度 36
	Id        string `gorm:"primaryKey;type:varchar(36)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
