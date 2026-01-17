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
)

// // 所有错误码定义（集中管理）
// var (
// 	SUCCESS            = CodeMessage{"0", "操作成功"}
// 	UNKNOWN_ERROR      = CodeMessage{"-1", "未知错误"}
// 	INVALID_PARAM      = CodeMessage{"400", "参数无效"}
// 	UNAUTHORIZED       = CodeMessage{"401", "未登录或登录已过期"}
// 	FORBIDDEN          = CodeMessage{"403", "权限不足"}
// 	NOT_FOUND          = CodeMessage{"404", "资源不存在"}
// 	METHOD_NOT_ALLOWED = CodeMessage{"405", "方法不支持"}
// 	SERVER_ERROR       = CodeMessage{"500", "服务器内部错误"}
// 	DB_ERROR           = CodeMessage{"501", "数据库操作失败"}
// 	VALIDATE_ERROR     = CodeMessage{"422", "数据验证失败"}
// )

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
