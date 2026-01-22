package hello

import (
	"errors"
	"net/http"

	"king-starter/internal/app"
	"king-starter/pkg/logx"

	"github.com/labstack/echo/v4"
)

// Register 只是用于最简单的代码测试而已
func Register(app *app.App) {
	e := app.Server.Engine()
	g := e.Group("/api/v1/hello")
	g.GET("/echo", func(c echo.Context) error {
		logx.Info("hello world")
		return c.String(http.StatusOK, "hello world")
	})
	g.GET("/error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, " echo.NewHTTPError(http.StatusBadRequest, xxxx) ")
	})
	g.GET("/error2", func(c echo.Context) error {
		return errors.New(" errors.New(xxxx) ")
	})
	g.GET("/panic", func(c echo.Context) error {
		panic("this will make panic")
	})
}
