package gormutil

import (
	"context"

	"gorm.io/gorm"
)

// BaseRepo 泛型
// T: 实体模型类型
// ID: 主键类型 (如 uint, int, string)
type BaseRepo[T any] struct {
	DB *gorm.DB
}

// NewBaseRepo 构造函数
func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{DB: db}
}

// Create 插入单条记录
func (r *BaseRepo[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

// GetByID 根据主键查询单条
func (r *BaseRepo[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, id).Error
	return &entity, err
}

// Update 根据主键更新 (会保存所有字段，包括零值)
func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

// Delete 根据主键删除
func (r *BaseRepo[T]) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(new(T), id).Error
}

// CreateBatch 批量插入
// entities: 实体切片
// batchSize: 每次插入的数量，例如 100
func (r *BaseRepo[T]) CreateBatch(ctx context.Context, entities []*T, batchSize int) error {
	if len(entities) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).CreateInBatches(entities, batchSize).Error
}

// DeleteBatch 根据主键切片批量删除
func (r *BaseRepo[T]) DeleteBatch(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Delete(new(T), ids).Error
}

// List 分页查询 (基础版，不带条件)
// page: 页码 (从 1 开始)
// pageSize: 每页数量
// 返回: 数据列表, 总数, 错误
func (r *BaseRepo[T]) List(ctx context.Context, page, pageSize int) ([]T, int64, error) {
	return r.ListByCondition(ctx, nil, page, pageSize)
}

// ListByCondition 根据条件分页查询 (增强版)
// conds: 查询条件，可以是 map[string]interface{} 或 string (如 "name = ?", "tom")
// page: 页码
// pageSize: 每页数量
func (r *BaseRepo[T]) ListByCondition(ctx context.Context, conds interface{}, page, pageSize int) ([]T, int64, error) {
	var results []T
	var total int64

	query := r.DB.WithContext(ctx).Model(new(T))

	// 如果有条件，应用条件
	if conds != nil {
		query = query.Where(conds)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	// 只有当 total > 0 时才去查列表，稍微优化一点性能
	if total > 0 {
		offset := (page - 1) * pageSize
		// 防止 pageSize 为负数或过大导致报错
		if offset < 0 {
			offset = 0
		}
		err := query.Offset(offset).Limit(pageSize).Find(&results).Error
		return results, total, err
	}

	return []T{}, total, nil
}

// GetDB 获取原始的 gorm.DB 实例
// 用法: repo.GetDB().Where("name like ?", "%jack%").Find(&users)
func (r *BaseRepo[T]) GetDB() *gorm.DB {
	return r.DB
}
