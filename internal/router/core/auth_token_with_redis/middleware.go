package auth_token

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	UserIDKey   = "user_id"
	UsernameKey = "username"
	RolesKey    = "roles"
)

// JWTMiddleware creates an Echo middleware for JWT authentication
func (at *AuthToken) JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code": 401,
					"msg":  "missing authorization header",
					"data": nil,
				})
			}

			// Split the header by space to get the token part
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code": 401,
					"msg":  "authorization header must be in the format 'Bearer {token}'",
					"data": nil,
				})
			}

			tokenString := parts[1]

			// Validate the token
			payload, err := at.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code": 401,
					"msg":  "invalid or expired token: " + err.Error(),
					"data": nil,
				})
			}

			// Store user information in the context
			c.Set(UserIDKey, payload.UserID)
			c.Set(UsernameKey, payload.Username)
			c.Set(RolesKey, payload.Roles)

			return next(c)
		}
	}
}