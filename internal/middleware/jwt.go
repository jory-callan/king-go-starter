package middleware

import (
	"king-starter/internal/common"
	"king-starter/pkg/jwt"

	"github.com/labstack/echo/v4"
)

func JWT(jwt *jwt.JWT) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, map[string]interface{}{
					"code": 401,
					"msg":  "missing authorization header",
					"data": nil,
				})
			}

			claims, err := jwt.ParseToken(authHeader)
			if err != nil {
				return c.JSON(401, map[string]interface{}{
					"code": 401,
					"msg":  "invalid or expired token",
					"data": nil,
				})
			}

			c.Set(string(common.UserIDKey), claims.UserID)
			c.Set(string(common.UsernameKey), claims.Username)

			return next(c)
		}
	}
}
