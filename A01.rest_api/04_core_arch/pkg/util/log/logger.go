package log

import (
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// InitLogger initializes the global logger based on configuration
func InitLogger(level string, format string) {
	once.Do(func() {
		var config zap.Config

		if strings.ToLower(format) == "json" {
			config = zap.NewProductionConfig()
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// Set Log Level
		var zapLevel zapcore.Level
		switch strings.ToLower(level) {
		case "debug":
			zapLevel = zap.DebugLevel
		case "warn":
			zapLevel = zap.WarnLevel
		case "error":
			zapLevel = zap.ErrorLevel
		default:
			zapLevel = zap.InfoLevel
		}
		config.Level = zap.NewAtomicLevelAt(zapLevel)

		// Disable stacktrace for info/warn
		config.DisableStacktrace = true

		logger, err := config.Build()
		if err != nil {
			panic(err)
		}

		globalLogger = logger
		zap.ReplaceGlobals(logger)
	})
}

// L returns the global logger
func L() *zap.Logger {
	if globalLogger == nil {
		// Fallback if not initialized
		l, _ := zap.NewProduction()
		return l
	}
	return globalLogger
}

// Sync flushes any buffered log entries
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}
