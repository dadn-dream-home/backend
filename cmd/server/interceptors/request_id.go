package interceptors

import (
	"context"

	"github.com/dadn-dream-home/x/server/telemetry"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func RequestIdUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// keep the request id short
	rid := uuid.New().String()[:6]
	log := telemetry.GetLogger(ctx).WithField("rid", rid)
	ctx = telemetry.ContextWithLogger(ctx, log)

	return handler(ctx, req)
}

func RequestIdStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	// keep the request id short
	rid := uuid.New().String()[:6]
	log := telemetry.GetLogger(ctx).WithField("rid", rid)
	ctx = telemetry.ContextWithLogger(ctx, log)

	return handler(srv, &serverStream{ss, ctx})
}
