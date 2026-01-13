package http

// import (
// 	"net/http"
// 	"king-starter/config"
// 	"king-starter/pkg/app"
// 	"king-starter/pkg/log"
// 	"testing"
// 	"time"

// 	"github.com/labstack/echo/v4"
// )

// func TestServer(t *testing.T) {
// 	// 初始化logger
// 	log := log.NewWithConfig(&log.Config{
// 		Level: "debug",
// 	})

// 	// 配置
// 	cfg := Config{
// 		Host: "0.0.0.0",
// 		Port: 8080,
// 	}
// 	// 创建核心
// 	core := app.NewWithConfig(&config.DefaultConfig())

// 	// 创建服务器
// 	server := NewWithConfig(&cfg, log, core)

// 	// 注册路由（直接使用Echo原生方法）
// 	server.echo.GET("/", func(c echo.Context) error {
// 		return c.String(http.StatusOK, "Hello, World!")
// 	})

// 	server.echo.GET("/users/:id", func(c echo.Context) error {
// 		id := c.Param("id")
// 		return c.JSON(http.StatusOK, map[string]string{
// 			"id":   id,
// 			"name": "John",
// 		})
// 	})

// 	// sleep 1s
// 	server.echo.GET("/sleep", func(c echo.Context) error {
// 		time.Sleep(1 * time.Second)
// 		return c.String(http.StatusOK, "Hello, World!")
// 	})

// 	// POST示例
// 	server.echo.POST("/users", func(c echo.Context) error {
// 		var user map[string]interface{}
// 		if err := c.Bind(&user); err != nil {
// 			return err
// 		}
// 		return c.JSON(http.StatusCreated, user)
// 	})

// 	// 启动服务器
// 	server.Start()

// 	// 等待信号优雅关闭
// 	server.WaitForSignal()
// }
