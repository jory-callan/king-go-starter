package echoutil

import "github.com/labstack/echo/v4"

const (
	UserIDKey = "userID"
)

// SetUserID 设置当前请求中的用户 ID
func SetUserID(c echo.Context, userID string) {
	c.Set(UserIDKey, userID)
}

// GetUserID 获取当前请求中的用户 ID
// 如果不存在，返回空字符串
func GetUserID(c echo.Context) string {
	if val := c.Get(UserIDKey); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}
