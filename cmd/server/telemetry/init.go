package telemetry

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(ctx context.Context) context.Context {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log, err := config.Build()
	if err != nil {
		panic(err)
	}

	return ContextWithLogger(ctx, log)
}
