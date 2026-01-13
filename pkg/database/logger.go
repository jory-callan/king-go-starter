package database

import (
	"context"
	"fmt"
	"king-starter/pkg/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// gormLogger GORM日志适配器
type gormLogger struct {
	logger        *logger.Logger
	level         gormlogger.LogLevel
	slowThreshold time.Duration
}

// newGormLogger 创建GORM日志适配器
func newGormLogger(logger *logger.Logger, level string, slowThreshold time.Duration) gormlogger.Interface {
	var gormLevel gormlogger.LogLevel
	switch level {
	case "silent":
		gormLevel = gormlogger.Silent
	case "error":
		gormLevel = gormlogger.Error
	case "warn":
		gormLevel = gormlogger.Warn
	case "info":
		gormLevel = gormlogger.Info
	default:
		gormLevel = gormlogger.Warn
	}

	return &gormLogger{
		logger:        logger,
		level:         gormLevel,
		slowThreshold: slowThreshold,
	}
}

// LogMode 实现gormlogger.Interface
func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}

// Info 实现gormlogger.Interface
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Info {
		l.logger.Info(fmt.Sprintf(msg, data...))
	}
}

// Warn 实现gormlogger.Interface
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Warn {
		l.logger.Warn(fmt.Sprintf(msg, data...))
	}
}

// Error 实现gormlogger.Interface
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Error {
		l.logger.Error(fmt.Sprintf(msg, data...))
	}
}

// Trace 实现gormlogger.Interface
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 构建日志字段
	fields := []zap.Field{
		zap.Duration("cost", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}

	switch {
	case err != nil && err != gorm.ErrRecordNotFound && l.level >= gormlogger.Error:
		l.logger.Error("gorm error", append(fields, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.slowThreshold > 0 && l.level >= gormlogger.Warn:
		l.logger.Warn("slow query", append(fields, zap.Duration("threshold", l.slowThreshold))...)
	case l.level >= gormlogger.Info:
		l.logger.Debug("gorm query", fields...)
	}
}
