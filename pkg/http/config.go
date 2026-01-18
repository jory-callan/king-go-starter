package http

import (
	"fmt"
)

// HttpConfig HTTP服务器配置
type HttpConfig struct {
	// 必填
	Host string `yaml:"host" json:"host"` // 监听地址
	Port int    `yaml:"port" json:"port"` // 监听端口

	// 可选，有默认值
	ReadTimeout     int // 读取超时（毫秒）
	WriteTimeout    int // 写入超时（毫秒）
	ShutdownTimeout int // 关闭超时（毫秒）
	IdleTimeout     int // 空闲连接超时时间（毫秒）
	MaxHeaderBytes  int // 最大请求头字节数
	MaxBodyBytes    int // 最大请求体字节数
	// CORE 跨域配置
	CORS CORSConfig // CORS 跨域配置
	// DEBUG 调试配置
	EnableDebug bool // 是否启用调试模式
	// RATE_LIMIT 限流配置
	RateLimit RateLimitConfig // 限流配置
}

type CORSConfig struct {
	AllowOrigins  []string // 允许的源
	AllowMethods  []string // 允许的 HTTP 方法
	AllowHeaders  []string // 允许的请求头
	ExposeHeaders []string // 暴露给客户端的响应头
	MaxAge        int      // 预检请求缓存时间(秒)
}

type RateLimitConfig struct {
	Enabled           bool // 是否启用限流
	RequestsPerSecond int  // 每秒允许的请求数
	Burst             int  // 突发请求允许的最大数量
}

// Validate 配置校验
func (c *HttpConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("[http] HttpConfig error: host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("[http] HttpConfig error: port %d is invalid", c.Port)
	}
	if c.RateLimit.Enabled {
		if c.RateLimit.RequestsPerSecond <= 0 {
			return fmt.Errorf("[http] HttpConfig error: requests_per_second %d is invalid", c.RateLimit.RequestsPerSecond)
		}
		if c.RateLimit.Burst <= 0 {
			return fmt.Errorf("[http] HttpConfig error: burst %d is invalid", c.RateLimit.Burst)
		}
	}
	return nil
}

// DefaultHttpConfig 默认配置
func DefaultHttpConfig() HttpConfig {
	return HttpConfig{
		Host:            "127.0.0.1",
		Port:            8080,
		ReadTimeout:     15000,
		WriteTimeout:    15000,
		ShutdownTimeout: 10000,
		IdleTimeout:     60000,
		MaxHeaderBytes:  1048576,
		MaxBodyBytes:    104857600,
		EnableDebug:     false,
		CORS: CORSConfig{
			AllowOrigins:  []string{"*"},
			AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"Content-Type", "Authorization"},
			ExposeHeaders: []string{"X-Request-ID"},
			MaxAge:        86400,
		},
		RateLimit: RateLimitConfig{
			Enabled:           true,
			RequestsPerSecond: 100,
			Burst:             20,
		},
	}
}

/*
http:
  host: "0.0.0.0"     # 监听地址
  port: 8080          # 监听端口
  read_timeout: 15000   # 读取请求超时时间
  write_timeout: 15000  # 写入响应超时时间
  shutdown_timeout: 10000   # 关闭超时时间
  idle_timeout: 60000   # 空闲连接超时时间
  max_header_bytes: 1048576   # 最大请求头字节数
  max_body_bytes: 104857600   # 最大允许的请求体大小 (支持 KB, MB, GB)
  # CORS 跨域配置
  cors:
    allow_origins: ["https://example.com", "http://localhost:3000"]  # 允许的源
    allow_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]      # 允许的 HTTP 方法
    allow_headers: ["Content-Type", "Authorization"]                # 允许的请求头
    expose_headers: ["X-Request-ID"]                                # 暴露给客户端的响应头
    max_age: 86400                                                  # 预检请求缓存时间(秒)
  # 限流配置 (可选)
  rate_limit:
    enabled: true       # 是否启用限流
    requests_per_second: 100  # 每秒允许的请求数
    burst: 20           # 突发请求允许的最大数量
*/
