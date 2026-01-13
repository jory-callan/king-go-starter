package http

import (
	"context"
	"fmt"
	"king-starter/config"
	"king-starter/pkg/http/middleware"
	"king-starter/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// Server HTTP服务器封装（精简版）
type Server struct {
	logger *logger.Logger
	config *config.HttpConfig
	echo   *echo.Echo
	server *http.Server
}

// New 创建Echo服务器
func New(cfg *config.HttpConfig, logger *logger.Logger) *Server {

	// 创建Echo实例
	e := echo.New()
	e.HideBanner = true

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           cfg.Addr(),
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Millisecond,
		MaxHeaderBytes: 1 << 22,
		Handler:        e,
	}

	server := &Server{
		echo:   e,
		server: srv,
		config: cfg,
		logger: logger.With(zap.String("component", "http")),
	}

	// 注册中间件
	server.registerMiddleware()

	// 注册健康检查
	e.GET("/health", server.healthCheck)

	server.logger.Info("http server created",
		zap.String("addr", cfg.Addr()),
		zap.Duration("read_timeout", time.Duration(cfg.ReadTimeout)*time.Millisecond),
	)

	return server
}

// 为了方便外部访问 echo 实例（用于注册路由）
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

// 直接允许操作 s.echo
// 这样在外部 router 里，你就像拿到了钥匙，可以直接进门布置
func (s *Server) SetupRoutes(setupFunc func(e *echo.Echo)) {
	setupFunc(s.echo)
}

// registerMiddleware 注册中间件
func (s *Server) registerMiddleware() {
	// 原生中间件
	// 使用Echo内置的Recover中间件
	s.echo.Use(echoMiddleware.Recover())
	// 添加RequestID
	s.echo.Use(echoMiddleware.RequestID())
	// s.echo.Use(echoMiddleware.Logger())
	s.echo.Use(echoMiddleware.CORS())

	// 自定义的中间件
	// 请求日志中间件（使用我们的logger）
	s.echo.Use(middleware.EchoLogger(s.logger))
	// 错误处理中间件
	s.echo.HTTPErrorHandler = middleware.EchoErrorHandler(s.logger)

}

// healthCheck 健康检查
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"service": "http-server",
	})
}

// Start 启动服务器
func (s *Server) Start() {
	go func() {
		s.logger.Info("http server starting", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("http server failed", zap.Error(err))
			panic(fmt.Sprintf("http server failed: %v", err))
		}
		s.logger.Info("http server started", zap.String("addr", s.server.Addr))
	}()
}

// StartTLS 启动HTTPS
func (s *Server) StartTLS(certFile, keyFile string) {
	go func() {
		s.logger.Info("https server starting", zap.String("addr", s.server.Addr))

		if err := s.server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			s.logger.Error("https server failed", zap.Error(err))
			panic(fmt.Sprintf("https server failed: %v", err))
		}
	}()
}

// Shutdown 优雅关闭
func (s *Server) Shutdown() error {
	s.logger.Info("http server shutting down")

	if s.server == nil {
		return fmt.Errorf("http server not started")
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.ShutdownTimeout)*time.Millisecond)

	defer cancel()

	return s.server.Shutdown(ctx)
}

// WaitForSignal 等待信号优雅关闭
func (s *Server) WaitForSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	s.logger.Info("received signal", zap.String("signal", sig.String()))

	if err := s.Shutdown(); err != nil {
		s.logger.Error("shutdown failed", zap.Error(err))
	}
}
