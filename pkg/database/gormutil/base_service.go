package gormutil

import (
	"gorm.io/gorm"
)

type BaseService[T any] struct {
	*BaseRepo[T] // 嵌入 BaseRepo 获取所有方法
}

func NewBaseService[T any](db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{
		BaseRepo: NewBaseRepo[T](db),
	}
}
