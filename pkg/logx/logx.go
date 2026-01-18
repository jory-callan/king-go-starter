package logx

import (
	"sync"
)

var (
	globalLogger Logger
	once         sync.Once
)

// Logger is the unified interface for all log operations.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Panic(msg string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
	With(args ...any) Logger
	Named(name string) Logger
	AddCallerSkip(skip int) Logger
	Close()
}

// NewZap sets the global logger and returns it (for DI compatibility).
func NewZap(cfg *LoggerConfig) (Logger, error) {
	var l Logger
	var err error
	once.Do(func() {
		l, err = newZapLogger(cfg)
		if err == nil {
			globalLogger = l
		}
	})
	return l, err
}

// NewSlog sets the global logger and returns it (for DI compatibility).
func NewSlog(cfg *LoggerConfig) (Logger, error) {
	var l Logger
	var err error
	once.Do(func() {
		l, err = newSlogLogger(cfg)
		if err == nil {
			globalLogger = l
		}
	})
	return l, err
}

// Info logs at Info level.
func Info(msg string, args ...any) { globalLogger.Info(msg, args...) }

// Warn logs at Warn level.
func Warn(msg string, args ...any) { globalLogger.Warn(msg, args...) }

// Error logs at Error level.
func Error(msg string, args ...any) { globalLogger.Error(msg, args...) }

// Debug logs at Debug level.
func Debug(msg string, args ...any) { globalLogger.Debug(msg, args...) }

// With adds key-value pairs to the logger context.
func With(args ...any) Logger { return globalLogger.With(args...) }

// Named adds a name to the logger (e.g., module name).
func Named(name string) Logger { return globalLogger.Named(name) }

// AddCallerSkip increases the caller skip depth.
func AddCallerSkip(skip int) Logger { return globalLogger.AddCallerSkip(skip) }

// Close flushes any buffered logs.
func Close() { globalLogger.Close() }
