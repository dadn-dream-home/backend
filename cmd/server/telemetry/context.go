package telemetry

import (
	"context"

	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	LoggerKey contextKey = iota
)

func GetLogger(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(LoggerKey).(logrus.FieldLogger)
}

func ContextWithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}
