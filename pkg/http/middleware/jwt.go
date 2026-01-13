package middleware

import (

	// "king-starter/pkg/app"
	"king-starter/pkg/http/resp"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logger"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTAuthMiddleware(jwt *jwt.JWT, logger *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return resp.ErrorJSON(c, resp.UNAUTHORIZED)
			}
			// 按空格分割
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				return resp.ErrorJSON(c, resp.UNAUTHORIZED)
			}
			// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它；也会自动校验过期时间
			payload, err := jwt.Parse(parts[1])
			if err != nil {
				return resp.ErrorJSON(c, resp.UNAUTHORIZED)
			}
			// 将当前请求的username信息保存到请求的上下文c上
			c.Set("UserID", payload.UserID)
			return next(c)
		}
	}
}

// func AuthMiddleware(ctx *app.AppContext) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			auth := c.Request().Header.Get("Authorization")
// 			if auth == "" {
// 				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
// 			}
// 			parts := strings.SplitN(auth, " ", 2)
// 			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
// 			}

// 			claims, err := ctx.Jwt.Parse(parts[1])
// 			if err != nil {
// 				ctx.Logger.Info("jwt parse failed", zap.Error(err)) // ✅ 用你封装的日志
// 				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
// 			}

// 			c.Set("user", claims)
// 			return next(c)
// 		}
// 	}
// }
