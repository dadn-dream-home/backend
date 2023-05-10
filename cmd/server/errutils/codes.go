package errutils

import (
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
)

func Internal(err error) error {
	s := status.New(codes.Internal, "")
	s.WithDetails(&epb.DebugInfo{
		StackEntries: []string{string(debug.Stack())},
		Detail:       err.Error(),
	})
	return s.Err()
}

func NotFound(resources ...*epb.ResourceInfo) error {
	s := status.New(codes.NotFound, resources[0].Description)
	for _, r := range resources {
		s.WithDetails(r)
	}
	return s.Err()
}

func AlreadyExists(resources ...*epb.ResourceInfo) error {
	s := status.New(codes.AlreadyExists, resources[0].Description)
	for _, r := range resources {
		s.WithDetails(r)
	}
	return s.Err()
}

func InvalidArgument(violations ...*epb.BadRequest_FieldViolation) error {
	s := status.New(codes.InvalidArgument, violations[0].Description)
	br := &epb.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, violations...)
	return s.Err()
}
