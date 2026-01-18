package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zapLogger 封装 zap.Logger，提供统一日志接口
type zapLogger struct {
	*zap.Logger
}

// New 创建 zapLogger 实例，使用提供的配置
func New(cfg *LoggerConfig) (*zapLogger, error) {
	// 解析日志级别
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	// 创建 encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // 颜色高亮：INFO, DEBUG
		//EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05")) // 人类可读时间
		},
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 只显示 main.go:12
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		//  JSON 格式,去掉颜色
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// 文本格式, 保留颜色
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
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
		return nil, fmt.Errorf("invalid output config: %s", cfg.Output)
	}

	// 创建 core
	core := zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(level))

	// 创建 logger
	zapLog := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(cfg.CallerSkip),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return &zapLogger{
		Logger: zapLog.Named("log"),
	}, nil
}

// HealthCheck 健康检查
func (l *zapLogger) HealthCheck(ctx context.Context) error {
	return nil
}

// Close 关闭 logger
func (l *zapLogger) Close() {
	// 直接忽略错误
	l.Logger.Sync()
}

// 封装常用方法给第三方使用

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, fields...)
}
func (l *zapLogger) Info(msg string, fields ...Field) {
	//按照时间
	l.Logger.Info(msg, fields...)
}
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, fields...)
}
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, fields...)
}
func (l *zapLogger) Any(msg string, args ...any) {
	l.Logger.Info(msg, Any("args", args))
}

func (l *zapLogger) Debugf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Debug(msg)
}
func (l *zapLogger) Infof(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg)
}
func (l *zapLogger) Warnf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Warn(msg)
}
func (l *zapLogger) Errorf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Error(msg)
}
func (l *zapLogger) Anyf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg, Any("args", args))
}
