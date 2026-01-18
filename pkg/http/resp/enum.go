package resp

// 所有错误码定义（集中管理）

// 成功
var OK = CodeMsg{Code: "0", Msg: "success", Status: StatusSuccess}

// 通用错误
var (
	// 客户端错误
	ErrBadRequest   = CodeMsg{Code: "400", Msg: "参数无效", Status: StatusClientError}
	ErrUnauthorized = CodeMsg{Code: "401", Msg: "未登录或登录已过期", Status: StatusClientError}
	ErrForbidden    = CodeMsg{Code: "403", Msg: "权限不足", Status: StatusClientError}
	ErrNotFound     = CodeMsg{Code: "404", Msg: "资源不存在", Status: StatusClientError}
	// 服务端错误
	ErrInternalServer = CodeMsg{Code: "500", Msg: "内部服务错误", Status: StatusServerError}
	ErrUnknown        = CodeMsg{Code: "500", Msg: "未知错误", Status: StatusServerError}
)

// CodeMsg 模拟枚举：错误码 + 消息 + 状态
// 状态 0 成功 1 客户端失败 2 服务端失败
type CodeMsg struct {
	Code   string
	Msg    string
	Status int
}

const (
	// 成功
	StatusSuccess = 0
	// 客户端错误
	StatusClientError = 1
	// 服务端错误
	StatusServerError = 2
)
