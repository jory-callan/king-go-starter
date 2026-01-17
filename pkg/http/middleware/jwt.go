package middleware

import (

	// "king-starter/pkg/app"
	"king-starter/pkg/http/resp"
	"king-starter/pkg/http/resp/echoresp"
	"king-starter/pkg/jwt"
	"king-starter/pkg/logger"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware(jwt *jwt.JWT, log *logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("缺少 Authorization 头")
				return echoresp.Error(c, http.StatusUnauthorized, resp.ErrUnauthorized)
			}

			// 按空格分割
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				log.Warn("Authorization 格式错误 => " + "authHeader=" + authHeader)
				return echoresp.Error(c, http.StatusUnauthorized, resp.ErrUnauthorized)
			}

			// parts[1] 是获取到的 tokenString，使用定义好的解析函数并校验过期时间
			claims, err := jwt.ParseToken(parts[1])
			if err != nil {
				log.Warn("JWT 解析失败", logger.Error(err))
				return echoresp.Error(c, http.StatusUnauthorized, resp.ErrUnauthorized)
			}

			// 将当前请求的 userID 保存到请求的上下文 c 上
			c.Set("UserID", claims.UserID)
			return next(c)
		}
	}
}
