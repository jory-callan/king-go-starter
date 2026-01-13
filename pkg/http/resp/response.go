package resp

import (
	"github.com/labstack/echo/v4"
)

// Response 统一返回结构
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

/* === 1. Response 的 With 方法（基础链式调用）=== */
func (r *Response) Success() *Response {
	r.Code = SUCCESS.Code
	return r
}

func (r *Response) Error(cm CodeMessage) *Response {
	r.Code = cm.Code
	return r
}

// WithCode 设置响应码
func (r *Response) WithCode(code string) *Response {
	r.Code = code
	return r
}

// WithMessage 设置响应消息
func (r *Response) WithMessage(message string) *Response {
	r.Message = message
	return r
}

// WithData 设置响应数据
func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

// Send 发送响应（HTTP 状态码根据 Code 动态判断）
func (r *Response) Send(c echo.Context) {
	status := r.statusCode()
	c.JSON(status, r)
}

// statusCode 根据 Code 映射 HTTP 状态码
func (r *Response) statusCode() int {
	switch r.Code {
	case SUCCESS.Code:
		return 200
	case UNAUTHORIZED.Code:
		return 401
	case FORBIDDEN.Code:
		return 403
	case NOT_FOUND.Code:
		return 404
	case INVALID_PARAM.Code, VALIDATE_ERROR.Code:
		return 422
	case SERVER_ERROR.Code, DB_ERROR.Code:
		return 500
	default:
		return 200
	}
}
