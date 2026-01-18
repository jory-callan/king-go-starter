package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapConfig mirrors your original LoggerConfig
type ZapConfig struct {
	Level      string
	Format     string // "json" or "console"
	Output     string // "stdout" or "file"
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	CallerSkip int
}

type zapLogger struct {
	*zap.Logger
	closer io.Closer
}

func myCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	callFullPath := caller.FullPath()
	gomodAbsPath, _ := filepath.Abs("go.mod")
	absDir := filepath.Dir(gomodAbsPath)
	//prefix := strings.TrimPrefix(callFullPath, absDir)
	path, _ := filepath.Rel(absDir, callFullPath)
	enc.AppendString(path)
}
func newZapLogger(cfg *LoggerConfig) (*zapLogger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	// 用 console 输出作为模板
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",                           // json下的时间的key,
		LevelKey:      "level",                          // json下的日志级别的key
		NameKey:       "logger",                         // 调用 Named 时 json下的key
		CallerKey:     "caller",                         // 调用者key
		MessageKey:    "msg",                            // 消息key
		StacktraceKey: "stacktrace",                     // 堆栈key
		LineEnding:    zapcore.DefaultLineEnding,        // 换行符
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // 颜色高亮：INFO, DEBUG
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { // 时间格式
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration:   zapcore.StringDurationEncoder, // 毫秒
		EncodeCaller:     myCallerEncoder,               // 只显示 main.go:12
		ConsoleSeparator: "  ",                          // 间隔符
	}

	if cfg.Format == "json" {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var writer zapcore.WriteSyncer
	var closer io.Closer = nil
	switch cfg.Output {
	case "stdout":
		writer = zapcore.AddSync(os.Stdout)
	case "file":
		if cfg.FilePath == "" {
			return nil, fmt.Errorf("file output requires FilePath")
		}
		lj := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		closer = lj
		writer = zapcore.AddSync(lj)
	default:
		return nil, fmt.Errorf("invalid output: %s", cfg.Output)
	}

	core := zapcore.NewCore(encoder, writer, level)
	logger := zap.New(core,
		zap.AddCaller(),
		//zap.AddCallerSkip(cfg.AddCallerSkip),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	return &zapLogger{logger, closer}, nil
}

// Implement Logger interface

func (z *zapLogger) Info(msg string, args ...any)  { z.Logger.Info(msg, kvToZapFields(args...)...) }
func (z *zapLogger) Warn(msg string, args ...any)  { z.Logger.Warn(msg, kvToZapFields(args...)...) }
func (z *zapLogger) Error(msg string, args ...any) { z.Logger.Error(msg, kvToZapFields(args...)...) }
func (z *zapLogger) Debug(msg string, args ...any) { z.Logger.Debug(msg, kvToZapFields(args...)...) }

func (z *zapLogger) Panic(msg string, args ...any) {
	z.Logger.Error(msg, kvToZapFields(args...)...)
	panic(msg)
}

func (z *zapLogger) Infof(format string, args ...any)  { z.Logger.Info(fmt.Sprintf(format, args...)) }
func (z *zapLogger) Warnf(format string, args ...any)  { z.Logger.Warn(fmt.Sprintf(format, args...)) }
func (z *zapLogger) Errorf(format string, args ...any) { z.Logger.Error(fmt.Sprintf(format, args...)) }
func (z *zapLogger) Debugf(format string, args ...any) { z.Logger.Debug(fmt.Sprintf(format, args...)) }
func (z *zapLogger) Panicf(format string, args ...any) {
	z.Logger.Error(fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
}

func (z *zapLogger) With(args ...any) Logger {
	return &zapLogger{z.Logger.With(kvToZapFields(args...)...), z.closer}
}

func (z *zapLogger) Named(name string) Logger {
	return &zapLogger{z.Logger.Named(name), z.closer}
}

func (z *zapLogger) AddCallerSkip(skip int) Logger {
	return &zapLogger{z.Logger.WithOptions(zap.AddCallerSkip(skip)), z.closer}
}

func (z *zapLogger) Close() {
	if z.closer != nil {
		err := z.closer.Close()
		if err != nil {
			return
		}
	}
}

// kvToZapFields converts key-value pairs to []zap.Field
func kvToZapFields(args ...any) []zap.Field {
	var fields []zap.Field
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		fields = append(fields, zap.Any(key, args[i+1]))
	}
	return fields
}
