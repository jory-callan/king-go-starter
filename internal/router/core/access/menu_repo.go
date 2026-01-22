package access

import (
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

type MenuRepo struct {
	*gormutil.BaseRepo[Menu]
}

func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{BaseRepo: gormutil.NewBaseRepo[Menu](db)}
}
