package middleware

import (
	"fmt"
	"net/http"

	"king-starter/pkg/http/resp"
	"king-starter/pkg/http/resp/echoresp"
	"king-starter/pkg/logx"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// EchoRecover 主协程panic处理。
func EchoRecover(logx logx.Logger) echo.MiddlewareFunc {
	return echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{
		StackSize: 2 << 10, // 2 KB
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			errMsg := fmt.Sprintf("内部服务错误: errmsg: %s , stack: %s", err.Error(), string(stack))
			logx.Error(errMsg)
			err = echoresp.Error(c, http.StatusInternalServerError, resp.ErrInternalServer)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
