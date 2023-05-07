package telemetry

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	LoggerKey contextKey = iota
	RequestIdKey
)

func GetLogger(ctx context.Context) *logrus.Entry {
	return ctx.Value(LoggerKey).(*logrus.Entry)
}

func ContextWithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func GetRequestId(ctx context.Context) string {
	return ctx.Value(RequestIdKey).(string)
}

func ContextWithRequestId(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, RequestIdKey, rid)
}
