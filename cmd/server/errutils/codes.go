package errutils

import (
	"context"
	"runtime/debug"
	"strings"

	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
)

func Internal(ctx context.Context, err error) error {
	log := telemetry.GetLogger(ctx)
	s := status.New(codes.Internal, "")
	s, _ = s.WithDetails(&epb.DebugInfo{
		StackEntries: strings.Split(string(debug.Stack()), "\n"),
		Detail:       err.Error(),
	})
	log.Error("internal error", zap.Error(err))
	return s.Err()
}

func NotFound(ctx context.Context, resources ...*epb.ResourceInfo) error {
	log := telemetry.GetLogger(ctx)
	s := status.New(codes.NotFound, resources[0].Description)
	for _, r := range resources {
		s.WithDetails(r)
	}
	log.Warn("resource not found", zap.Error(s.Err()))
	return s.Err()
}

func AlreadyExists(ctx context.Context, resources ...*epb.ResourceInfo) error {
	log := telemetry.GetLogger(ctx)
	s := status.New(codes.AlreadyExists, resources[0].Description)
	for _, r := range resources {
		s.WithDetails(r)
	}
	log.Warn("resource already exists", zap.Error(s.Err()))
	return s.Err()
}

func InvalidArgument(ctx context.Context, violations ...*epb.BadRequest_FieldViolation) error {
	log := telemetry.GetLogger(ctx)
	s := status.New(codes.InvalidArgument, violations[0].Description)
	br := &epb.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, violations...)
	log.Warn("invalid argument", zap.Error(s.Err()))
	return s.Err()
}
