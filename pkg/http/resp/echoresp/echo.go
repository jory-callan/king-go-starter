// Package echoresp provides HTTP response helpers for Echo framework.
package echoresp

import (
	"king-starter/pkg/http/resp"
	"net/http"

	"github.com/labstack/echo/v4"
)

// func ToEchoStatus(status int) int {
// 	switch status {
// 	case resp.StatusSuccess:
// 		return echo.StatusOK
// 	case resp.StatusClientError:
// 		return echo.StatusBadRequest
// 	case resp.StatusServerError:
// 		return echo.StatusInternalServerError
// 	default:
// 		return echo.StatusInternalServerError
// 	}
// }

func ToHTTPStatus(c resp.CodeMsg) int {
	switch c.Status {
	case resp.StatusSuccess:
		return http.StatusOK
	case resp.StatusClientError:
		if c.Code == "401" {
			return http.StatusUnauthorized
		}
		if c.Code == "403" {
			return http.StatusForbidden
		}
		if c.Code == "404" {
			return http.StatusNotFound
		}
		return http.StatusBadRequest
	case resp.StatusServerError:
		if c.Code == "500" {
			return http.StatusInternalServerError
		}
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}

// JSON sends a resp.Response as JSON with given HTTP status code.
func JSON(c echo.Context, statusCode int, r resp.Response) error {
	return c.JSON(statusCode, r)
}

// Success sends a success response with data (code=0).
func Success(c echo.Context, statusCode int, data interface{}) error {
	return JSON(c, statusCode, resp.Success(resp.WithData(data)))
}

// SuccessMsg sends a success response with custom message.
func SuccessMsg(c echo.Context, statusCode int, msg string) error {
	return JSON(c, statusCode, resp.Success(resp.WithMsg(msg)))
}

// Error sends a failure response using predefined CodeMsg.
func Error(c echo.Context, statusCode int, codeMsg resp.CodeMsg) error {
	return JSON(c, statusCode, resp.Error(codeMsg))
}

// ErrorMsg sends a failure response with custom message.
func ErrorMsg(c echo.Context, statusCode int, codeMsg resp.CodeMsg, msg string) error {
	return JSON(c, statusCode, resp.Error(codeMsg, resp.WithMsg(msg)))
}
