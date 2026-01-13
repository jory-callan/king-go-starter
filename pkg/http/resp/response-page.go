package resp

// === 3. 分页结构 ===

// ResponsePagination 分页返回结构
type ResponsePagination struct {
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Rows  interface{} `json:"rows"`
}
