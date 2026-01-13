package http

import (
	"context"
	"golang.org/x/time/rate"
	"king-starter/pkg/http/middleware"
	"king-starter/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// Server HTTP服务器封装（精简版）
type Server struct {
	log    *logger.Logger
	echo   *echo.Echo
	config *HttpConfig
}

// New 创建Echo服务器
func New(cfg *HttpConfig, log *logger.Logger) *Server {
	// 创建Echo实例
	e := echo.New()
	// 隐藏Banner
	e.HideBanner = true
	e.HidePort = true
	e.Debug = cfg.EnableDebug
	// 直接使用 echo 内置 Server（无需手动创建 http.Server）
	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Millisecond
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Millisecond
	e.Server.MaxHeaderBytes = cfg.MaxHeaderBytes

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// 创建服务器实例
	server := &Server{
		log:    log.With(logger.String("component", "http")),
		echo:   e, // Echo实例
		config: cfg,
	}

	// 注册中间件
	server.registerMiddleware()

	// 注册健康检查
	e.GET("/health", server.healthCheck)

	server.log.Info("http server created",
		logger.String("addr", cfg.Addr()),
		logger.Duration("read_timeout", time.Duration(cfg.ReadTimeout)*time.Millisecond),
	)

	return server
}

// Engine 为了方便外部访问 echo 实例（用于注册路由）
func (s *Server) Engine() *echo.Echo {
	return s.echo
}

// RegisterRoutes 注册路由（符合 Go 惯例）
func (s *Server) RegisterRoutes(registerFunc func(e *echo.Echo)) {
	registerFunc(s.echo)
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
	s.echo.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	// 自定义的中间件
	// 请求日志中间件（使用我们的logger）
	s.echo.Use(middleware.EchoLogger(s.log))
	// 错误处理中间件
	s.echo.HTTPErrorHandler = middleware.EchoErrorHandler(s.log)

}

// healthCheck 健康检查
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"service": "http-server",
	})
}

func (s *Server) Start() error {
	// 启动信号监听（非阻塞）
	s.startSignalHandler()
	// 打印启动信息
	s.log.Info("⇨ http server started on " + s.config.Addr())
	// 阻塞直到服务关闭
	return s.echo.Start(s.config.Addr())
}

func (s *Server) startSignalHandler() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		s.log.Info("received signal", logger.String("signal", sig.String()))
		s.Shutdown()
	}()
}

func (s *Server) Shutdown() {
	s.log.Info("http server shutting down")
	if s.echo.Server == nil {
		s.log.Error("http server not started")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.ShutdownTimeout)*time.Millisecond)
	defer cancel()
	err := s.echo.Shutdown(ctx)
	if err != nil {
		s.log.Error("http server echo shutdown error", logger.Error(err))
		return
	}
	s.log.Info("http server stopped")
}
