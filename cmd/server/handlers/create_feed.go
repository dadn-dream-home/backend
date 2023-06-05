package handlers

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type CreateFeedHandler struct {
	state.State
}

func (h CreateFeedHandler) CreateFeed(ctx context.Context, req *pb.CreateFeedRequest) (res *pb.CreateFeedResponse, err error) {
	log := telemetry.GetLogger(ctx)

	if req.Feed == nil {
		return nil, errutils.InvalidArgument(ctx, &errdetails.BadRequest_FieldViolation{
			Field:       "feed",
			Description: "Request field 'feed' is nil, expected object",
		})
	}

	if req.Feed.Id == "" {
		return nil, errutils.InvalidArgument(ctx, &errdetails.BadRequest_FieldViolation{
			Field:       "feed.id",
			Description: "Request field 'feed.id' is empty, expected non-empty string",
		})
	}

	log = log.With(
		zap.String("feed.id", req.Feed.Id),
		zap.String("feed.type", req.Feed.Type.String()),
	)
	ctx = telemetry.ContextWithLogger(ctx, log)

	if err := h.Repository().CreateFeed(ctx, req.Feed); err != nil {
		return nil, err
	}

	return &pb.CreateFeedResponse{}, nil
}
