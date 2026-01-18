package resp

// Response 统一返回结构
type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type RespOption func(*Response)

func WithMsg(msg string) RespOption {
	return func(r *Response) { r.Msg = msg }
}
func WithCode(code string) RespOption {
	return func(r *Response) { r.Code = code }
}
func WithData(data interface{}) RespOption {
	return func(r *Response) { r.Data = data }
}

// New 调用方式
//
//	response.New(response.OK)
//	response.New(response.OK, response.WithData(user))
//	response.New(response.ErrUserNotFound)
//	response.New(response.ErrUserNotFound, response.WithMsg("用户不存在，请检查ID"))
//	response.New(response.OK, response.WithMsg("操作完成"), response.WithData(result))
func New(codeMsg CodeMsg, opts ...RespOption) Response {
	resp := Response{
		Code: codeMsg.Code,
		Msg:  codeMsg.Msg, // 默认用 CodeMsg 的 Msg
		Data: nil,         // 默认无 data
	}
	for _, opt := range opts {
		opt(&resp)
	}
	return resp
}

// 快捷函数

func Success(opts ...RespOption) Response                { return New(OK, opts...) }
func Error(codeMsg CodeMsg, opts ...RespOption) Response { return New(codeMsg, opts...) }
