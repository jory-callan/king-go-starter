package database

import (
	"context"
	"errors"
	"fmt"
	"king-starter/pkg/logger"
	"strings"
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
	switch strings.ToLower(level) {
	case "silent":
		gormLevel = gormlogger.Silent
	case "error":
		gormLevel = gormlogger.Error
	case "warn":
		gormLevel = gormlogger.Warn
	case "info":
		gormLevel = gormlogger.Info
	default:
		gormLevel = gormlogger.Info
	}

	return &gormLogger{
		logger:        logger,
		level:         gormLevel,
		slowThreshold: slowThreshold,
	}
}

func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.level = level
	return &newLogger
}
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Info {
		l.logger.Info(fmt.Sprintf(msg, data...))
	}
}
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Warn {
		l.logger.Warn(fmt.Sprintf(msg, data...))
	}
}
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level >= gormlogger.Error {
		l.logger.Error(fmt.Sprintf(msg, data...))
	}
}
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 如果日志级别低于Silent，则不输出日志
	if l.level <= gormlogger.Silent {
		return
	}
	// 计算执行时间
	elapsed := time.Since(begin)
	// 获取SQL和行数
	sql, rows := fc()
	// 构建日志字段
	fields := []zap.Field{
		zap.Duration("cost", elapsed),
		zap.String("sql", sql),
		zap.Int64("rows", rows),
	}
	// 根据错误和执行时间输出日志
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.level >= gormlogger.Error:
		l.logger.Error("gorm error", append(fields, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.slowThreshold > 0 && l.level >= gormlogger.Warn:
		l.logger.Warn("slow query", append(fields, zap.Duration("threshold", l.slowThreshold))...)
	case l.level >= gormlogger.Info:
		l.logger.Debug("gorm query", fields...)
	}
}
