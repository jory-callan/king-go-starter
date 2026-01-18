package logx

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type SlogConfig struct {
	Level      string
	Format     string // "json" or "text"
	Output     string // "stdout" or "file"
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	// Note: Slog does not support dynamic caller skip
}

type slogLogger struct {
	*slog.Logger
	closer  io.Closer // 用于 lumberjack.Close()
	logName string
}

func newSlogLogger(cfg *LoggerConfig) (*slogLogger, error) {
	level := parseSlogLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level: level,
		//AddSource: true,
	}

	var handler slog.Handler
	//var writer = os.Stdout
	var writer io.Writer = os.Stdout
	var closer io.Closer = nil // 默认无

	if cfg.Output == "file" {
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
		writer = lj
		closer = lj // lumberjack.Logger 实现了 io.closer
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		//handler = slog.NewTextHandler(writer, opts)
		handler = newZapLikeHandler(writer, opts)
	}

	loggerName := cfg.LogName
	return &slogLogger{
		Logger:  slog.New(handler),
		closer:  closer,
		logName: loggerName,
	}, nil
}

func parseSlogLevel(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Implement Logger interface

func (s *slogLogger) Info(msg string, args ...any)  { s.Logger.Info(msg, args...) }
func (s *slogLogger) Warn(msg string, args ...any)  { s.Logger.Warn(msg, args...) }
func (s *slogLogger) Error(msg string, args ...any) { s.Logger.Error(msg, args...) }
func (s *slogLogger) Debug(msg string, args ...any) { s.Logger.Debug(msg, args...) }
func (s *slogLogger) Panic(msg string, args ...any) {
	s.Logger.Error(msg, args...)
	panic(msg)
}
func (s *slogLogger) Infof(format string, args ...any)  { s.Logger.Info(fmt.Sprintf(format, args...)) }
func (s *slogLogger) Warnf(format string, args ...any)  { s.Logger.Warn(fmt.Sprintf(format, args...)) }
func (s *slogLogger) Errorf(format string, args ...any) { s.Logger.Error(fmt.Sprintf(format, args...)) }
func (s *slogLogger) Debugf(format string, args ...any) { s.Logger.Debug(fmt.Sprintf(format, args...)) }
func (s *slogLogger) Panicf(format string, args ...any) {
	s.Logger.Error(fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
}

func (s *slogLogger) With(args ...any) Logger {
	return &slogLogger{s.Logger.With(args...), s.closer, s.logName}
}

func (s *slogLogger) Named(name string) Logger {
	//return &slogLogger{s.Logger.WithGroup(name), s.closer}
	newName := s.logName + "." + name
	return &slogLogger{
		Logger:  s.Logger.With("logger", newName),
		closer:  s.closer,
		logName: newName,
	}
}

// AddCallerSkip is NOT supported in slog.
// We return self and document the limitation.
func (s *slogLogger) AddCallerSkip(skip int) Logger {
	// Optionally: log a warning on first call, or panic in debug mode
	return s
}

func (s *slogLogger) Close() {
	if s.closer != nil {
		err := s.closer.Close()
		if err != nil {
			return
		}
	}
}
