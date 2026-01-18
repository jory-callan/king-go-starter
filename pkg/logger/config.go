package logger

import "fmt"

type LoggerConfig struct {
	Level      string // 日志级别: debug, info, warn, error, fatal
	Format     string // 日志格式: json, text
	Output     string // 输出目标: stdout, 文件路径（如 /var/log/app.log）
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小(MB)
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 旧日志文件保留天数(天)
	Compress   bool   // 是否压缩旧日志文件
	CallerSkip int    // 跳过调用栈的层数
}

func (c *LoggerConfig) Validate() error {
	// 验证日志级别
	switch c.Level {
	case "debug", "info", "warn", "error", "fatal":
	default:
		return fmt.Errorf("[logger] config invalid log level: %s", c.Level)
	}

	// 验证日志格式
	switch c.Format {
	case "json", "text":
	default:
		return fmt.Errorf("[logger] config invalid log format: %s", c.Format)
	}

	// 验证输出目标
	switch c.Output {
	case "stdout", "file":
	default:
		return fmt.Errorf("[logger] config invalid log output: %s", c.Output)
	}

	// 如果输出目标是 file 则必须配置日志文件路径
	if c.Output == "file" && c.FilePath == "" {
		return fmt.Errorf("[logger] config file path is required for file output")
	}

	return nil
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		FilePath:   "./app-king.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
		CallerSkip: 1,
	}
}

/*
logger:
  level: "info"
  format: "json"
  output: "stdout"
  file_path: "/var/log/app.log"
  max_size: 100
  max_backups: 5
  max_age: 7
  compress: true
  caller_skip: 2
*/
