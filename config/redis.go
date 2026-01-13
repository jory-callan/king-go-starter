package config

import (
	"fmt"
)

// RedisConfig Redis配置
type RedisConfig struct {
	// 必填字段
	Addr     string `yaml:"addr" json:"addr"`         // 地址: localhost:6379
	Password string `yaml:"password" json:"password"` // 密码
	Prefix   string `yaml:"prefix" json:"prefix"`     // 键前缀（强制要求）
	// 可选字段（有默认值）
	DB         int `yaml:"db" json:"db"`                   // 数据库编号
	MaxRetries int `yaml:"max_retries" json:"max_retries"` // 最大重试次数
	// 连接池配置 默认是线程数量, 不够会自动扩容 可以通过配置 max_active_conns 来调整最大连接池
	PoolSize        int `yaml:"pool_size" json:"pool_size"`
	MaxActiveConns  int `yaml:"max_active_conns" json:"max_active_conns"`     // 最大活跃连接数
	DialTimeout     int `yaml:"dial_timeout" json:"dial_timeout"`             // 连接超时（毫秒）
	ReadTimeout     int `yaml:"read_timeout" json:"read_timeout"`             // 读取超时（毫秒）
	WriteTimeout    int `yaml:"write_timeout" json:"write_timeout"`           // 写入超时（毫秒）
	MinIdleConns    int `yaml:"min_idle_conns" json:"min_idle_conns"`         // 最小空闲连接
	MaxIdleConns    int `yaml:"max_idle_conns" json:"max_idle_conns"`         // 最大空闲连接
	ConnMaxLifetime int `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`   // 连接最大生命周期（毫秒）
	ConnMaxIdleTime int `yaml:"conn_max_idle_time" json:"conn_max_idle_time"` // 连接最大空闲时间（毫秒）
}

// Validate 配置校验
func (c *RedisConfig) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("[redis] RedisConfig error: addr is required")
	}
	if c.Prefix == "" {
		return fmt.Errorf("[redis] RedisConfig error: prefix is required")
	}
	return nil
}

func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Addr:            "localhost:6379",
		Password:        "",
		Prefix:          "app:cache:",
		DB:              0,
		MaxRetries:      3,
		DialTimeout:     5000,
		ReadTimeout:     3000,
		WriteTimeout:    3000,
		MaxActiveConns:  100,
		MinIdleConns:    10,
		MaxIdleConns:    100,
		ConnMaxLifetime: 500000,
		ConnMaxIdleTime: 30000,
	}
}

/*
redis:
  # 默认缓存实例
  cache:
    addr: "localhost:6379"
    password: ""
    prefix: "app:cache:"      # 强制前缀，键隔离
    db: 0
    pool_size: 100
    min_idle_conns: 10
    dial_timeout: "5s"
    read_timeout: "3s"
    write_timeout: "3s"

  # 会话存储实例
  session:
    addr: "localhost:6380"
    password: "secret123"
    prefix: "app:session:"
    db: 1
    pool_size: 50
    min_idle_conns: 5

  # 消息队列实例
  queue:
    addr: "localhost:6381"
    password: ""
    prefix: "app:queue:"
    db: 2
    pool_size: 30
    min_idle_conns: 3
    conn_max_lifetime: "30m"
*/
