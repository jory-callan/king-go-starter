package page

// Page 分页结果，兼容 MyBatis-Plus 风格
type Page[T any] struct {
	Total   int64 `json:"total"`   // 总记录数
	Pages   int64 `json:"pages"`   // 总页数（自动计算）
	Size    int64 `json:"size"`    // 每页大小
	Page    int64 `json:"page"`    // 当前页码（从 1 开始）
	Records []T   `json:"records"` // 当前页数据列表
}

// New 创建分页结果
//
// 参数说明：
//   - total: 总记录数（从数据库 COUNT 得到）
//   - size: 每页条数（如 10, 20）默认 10，最大 100
//   - page: 当前页码（从 1 开始，前端传入）
//   - records: 当前页实际数据（切片）
//
// 自动计算 pages = ceil(total / size)
func New[T any](total, size, page int64, records []T) Page[T] {
	if size <= 0 {
		size = 10 // 默认每页 10 条
	}
	if size > 100 {
		size = 100 // 最大每页 100 条
	}
	if page <= 0 {
		page = 1
	}

	pages := total / size
	if total%size != 0 {
		pages++
	}
	if pages == 0 {
		pages = 1
	}

	return Page[T]{
		Total:   total,
		Pages:   pages,
		Size:    size,
		Page:    page,
		Records: records,
	}
}
