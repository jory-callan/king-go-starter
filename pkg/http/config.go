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
	ReadTimeout     int `yaml:"read_timeout" json:"read_timeout"`         // 读取超时（毫秒）
	WriteTimeout    int `yaml:"write_timeout" json:"write_timeout"`       // 写入超时（毫秒）
	ShutdownTimeout int `yaml:"shutdown_timeout" json:"shutdown_timeout"` // 关闭超时（毫秒）
	IdleTimeout     int `yaml:"idle_timeout" json:"idle_timeout"`         // 空闲连接超时时间（毫秒）
	MaxHeaderBytes  int `yaml:"max_header_bytes" json:"max_header_bytes"` // 最大请求头字节数
	MaxBodyBytes    int `yaml:"max_body_bytes" json:"max_body_bytes"`     // 最大请求体字节数
	// CORE 跨域配置
	CORS CORSConfig `yaml:"cors" json:"cors"` // CORS 跨域配置
	// DEBUG 调试配置
	EnableDebug bool `yaml:"enable_debug" json:"enable_debug"` // 是否启用调试模式
	// RATE_LIMIT 限流配置
	RateLimit RateLimitConfig `yaml:"rate_limit" json:"rate_limit"` // 限流配置
}

type CORSConfig struct {
	AllowOrigins  []string `yaml:"allow_origins" json:"allow_origins"`   // 允许的源
	AllowMethods  []string `yaml:"allow_methods" json:"allow_methods"`   // 允许的 HTTP 方法
	AllowHeaders  []string `yaml:"allow_headers" json:"allow_headers"`   // 允许的请求头
	ExposeHeaders []string `yaml:"expose_headers" json:"expose_headers"` // 暴露给客户端的响应头
	MaxAge        int      `yaml:"max_age" json:"max_age"`               // 预检请求缓存时间(秒)
}

type RateLimitConfig struct {
	Enabled           bool `yaml:"enabled" json:"enabled"`                         // 是否启用限流
	RequestsPerSecond int  `yaml:"requests_per_second" json:"requests_per_second"` // 每秒允许的请求数
	Burst             int  `yaml:"burst" json:"burst"`                             // 突发请求允许的最大数量
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

// Addr 获取监听地址
func (c *HttpConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
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
		EnableDebug:     true,
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
