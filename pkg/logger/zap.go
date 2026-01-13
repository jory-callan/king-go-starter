package logger

import (
	"fmt"
	"king-starter/config"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 封装 zap.Logger，提供统一日志接口
type Logger struct {
	*zap.Logger
	once sync.Once
}

// New 创建 Logger 实例，使用提供的配置
func New(cfg *config.LoggerConfig) *Logger {
	// 解析日志级别
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		panic(fmt.Sprintf("invalid log level: %v", err))
	}

	// 创建 encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建写入器
	var writer zapcore.WriteSyncer
	if cfg.Output == "stdout" {
		writer = zapcore.AddSync(os.Stdout)
	} else if cfg.Output == "file" && cfg.FilePath != "" {
		writer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
	} else {
		panic(fmt.Sprintf("invalid output config: %s", cfg.Output))
	}

	// 创建 core
	core := zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(level))

	// 创建 logger
	zapLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(cfg.CallerSkip),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &Logger{
		Logger: zapLogger.Named("logger"),
	}
}

// NewWithDefaultConfig 使用默认配置创建 Logger 实例
func NewWithDefaultConfig() *Logger {
	cfg := config.DefaultLoggerConfig()
	return New(&cfg)
}

// With 添加自定义字段到 logger
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}

// Named 添加 logger 名称前缀
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		Logger: l.Logger.Named(name),
	}
}
