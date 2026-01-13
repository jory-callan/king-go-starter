package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestZapOrigin(t *testing.T) {
	log, _ := zap.NewDevelopment()
	log.Info("hello world")
	log.Warn("hello world")
	log.Error("hello world")

	log1 := log.Named("database")
	log1.Info("hello world")
	log1.Warn("hello world")
	log1.Error("hello world")

	log2 := log.With(
		zap.String("module", "redis"),
		zap.String("component", "instance"),
	)
	log2.Info("hello world")
	log2.Warn("hello world")
	log2.Error("hello world")

	log10, _ := zap.NewProduction()
	log10.Info("hello world")
	log10.Warn("hello world")
	log10.Error("hello world")

	log11 := log10.Named("database")
	log11.Info("hello world")
	log11.Warn("hello world")
	log11.Error("hello world")

	log12 := log10.With(
		zap.String("module", "redis"),
		zap.String("component", "instance"),
	)
	log12.Info("hello world")
	log12.Warn("hello world")
	log12.Error("hello world")
}
