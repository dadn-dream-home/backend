package telemetry

import (
	"context"

	"go.uber.org/zap"
)

type contextKey int

const (
	LoggerKey contextKey = iota
	RequestIdKey
)

func GetLogger(ctx context.Context) *zap.Logger {
	return ctx.Value(LoggerKey).(*zap.Logger)
}

func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func GetRequestId(ctx context.Context) string {
	return ctx.Value(RequestIdKey).(string)
}

func ContextWithRequestId(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, RequestIdKey, rid)
}
