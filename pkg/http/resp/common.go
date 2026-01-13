package resp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func JSON(c echo.Context, resp Response) {
	c.JSON(http.StatusOK, resp)
}

func SuccessJSON(c echo.Context, data interface{}) error {
	JSON(c, Response{
		Code:    SUCCESS.Code,
		Message: SUCCESS.Message,
		Data:    data,
	})
	return nil
}

func ErrorJSON(c echo.Context, cm CodeMessage) error {
	JSON(c, Response{
		Code:    cm.Code,
		Message: cm.Message,
		Data:    nil,
	})
	return nil
}
