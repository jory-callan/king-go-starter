package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
}

func (User) TableName() string {
	return "users"
}
