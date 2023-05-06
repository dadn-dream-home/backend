package handlers

import (
	"context"
	"fmt"

	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListFeedsHandler struct {
	state.State
}

func (h ListFeedsHandler) ListFeeds(ctx context.Context, req *pb.ListFeedsRequest) (*pb.ListFeedsResponse, error) {
	feeds, err := h.getFeeds(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListFeedsResponse{Feeds: feeds}, nil
}

func (h ListFeedsHandler) getFeeds(ctx context.Context) (feeds []*pb.Feed, err error) {
	log := telemetry.GetLogger(ctx)

	log.Debugf("selecting feeds from database")

	rows, err := h.DB().Query("SELECT id, type FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("error querying feeds: %w", err)
	}
	defer rows.Close()

	log.Debugf("selected feeds from database successfully")

	for rows.Next() {
		var feed pb.Feed
		if err := rows.Scan(&feed.Id, &feed.Type); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		feeds = append(feeds, &feed)
	}

	log.WithField("count", len(feeds)).Infof("get feeds from database successfully")

	return feeds, nil
}
