package access

import (
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	*gormutil.BaseRepo[Permission]
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{BaseRepo: gormutil.NewBaseRepo[Permission](db)}
}
