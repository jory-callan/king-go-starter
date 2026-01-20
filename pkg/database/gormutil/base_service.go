package gormutil

import (
	"gorm.io/gorm"
)

// BaseService 基础服务层, 提供基础的 CRUD 操作
//
//	T: 实体模型类型
type BaseService[T any] struct {
	*BaseRepo[T] // 嵌入 BaseRepo 获取所有方法
}

func NewBaseService[T any](db *gorm.DB) *BaseService[T] {
	return &BaseService[T]{
		BaseRepo: NewBaseRepo[T](db),
	}
}
