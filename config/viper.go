package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	globalConfig *Config
	once         sync.Once
)

// LoadWithFile 加载配置，失败则panic
func LoadWithFile(configPath string) *Config {
	once.Do(func() {
		globalConfig = load(configPath)
		if globalConfig == nil {
			panic("[config] load config failed")
		}
		log.Printf("[config] load config success: %s", globalConfig.Summary())
	})
	return globalConfig
}

func Load() *Config {
	once.Do(func() {
		globalConfig = load("")
		if globalConfig == nil {
			panic("[config] load config failed")
		}
		log.Printf("[config] load config success: %s", globalConfig.Summary())
	})
	return globalConfig
}

// load 实际的配置加载逻辑
func load(configPath string) *Config {
	v := setupViper(configPath)

	// 尝试读取配置
	if err := v.ReadInConfig(); err != nil {
		// 如果明确指定了配置文件但找不到，直接失败
		if configPath != "" {
			panic(fmt.Sprintf("[config] 指定的配置文件不存在: %s, 错误: %v", configPath, err))
		}
		// 未指定配置文件时，使用默认配置
		log.Println("[config] 未找到配置文件，使用默认配置")
	} else {
		log.Printf("[config] 使用配置文件: %s", v.ConfigFileUsed())
	}

	// 先获取默认的配置值
	cfg := DefaultConfig()

	// 解析配置到结构体，覆盖默认值
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Sprintf("[config] 配置解析失败: %v", err))
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("[config] 配置验证失败: %v", err))
	}

	return &cfg
}

// setupViper 配置Viper实例
func setupViper(configPath string) *viper.Viper {
	v := viper.New()

	// 配置文件处理
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// 自动搜索
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath(filepath.Join(os.Getenv("HOME"), "config"))
	}

	// 环境变量配置, 前缀为 APP_
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return v
}
