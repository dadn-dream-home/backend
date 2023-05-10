package telemetry

import (
	"context"

	"go.uber.org/zap"
)

func InitLogger(ctx context.Context) context.Context {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return ContextWithLogger(ctx, log)
}
