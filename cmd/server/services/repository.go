package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/errutils"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type repository struct {
	state.State
	db *sql.DB
}

func NewRepository(ctx context.Context, state state.State, db *sql.DB) state.Repository {
	return repository{
		State: state,
		db:    db,
	}
}

func (r repository) CreateFeed(ctx context.Context, feed *pb.Feed) error {
	log := telemetry.GetLogger(ctx)

	if res, err := r.db.Exec("INSERT OR IGNORE INTO feeds (id, type) VALUES (?, ?)", feed.Id, feed.Type); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting feed into database: %w", err))
	} else if n, _ := res.RowsAffected(); n == 0 {
		return errutils.AlreadyExists(ctx, &errdetails.ResourceInfo{
			ResourceType: "Feed",
			ResourceName: feed.Id,
			Description:  fmt.Sprintf("Feed '%s' already exists", feed.Id),
		})
	}

	log.Info("inserted feed into database successfully")

	return nil
}

func (r repository) DeleteFeed(ctx context.Context, feedId string) error {
	log := telemetry.GetLogger(ctx)

	res, err := r.db.Exec("DELETE FROM feeds WHERE id = ?", feedId)
	if err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error deleting feed from database: %w", err))
	}

	if n, err := res.RowsAffected(); err != nil {
		log.Fatal("database driver not support rows affected to check if feed exists", zap.Error(err))
	} else if n == 0 {
		return errutils.NotFound(ctx, &errdetails.ResourceInfo{
			ResourceType: "Feed",
			ResourceName: feedId,
			Description:  fmt.Sprintf("Feed '%s' not found", feedId),
		})
	}

	log.Info("deleted feed from database successfully", zap.String("feed.id", feedId))

	return nil
}

func (r repository) GetFeed(ctx context.Context, feedId string) (*pb.Feed, error) {
	log := telemetry.GetLogger(ctx)

	var feed pb.Feed
	if err := r.db.QueryRow("SELECT id, type FROM feeds WHERE id = ?", feedId).Scan(&feed.Id, &feed.Type); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutils.NotFound(ctx, &errdetails.ResourceInfo{
				ResourceType: "Feed",
				ResourceName: feedId,
				Description:  fmt.Sprintf("Feed '%s' not found", feedId),
			})
		}
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying feed from database: %w", err))
	}

	log.Info("got feed %s from database successfully", zap.String("feed.id", feedId))

	return &feed, nil
}

func (r repository) ListFeeds(ctx context.Context) ([]*pb.Feed, error) {
	log := telemetry.GetLogger(ctx)

	rows, err := r.db.Query("SELECT id, type FROM feeds")
	if err != nil {
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying feeds from database: %w", err))
	}
	defer rows.Close()

	feeds := make([]*pb.Feed, 0)
	for rows.Next() {
		var feed pb.Feed
		if err := rows.Scan(&feed.Id, &feed.Type); err != nil {
			return nil, errutils.Internal(ctx, fmt.Errorf(
				"error scanning feed from database: %w", err))
		}
		feeds = append(feeds, &feed)
	}

	log.Info("listed feeds from database successfully", zap.Int("len", len(feeds)))

	return feeds, nil
}

func (r repository) InsertFeedValue(ctx context.Context, feedId string, value []byte) error {
	log := telemetry.GetLogger(ctx)

	if _, err := r.db.Exec(
		"INSERT INTO feed_values (feed_id, value) VALUES (?, ?)",
		feedId, value,
	); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting feed value into database: %w", err))
	}

	log.Info("inserted feed value into database successfully")

	return nil
}

func (r repository) GetFeedLatestValue(ctx context.Context, feedId string) ([]byte, error) {
	log := telemetry.GetLogger(ctx)

	var value []byte
	if err := r.db.QueryRow(
		"SELECT value FROM feed_values WHERE feed_id = ? ORDER BY time DESC LIMIT 1",
		feedId,
	).Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying feed latest value from database: %w", err))
	}

	log.Info("got feed latest value from database successfully")

	return value, nil
}

func (r repository) InsertNotification(ctx context.Context, notification *pb.Notification) error {
	log := telemetry.GetLogger(ctx)

	if _, err := r.db.Exec(
		"INSERT INTO notifications (feed_id, message) VALUES (?, ?)",
		notification.Feed.Id, notification.Message,
	); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting notification into database: %w", err))
	}

	log.Info("inserted notification into database successfully")

	return nil
}

func (r repository) GetLatestNotification(ctx context.Context, feedId string) (*pb.Notification, error) {
	log := telemetry.GetLogger(ctx)

	var notification pb.Notification
	if err := r.db.QueryRow(
		"SELECT feed_id, message, time FROM notifications WHERE feed_id = ? ORDER BY time DESC LIMIT 1",
		feedId,
	).Scan(&notification.Feed.Id, &notification.Message, notification.Timestamp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying latest notification from database: %w", err))
	}

	log.Info("got latest notification from database successfully")

	return &notification, nil
}
