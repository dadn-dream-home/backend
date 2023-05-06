package handlers

import (
	"context"
	"fmt"

	pb "github.com/dadn-dream-home/x/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type CreateFeedHandler struct {
	state.State
}

type ErrFeedAlreadyExisted struct {
	Id string
}

func (h CreateFeedHandler) CreateFeed(ctx context.Context, req *pb.CreateFeedRequest) (*pb.CreateFeedResponse, error) {
	log := telemetry.GetLogger(ctx)
	log = log.WithField("feed_id", req.Id)
	log = log.WithField("feed_type", req.Type)
	ctx = telemetry.ContextWithLogger(ctx, log)

	err := h.createFeedInDatabase(ctx, req.Id, req.Type)
	if err != nil {
		log.WithError(err).Errorf("error creating feed in database")
		return nil, err
	}

	return &pb.CreateFeedResponse{
		Id:   req.Id,
		Type: req.Type,
	}, nil
}

func (h CreateFeedHandler) createFeedInDatabase(ctx context.Context, id string, ty pb.FeedType) error {
	log := telemetry.GetLogger(ctx)

	log.Debugf("inserting feed into database")

	if res, err := h.DB().Exec("INSERT OR IGNORE INTO feeds (id, type) VALUES (?, ?)", id, ty); err != nil {
		return fmt.Errorf("error inserting feed into database: %w", err)
	} else if n, _ := res.RowsAffected(); n == 0 {
		return ErrFeedAlreadyExisted{}
	}

	log.Infof("inserted feed into database successfully")

	return nil
}

func (e ErrFeedAlreadyExisted) Error() string {
	return fmt.Sprintf("feed %s already existed", e.Id)
}

func (e ErrFeedAlreadyExisted) GRPCStatus() *status.Status {
	return status.New(codes.AlreadyExists, e.Error())
}
