package middleware

import (
	"errors"
	"time"

	"king-starter/pkg/logx"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// EchoLogger 日志中间件
func EchoLogger(logx logx.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			// 提取状态码时要考虑错误可能修改的状态码
			// Extract valuable fields, exactly like upstream logger
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			//status := res.Status
			//if err != nil {
			//	// c.Error(err) may have updated status (e.g. via HTTPErrorHandler)
			//	// echo sets res.Status on error, so this is safe
			//	status = res.Status
			//}
			status := res.Status
			if err != nil {
				var he *echo.HTTPError
				if errors.As(err, &he) {
					status = he.Code
				}
			}
			// 构建基础字段（固定顺序，便于阅读）
			args := []any{
				"id", id,
				"request_id", id,
				"remote_ip", c.RealIP(),
				"host", req.Host,
				"proto", req.Proto,
				"method", req.Method,
				"uri", req.RequestURI,
				"path", req.URL.Path,
				"route", c.Path(),
				"user_agent", req.UserAgent(),
				"referer", req.Referer(),
				"status", status,
				"latency", latency,
				"latency_human", latency.String(),
				"bytes_in", req.ContentLength,
				"bytes_out", res.Size,
			}

			// 动态字段：按需追加
			if qs := req.URL.RawQuery; qs != "" {
				args = append(args, "query_string", qs)
			}
			if contentType := req.Header.Get("Content-Type"); contentType != "" {
				args = append(args, "content_type", contentType)
			}

			// 使用 logx 记录日志
			logx.Info("HTTP request completed", args...)

			if err != nil {
				// Do NOT use %v — may leak stack or secrets
				args = append(args, zap.Error(err))
			}

			// 根据状态码决定日志级别
			switch {
			case status >= 500:
				logx.Error("http 500 error", args...)
			case status >= 400:
				logx.Warn("http 400 error", args...)
			default:
				logx.Info("http request", args...)
			}

			// 如果有错误，则返回错误，让下一个中间件处理
			return err
		}
	}
}
