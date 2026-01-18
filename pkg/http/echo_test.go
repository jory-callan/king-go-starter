package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestEcho(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	go func() {
		// 创建信号通道
		quit := make(chan os.Signal, 1)
		// 监听 Ctrl+C (SIGINT) 和 终止信号 (SIGTERM)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// 这里会阻塞，直到收到信号
		<-quit
		fmt.Println("正在优雅关闭服务器...")
		// 设置一个 5秒 的超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 执行关闭操作
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}()
	// 启动服务（主 goroutine 在这里阻塞）
	// 注意：不需要用 if err != nil 判断，因为 Shutdown 正常执行时，Start() 会返回 http.ErrServerClosed
	fmt.Println("服务启动在 :8081")
	if err := e.Start(":8081"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("服务启动失败: " + err.Error())
	}

	fmt.Println("服务器已完全停止")
}
