package services

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

type repository struct {
	state.State
	db *sql.DB
}

type ErrFeedAlreadyExisted struct {
	Id string
}

type ErrFeedNotExists struct{}

func NewRepository(ctx context.Context, state state.State) (state.Repository, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return repository{
		State: state,
		db:    db,
	}, nil
}

func (r repository) CreateFeed(ctx context.Context, feed *pb.Feed) error {
	log := telemetry.GetLogger(ctx)

	log.Debugf("inserting feed into database")

	if res, err := r.db.Exec("INSERT OR IGNORE INTO feeds (id, type) VALUES (?, ?)", feed.Id, feed.Type); err != nil {
		return fmt.Errorf("error inserting feed into database: %w", err)
	} else if n, _ := res.RowsAffected(); n == 0 {
		return ErrFeedAlreadyExisted{}
	}

	log.Infof("inserted feed into database successfully")

	return nil
}

func (r repository) DeleteFeed(ctx context.Context, feedId string) error {
	log := telemetry.GetLogger(ctx)

	log.Debugf("deleting feed from database")

	res, err := r.db.Exec("DELETE FROM feeds WHERE id = ?", feedId)
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

func (r repository) GetFeed(ctx context.Context, feedId string) (*pb.Feed, error) {
	log := telemetry.GetLogger(ctx)

	log.Debugf("getting feed from database")

	var feed pb.Feed
	if err := r.db.QueryRow("SELECT id, type FROM feeds WHERE id = ?", feedId).Scan(&feed.Id, &feed.Type); err != nil {
		return nil, fmt.Errorf("error getting feed from database: %w", err)
	}

	log.Infof("got feed %s from database successfully", feedId)

	return &feed, nil
}

func (r repository) ListFeeds(ctx context.Context) ([]*pb.Feed, error) {
	log := telemetry.GetLogger(ctx)

	log.Debugf("listing feeds from database")

	rows, err := r.db.Query("SELECT id, type FROM feeds")
	if err != nil {
		return nil, fmt.Errorf("error listing feeds from database: %w", err)
	}
	defer rows.Close()

	feeds := make([]*pb.Feed, 0)
	for rows.Next() {
		var feed pb.Feed
		if err := rows.Scan(&feed.Id, &feed.Type); err != nil {
			return nil, fmt.Errorf("error scanning feed from database: %w", err)
		}
		feeds = append(feeds, &feed)
	}

	log.Infof("listed %d feeds from database successfully", len(feeds))

	return feeds, nil
}

func (r repository) InsertFeedValue(ctx context.Context, feedId string, value []byte) error {
	log := telemetry.GetLogger(ctx)

	log.Debugf("inserting feed value into database")

	if _, err := r.db.Exec(
		"INSERT INTO feed_values (feed_id, value) VALUES (?, ?)",
		feedId, value,
	); err != nil {
		return fmt.Errorf("error inserting feed value into database: %w", err)
	}

	log.Infof("inserted feed value into database successfully")

	return nil
}

func (r repository) GetFeedLatestValue(ctx context.Context, feedId string) ([]byte, error) {
	log := telemetry.GetLogger(ctx)

	log.Debugf("getting feed latest value from database")

	var value []byte
	if err := r.db.QueryRow(
		"SELECT value FROM feed_values WHERE feed_id = ? ORDER BY time DESC LIMIT 1",
		feedId,
	).Scan(&value); err != nil {
		return nil, fmt.Errorf("error getting feed latest value from database: %w", err)
	}

	log.Infof("got feed latest value from database successfully")

	return value, nil
}

func (e ErrFeedAlreadyExisted) Error() string {
	return fmt.Sprintf("feed %s already existed", e.Id)
}

func (e ErrFeedAlreadyExisted) GRPCStatus() *status.Status {
	return status.New(codes.AlreadyExists, e.Error())
}

func (e ErrFeedNotExists) Error() string {
	return "feed not exists"
}

func (e ErrFeedNotExists) GRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}
