package log

import (
	"context"

	"go.uber.org/zap"
)

type contextKey struct{}

var loggerKey = contextKey{}

// WithLogger returns a new context with the provided logger attached
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context, or the global logger if none exists
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return L()
}

// With creates a child logger with the provided fields and returns a new context containing it
func With(ctx context.Context, fields ...zap.Field) context.Context {
	logger := FromContext(ctx).With(fields...)
	return WithLogger(ctx, logger)
}
