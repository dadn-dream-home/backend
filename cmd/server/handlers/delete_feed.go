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

type ErrFeedNotExists struct{}

func (e ErrFeedNotExists) Error() string {
	return "feed not exists"
}

func (e ErrFeedNotExists) GRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

type DeleteFeedHandler struct {
	state.State
}

func (h DeleteFeedHandler) DeleteFeed(ctx context.Context, req *pb.DeleteFeedRequest) (*pb.DeleteFeedResponse, error) {
	log := telemetry.GetLogger(ctx).WithField("feed_id", req.Id)
	ctx = telemetry.ContextWithLogger(ctx, log)

	err := h.deleteFeedInDatabase(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.DeleteFeedResponse{}, nil
}

func (h DeleteFeedHandler) deleteFeedInDatabase(ctx context.Context, feedId string) error {
	log := telemetry.GetLogger(ctx)

	log.Debugf("deleting feed from database")

	res, err := h.DB().Exec("DELETE FROM feeds WHERE id = ?", feedId)
	if err != nil {
		return fmt.Errorf("error deleting feed from database: %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		panic(fmt.Errorf("database driver not support rows affected to check if feed exists: %w", err))
	} else if n == 0 {
		return ErrFeedNotExists{}
	}

	log.Infof("deleted feed %s from database successfully", feedId)

	return nil
}
