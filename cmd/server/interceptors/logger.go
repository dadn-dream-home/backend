package interceptors

import (
	"context"

	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/grpc"
)

func AddLoggerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// keep the request id short
	ctx = telemetry.InitLogger(ctx)
	return handler(ctx, req)
}

func AddLoggerStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	ctx = telemetry.InitLogger(ctx)
	return handler(srv, &serverStream{ss, ctx})
}
