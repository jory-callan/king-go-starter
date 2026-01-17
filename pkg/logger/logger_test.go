package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := LoggerConfig{
			Level:      "debug",
			Format:     "text",
			Output:     "stdout",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     5,
			Compress:   false,
			CallerSkip: 1,
		}
		err := cfg.Validate()
		require.NoError(t, err)
		logger, err := New(&cfg)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		require.NoError(t, err)
		assert.NotNil(t, logger)
		defer logger.Sync()
	})

	t.Run("invalid log level", func(t *testing.T) {
		cfg := LoggerConfig{
			Level:  "invalid",
			Format: "json",
			Output: "stdout",
		}

		logger, err := New(&cfg)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		assert.Panics(t, func() {
			logger.Info("test message from default config")
		})
	})

	t.Run("invalid output", func(t *testing.T) {
		cfg := LoggerConfig{
			Level:  "info",
			Format: "json",
			Output: "invalid",
		}

		logger, err := New(&cfg)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		assert.Panics(t, func() {
			logger.Info("test message from default config")
		})
	})
}

func TestNewWithDefaultConfig(t *testing.T) {
	cfg := DefaultLoggerConfig()
	logger, err := New(&cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	assert.NotNil(t, logger)
	defer logger.Sync()

	// 测试日志输出
	logger.Info("test message from default config")
	logger.Debug("debug message")
}

func TestLoggerLoggingMethods(t *testing.T) {
	cfg := DefaultLoggerConfig()
	logger, err := New(&cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer logger.Sync()

	t.Run("Info", func(t *testing.T) {
		logger.Info("info message")
	})

	t.Run("InfoWithFields", func(t *testing.T) {
		logger.Info("info with fields",
			zap.String("key", "value"),
			zap.Int("count", 42),
		)
	})

	t.Run("Debug", func(t *testing.T) {
		logger.Debug("debug message")
	})

	t.Run("Warn", func(t *testing.T) {
		logger.Warn("warn message")
	})

	t.Run("Error", func(t *testing.T) {
		logger.Error("error message")
	})
}

func TestLoggerWithFileOutput(t *testing.T) {
	cfg := LoggerConfig{
		Level:    "debug",
		Format:   "json",
		Output:   "file",
		FilePath: "/tmp/test_logger.log",
		MaxSize:  1,
		MaxAge:   1,
	}

	logger, err := New(&cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	assert.NotNil(t, logger)
	defer logger.Sync()

	logger.Info("test file output message")
}

func TestLoggerSync(t *testing.T) {
	cfg := DefaultLoggerConfig()
	logger, err := New(&cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	assert.NotNil(t, logger)
	defer logger.Sync()

	// Note: Sync may return an error when syncing stdout/stderr
	// This is expected and not a critical error in most scenarios
	_ = logger.Sync()
}
