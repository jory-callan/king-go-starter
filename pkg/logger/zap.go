package logger

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"syscall"
	"time"
)

// Logger 封装 zap.Logger，提供统一日志接口
type Logger struct {
	*zap.Logger
}

// New 创建 Logger 实例，使用提供的配置
func New(cfg *LoggerConfig) (*Logger, error) {
	// 解析日志级别
	level := zapcore.InfoLevel

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
		Logger: zapLogger.Named("log"),
	}, nil
}

// HealthCheck 健康检查
func (l *Logger) HealthCheck(ctx context.Context) error {
	return nil
}

// Close 关闭 logger（安全忽略 stdout/stderr 的 sync 错误）
func (l *Logger) Close() {
	err := l.Logger.Sync()
	if err != nil {
		// 忽略 stdout/stderr 的 sync 错误
		if isStdoutStderrSyncError(err) {
			return // 视为成功
		}
		l.Logger.Error("close logger failed", zap.Error(err))
		return
	}
	return
}

// isStdoutStderrSyncError 判断是否为 stdout/stderr 的 sync 错误
func isStdoutStderrSyncError(err error) bool {
	// 检查是否为 syscall.ENOTTY（inappropriate ioctl for device）
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errors.Is(errno, syscall.ENOTTY)
	}
	// 兼容其他系统可能的错误信息（如 macOS）
	return err.Error() == "sync /dev/stdout: inappropriate ioctl for device" ||
		err.Error() == "sync /dev/stderr: inappropriate ioctl for device"
}

// With 添加自定义字段到 logger
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{
		Logger: l.Logger.Named(name),
	}
}

// 封装常用方法给第三方使用

func (l *Logger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, fields...)
}
func (l *Logger) Info(msg string, fields ...Field) {
	//按照时间
	l.Logger.Info(msg, fields...)
}
func (l *Logger) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, fields...)
}
func (l *Logger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, fields...)
}
func (l *Logger) Any(msg string, args ...any) {
	l.Logger.Info(msg, Any("args", args))
}

func (l *Logger) Debugf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Debug(msg)
}
func (l *Logger) Infof(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg)
}
func (l *Logger) Warnf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Warn(msg)
}
func (l *Logger) Errorf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Error(msg)
}
func (l *Logger) Anyf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg, Any("args", args))
}
