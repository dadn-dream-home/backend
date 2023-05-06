package handlers

import (
	"context"
	"database/sql"
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

func (h CreateFeedHandler) CreateFeed(ctx context.Context, req *pb.CreateFeedRequest) (res *pb.CreateFeedResponse, err error) {
	log := telemetry.GetLogger(ctx)
	log = log.WithField("feed_id", req.Id)
	log = log.WithField("feed_type", req.Type)
	ctx = telemetry.ContextWithLogger(ctx, log)

	tx, err := h.createFeedInDatabase(ctx, req.Id, req.Type)
	if err != nil {
		log.WithError(err).Errorf("error creating feed in database")
		return nil, err
	}
	defer commitOrRollback(tx, err)

	err = h.PubSub().AddFeed(ctx, req.Id)
	if err != nil {
		log.WithError(err).Errorf("error adding feed to pubsub")
		return nil, err
	}

	return &pb.CreateFeedResponse{
		Id:   req.Id,
		Type: req.Type,
	}, nil
}

func (h CreateFeedHandler) createFeedInDatabase(ctx context.Context, id string, ty pb.FeedType) (*sql.Tx, error) {
	log := telemetry.GetLogger(ctx)

	log.Debugf("inserting feed into database")

	tx, err := h.DB().Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	if res, err := tx.Exec("INSERT OR IGNORE INTO feeds (id, type) VALUES (?, ?)", id, ty); err != nil {
		return nil, fmt.Errorf("error inserting feed into database: %w", err)
	} else if n, _ := res.RowsAffected(); n == 0 {
		return nil, ErrFeedAlreadyExisted{}
	}

	log.Infof("inserted feed into database successfully")

	return tx, nil
}

func (e ErrFeedAlreadyExisted) Error() string {
	return fmt.Sprintf("feed %s already existed", e.Id)
}

func (e ErrFeedAlreadyExisted) GRPCStatus() *status.Status {
	return status.New(codes.AlreadyExists, e.Error())
}
