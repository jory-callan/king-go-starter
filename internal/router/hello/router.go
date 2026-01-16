package hello

import (
	"king-starter/internal/app"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Router struct{}

func New() *Router {
	return &Router{}
}

func (m *Router) Name() string {
	return "hello"
}

func (m *Router) Register(app *app.App) {
	e := app.Server.Engine()
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world")
	})
}
