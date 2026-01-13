package middleware

import (
	"king-starter/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func EchoLogger(log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			// Extract valuable fields, exactly like upstream logger
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			status := res.Status
			if err != nil {
				// c.Error(err) may have updated status (e.g. via HTTPErrorHandler)
				// echo sets res.Status on error, so this is safe
				status = res.Status
			}

			fields := []zap.Field{
				//id 采用 X-Request-ID
				zap.String("id", id),
				// 请求ID X-Request-ID
				zap.String("request_id", id),
				// 请求IP 127.0.0.1
				zap.String("remote_ip", c.RealIP()),
				// 请求主机 host:port
				zap.String("host", req.Host),
				// 请求协议 HTTP/1.1
				zap.String("proto", req.Proto),
				// 请求方法 GET
				zap.String("method", req.Method),
				// 请求URI /api/v1/health
				zap.String("uri", req.RequestURI),
				// 请求路径 /api/v1/health
				zap.String("path", req.URL.Path),
				// 路由路径 /api/v1/health
				zap.String("route", c.Path()), // actual matched route, e.g. "/users/:id"
				// 用户代理 User-Agent
				zap.String("user_agent", req.UserAgent()),
				// 引用者 Referer
				zap.String("referer", req.Referer()),
				// 状态码 200
				zap.Int("status", status),
				// 耗时 1.234ms
				zap.Duration("latency", latency),
				// 耗时人类可读 1.234ms
				zap.String("latency_human", latency.String()),
				// 请求体大小 0
				zap.Int64("bytes_in", req.ContentLength), // raw Content-Length, -1 if absent
				// 响应体大小 123
				zap.Int64("bytes_out", res.Size),
			}

			// Optional: extract dynamic fields only if present
			if qs := req.URL.RawQuery; qs != "" {
				fields = append(fields, zap.String("query_string", qs))
			}

			if contentType := req.Header.Get("Content-Type"); contentType != "" {
				fields = append(fields, zap.String("content_type", contentType))
			}

			// Error field: only non-nil
			if err != nil {
				// Do NOT use %v — may leak stack or secrets
				// Prefer safe string conversion
				fields = append(fields, zap.Error(err))
			}

			// 根据状态码决定日志级别
			switch {
			case status >= 500:
				log.Error("http 500 error", fields...)
			case status >= 400:
				log.Warn("http 400 error", fields...)
			default:
				log.Info("http request", fields...)
			}
			return nil
		}
	}
}
