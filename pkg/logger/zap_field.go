package logger

import (
	"time"

	"go.uber.org/zap"
)

// 将 zap.String / zap.Int 等调用替换为 logger.String / logger.Int 等封装方法
// 例如：
//
//	zap.String("key", "value")  →  logger.String("key", "value")
//	zap.Int("count", 1)         →  logger.Int("count", 1)
//
// 保持其余逻辑不变，仅替换 zap 字段构造器为 logger 提供的同名方法。

type Field = zap.Field

func String(key, val string) Field {
	return zap.String(key, val)
}
func Int(key string, val int) Field {
	return zap.Int(key, val)
}
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}
func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}
func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}
func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}
func Any(key string, val any) Field {
	return zap.Any(key, val)
}
func Error(err error) Field {
	return zap.Error(err)
}
