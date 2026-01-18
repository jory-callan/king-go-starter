package middleware

import (
	"time"

	"king-starter/pkg/logx"

	"github.com/labstack/echo/v4"
)

const rateLimitWindow = time.Minute

func RateLimit() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logx.Info("middleware ==> rate limit , ip " + c.RealIP())
			//key := fmt.Sprintf("rate_limit:%s", c.RealIP())
			//ctx := context.Background()
			//pipe := rdb.Pipeline()
			//
			//incr := pipe.Incr(ctx, key)
			//pipe.Expire(ctx, key, rateLimitWindow)
			//
			//if _, err := pipe.Exec(ctx); err != nil {
			//	return next(c)
			//}
			//
			//count, _ := incr.Result()
			//if count > int64(requestsPerMinute) {
			//	return c.JSON(429, map[string]interface{}{
			//		"code": 429,
			//		"msg":  "too many requests",
			//		"data": nil,
			//	})
			//}

			return next(c)
		}
	}
}
