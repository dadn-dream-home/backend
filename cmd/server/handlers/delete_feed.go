package handlers

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type DeleteFeedHandler struct {
	state.State
}

func (h DeleteFeedHandler) DeleteFeed(ctx context.Context, req *pb.DeleteFeedRequest) (*pb.DeleteFeedResponse, error) {
	log := telemetry.GetLogger(ctx)

	if req.Feed == nil {
		log.Error("feed is nil")
		return nil, status.Error(codes.InvalidArgument, "feed cannot be nil")
	}

	log = log.With(zap.String("feed.id", req.Feed.Id))
	ctx = telemetry.ContextWithLogger(ctx, log)

	err := h.PubSubFeeds().DeleteFeed(ctx, req.Feed.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteFeedResponse{}, nil
}
