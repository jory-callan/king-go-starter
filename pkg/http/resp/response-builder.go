package resp

import (
	"github.com/labstack/echo/v4"
)

// === 2. ResponseBuilder（高级构建，可扩展）===

// ResponseBuilder 构建器，用于复杂场景
type ResponseBuilder struct {
	response *Response
}

// NewBuilder 创建新的构建器
func NewBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		response: &Response{
			Code:    SUCCESS.Code,
			Message: SUCCESS.Message,
		},
	}
}

// Success 设置为成功状态
func (rb *ResponseBuilder) Success() *ResponseBuilder {
	rb.response.WithCode(SUCCESS.Code).WithMessage(SUCCESS.Message)
	return rb
}

// ErrorJSON 设置为错误状态
func (rb *ResponseBuilder) Error(cm CodeMessage) *ResponseBuilder {
	rb.response.WithCode(cm.Code).WithMessage(cm.Message).WithData(nil)
	return rb
}

// Data 设置数据
func (rb *ResponseBuilder) Data(data interface{}) *ResponseBuilder {
	rb.response.WithData(data)
	return rb
}

// Message 自定义消息
func (rb *ResponseBuilder) Message(msg string) *ResponseBuilder {
	rb.response.WithMessage(msg)
	return rb
}

// Code 自定义码
func (rb *ResponseBuilder) Code(code string) *ResponseBuilder {
	rb.response.WithCode(code)
	return rb
}

// Page 分页数据
func (rb *ResponseBuilder) Page(total int64, page, size int, rows interface{}) *ResponseBuilder {
	rb.Success()
	rb.response.Data = ResponsePagination{
		Total: total,
		Page:  page,
		Size:  size,
		Rows:  rows,
	}
	return rb
}

// Send 发送响应
func (rb *ResponseBuilder) Send(c echo.Context) {
	rb.response.Send(c)
}
