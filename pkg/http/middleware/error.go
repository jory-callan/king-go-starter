package middleware

import (
	"errors"
	"net/http"

	"king-starter/pkg/logger"

	"github.com/labstack/echo/v4"
)

// EchoErrorHandler 统一错误处理
func EchoErrorHandler(log *logger.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		var he *echo.HTTPError
		// 禁用 HEAD 请求
		if c.Request().Method == http.MethodHead {
			return
		}
		// 如果是 echo 的错误，按照 echo 的风格返回
		if errors.As(err, &he) {
			err = c.JSON(he.Code, map[string]any{
				"code": he.Code,
				"msg":  he.Message,
				"data": nil,
			})
		}
		// 返回错误内容，后续可以更改为指定的业务错误
		// 返回 errors.New() 里面的字符串内容
		if err != nil {
			err = c.JSON(http.StatusInternalServerError, map[string]any{
				"code": "-1",
				"msg":  err.Error(),
				"data": nil,
			})
		}
		// 记录日志，不用记录了有日志中间件记录
		//logContent := fmt.Sprintf("method: %s, path: %s, status: %d, client_ip: %s, error: %v", c.Request().Method, c.Path(), code, c.RealIP(), err)
		//log.Error("HTTP error response: " + logContent)
	}
}
