package access

import (
	"king-starter/pkg/database/gormutil"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	*gormutil.BaseRepo[CorePermission]
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{BaseRepo: gormutil.NewBaseRepo[CorePermission](db)}
}
