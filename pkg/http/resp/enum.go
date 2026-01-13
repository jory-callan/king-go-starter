package resp

// CodeMessage 模拟枚举：错误码 + 消息
type CodeMessage struct {
	Code    string
	Message string
}

func (cm CodeMessage) String() string {
	return "Code: " + cm.Code + ", Message: " + cm.Message
}

// 所有错误码定义（集中管理）
var (
	SUCCESS            = CodeMessage{"0", "操作成功"}
	UNKNOWN_ERROR      = CodeMessage{"-1", "未知错误"}
	INVALID_PARAM      = CodeMessage{"400", "参数无效"}
	UNAUTHORIZED       = CodeMessage{"401", "未登录或登录已过期"}
	FORBIDDEN          = CodeMessage{"403", "权限不足"}
	NOT_FOUND          = CodeMessage{"404", "资源不存在"}
	METHOD_NOT_ALLOWED = CodeMessage{"405", "方法不支持"}
	SERVER_ERROR       = CodeMessage{"500", "服务器内部错误"}
	DB_ERROR           = CodeMessage{"501", "数据库操作失败"}
	VALIDATE_ERROR     = CodeMessage{"422", "数据验证失败"}
)
