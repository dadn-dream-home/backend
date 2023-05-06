package handlers

import (
	"context"
	"database/sql"
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

	tx, err := h.deleteFeedInDatabase(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	defer commitOrRollback(tx, err)

	err = h.PubSub().RemoveFeed(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteFeedResponse{}, nil
}

func (h DeleteFeedHandler) deleteFeedInDatabase(ctx context.Context, feedId string) (*sql.Tx, error) {
	log := telemetry.GetLogger(ctx)

	log.Debugf("deleting feed from database")

	tx, err := h.DB().Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	res, err := tx.Exec("DELETE FROM feeds WHERE id = ?", feedId)
	if err != nil {
		return nil, fmt.Errorf("error deleting feed from database: %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		panic(fmt.Errorf("database driver not support rows affected to check if feed exists: %w", err))
	} else if n == 0 {
		return nil, ErrFeedNotExists{}
	}

	log.Infof("deleted feed %s from database successfully", feedId)

	return tx, nil
}
