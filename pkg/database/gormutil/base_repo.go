package gormutil

import (
	"context"
	"king-starter/internal/response"

	"gorm.io/gorm"
)

// BaseRepo 泛型, 基础仓库层, 提供基础的 CRUD 操作
//
//	T: 实体模型类型
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

// GetDB 获取原始的 gorm.DB 实例
// 用法: repo.GetDB().Where("name like ?", "%jack%").Find(&users)
func (r *BaseRepo[T]) GetDB(ctx context.Context) *gorm.DB {
	r.DB.Where("deleted_at IS NULL")
	return r.DB.WithContext(ctx)
}

// PaginationWithScopes 分页查询 (增强版)
// ctx: 上下文
// pq: 分页查询参数 @see: response.PageQuery
// scopes: 可选的查询范围函数，例如过滤、关联预加载等
// 返回: 分页结果, 错误
func (r *BaseRepo[T]) PaginationWithScopes(
	ctx context.Context,
	pq *response.PageQuery,
	scopes ...func(*gorm.DB) *gorm.DB,
) (*response.PageResult[T], error) {

	if pq == nil {
		pq = &response.PageQuery{}
	}
	pq.Normalize()

	db := r.DB.WithContext(ctx).Model(new(T))

	// apply scopes
	for _, scope := range scopes {
		db = db.Scopes(scope)
	}

	// apply order
	for _, o := range pq.Order {
		if o.Desc {
			db = db.Order(o.Field + " DESC")
		} else {
			db = db.Order(o.Field + " ASC")
		}
	}

	// list
	var list []T

	// 如果不需要 count，为了计算 hasMore，我们需要多取一条
	limit := pq.Size
	if !pq.NeedCount {
		limit = pq.Size + 1
	}

	offset := (pq.Page - 1) * pq.Size
	if err := db.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}

	// 处理 hasMore
	hasMore := false
	if !pq.NeedCount {
		if len(list) > pq.Size {
			hasMore = true
			list = list[:pq.Size] // 截断多取的那条
		}
	}

	// build result
	res := &response.PageResult[T]{
		Items:   list,
		Page:    pq.Page,
		Size:    pq.Size,
		HasMore: hasMore,
	}

	// count if needed
	if pq.NeedCount {
		var total int64
		if err := db.Count(&total).Error; err != nil {
			return nil, err
		}
		res.Total = total
		res.HasMore = int64(pq.Page*pq.Size) < total
	}

	return res, nil
}

// Pagination 分页查询 (方便版)
// ctx: 上下文
// pq: 分页查询参数 @see: response.PageQuery
// customDB: 自定义的 gorm.DB 实例，用于添加额外的查询条件或关联预加载，只要最终能 Find(&list) 即可
// 返回: 分页结果, 错误
func (r *BaseRepo[T]) Pagination(
	ctx context.Context,
	pq *response.PageQuery,
	customDB *gorm.DB,
) (*response.PageResult[T], error) {

	if pq == nil {
		pq = &response.PageQuery{}
	}
	pq.Normalize()

	// 如果 customDB 为 nil，使用 r.db
	db := customDB
	if db == nil {
		db = r.DB
	}

	db = db.WithContext(ctx).Model(new(T))

	// apply order
	for _, o := range pq.Order {
		if o.Desc {
			db = db.Order(o.Field + " DESC")
		} else {
			db = db.Order(o.Field + " ASC")
		}
	}

	// pagination list
	limit := pq.Size
	if !pq.NeedCount {
		limit = pq.Size + 1
	}

	var list []T
	offset := (pq.Page - 1) * pq.Size
	if err := db.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}

	hasMore := false
	if !pq.NeedCount && len(list) > pq.Size {
		hasMore = true
		list = list[:pq.Size]
	}

	res := &response.PageResult[T]{
		Items:   list,
		Page:    pq.Page,
		Size:    pq.Size,
		HasMore: hasMore,
	}

	if pq.NeedCount {
		var total int64
		if err := db.Count(&total).Error; err != nil {
			return nil, err
		}
		res.Total = total
		res.HasMore = int64(pq.Page*pq.Size) < total
	}

	return res, nil
}
