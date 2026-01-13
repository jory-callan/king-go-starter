package middleware

import (
	"errors"
	"king-starter/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

// echoErrorHandler 统一错误处理
func EchoErrorHandler(log *logger.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		// Send response
		code := http.StatusInternalServerError
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
		}
		//log.Errorf("%+v", err)
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, map[string]any{
				"code": code,
				"msg":  "内部未知错误",
				"data": nil,
			})
		}
		if err != nil {
			log.Error(err.Error())
		}
	}
}
