package handlers

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type CreateFeedHandler struct {
	state.State
}

func (h CreateFeedHandler) CreateFeed(ctx context.Context, req *pb.CreateFeedRequest) (res *pb.CreateFeedResponse, err error) {
	log := telemetry.GetLogger(ctx).
		WithField("feed_id", req.Feed.Id).
		WithField("feed_type", req.Feed.Type)
	ctx = telemetry.ContextWithLogger(ctx, log)
	rid := telemetry.GetRequestId(ctx)

	err = h.PubSubFeeds().CreateFeed(ctx, rid, req.Feed)
	if err != nil {
		log.WithError(err).Errorf("error creating feed")
		return nil, err
	}

	return &pb.CreateFeedResponse{}, nil
}
