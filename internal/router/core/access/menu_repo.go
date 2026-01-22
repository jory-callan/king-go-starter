package access

import (
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

type MenuRepo struct {
	*gormutil.BaseRepo[CoreMenu]
}

func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{BaseRepo: gormutil.NewBaseRepo[CoreMenu](db)}
}
