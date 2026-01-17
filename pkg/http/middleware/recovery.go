package middleware

import (
	"king-starter/pkg/http/resp"
	"king-starter/pkg/http/resp/echoresp"
	"king-starter/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// EchoRecover 主协程panic处理。
func EchoRecover(log *logger.Logger) echo.MiddlewareFunc {
	return echoMiddleware.RecoverWithConfig(echoMiddleware.RecoverConfig{
		StackSize: 2 << 10, // 2 KB
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error("内部服务错误: %v\n%s", zap.Error(err))
			err = echoresp.Error(c, http.StatusInternalServerError, resp.ErrInternalServer)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
