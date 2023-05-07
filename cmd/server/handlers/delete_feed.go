package handlers

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type DeleteFeedHandler struct {
	state.State
}

func (h DeleteFeedHandler) DeleteFeed(ctx context.Context, req *pb.DeleteFeedRequest) (*pb.DeleteFeedResponse, error) {
	log := telemetry.GetLogger(ctx).WithField("feed_id", req.Id)
	ctx = telemetry.ContextWithLogger(ctx, log)
	rid := telemetry.GetRequestId(ctx)

	err := h.PubSubFeeds().DeleteFeed(ctx, rid, req.Id)
	if err != nil {
		log.WithError(err).Errorf("error deleting feed from pubsub")
		return nil, err
	}

	return &pb.DeleteFeedResponse{}, nil
}
