package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/dadn-dream-home/x/server/state"
)

type feedRepository struct {
	baseRepository
}

var _ state.FeedRepository = feedRepository{}

func (r feedRepository) CreateFeed(ctx context.Context, feed *pb.Feed) error {
	log := telemetry.GetLogger(ctx)

	if res, err := r.conn.ExecContext(ctx, "INSERT OR IGNORE INTO feeds (id, type) VALUES (?, ?)", feed.Id, feed.Type); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting feed into database: %w", err))
	} else if n, _ := res.RowsAffected(); n == 0 {
		return errutils.AlreadyExists(ctx, &errdetails.ResourceInfo{
			ResourceType: "Feed",
			ResourceName: feed.Id,
			Description:  fmt.Sprintf("Feed '%s' already exists", feed.Id),
		})
	}

	// create feed config
	switch feed.Type {
	case pb.FeedType_TEMPERATURE, pb.FeedType_HUMIDITY:
		if _, err := r.conn.ExecContext(ctx, "INSERT OR IGNORE INTO sensor_configs (feed_id) VALUES (?)", feed.Id); err != nil {
			return errutils.Internal(ctx, fmt.Errorf(
				"error inserting sensor config into database: %w", err))
		}
	case pb.FeedType_LIGHT:
		if _, err := r.conn.ExecContext(ctx, "INSERT OR IGNORE INTO actuator_configs (feed_id) VALUES (?)", feed.Id); err != nil {
			return errutils.Internal(ctx, fmt.Errorf(
				"error inserting actuator config into database: %w", err))
		}
	default:
		return errutils.Internal(ctx, errors.New("unhandled feed type"))
	}

	log.Info("inserted feed into database successfully")

	return nil
}

func (r feedRepository) DeleteFeed(ctx context.Context, feedId string) error {
	log := telemetry.GetLogger(ctx)

	res, err := r.conn.ExecContext(ctx, "DELETE FROM feeds WHERE id = ?", feedId)
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

func (r feedRepository) updateFeed(ctx context.Context, feed *pb.Feed) error {
	log := telemetry.GetLogger(ctx)

	if res, err := r.conn.ExecContext(ctx, "UPDATE feeds SET type = ? WHERE id = ?", feed.Type, feed.Id); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error updating feed in database: %w", err))
	} else if n, _ := res.RowsAffected(); n == 0 {
		return errutils.NotFound(ctx, &errdetails.ResourceInfo{
			ResourceType: "Feed",
			ResourceName: feed.Id,
			Description:  fmt.Sprintf("Feed '%s' not found", feed.Id),
		})
	}

	log.Info("updated feed in database successfully", zap.String("feed.id", feed.Id))

	return nil
}

func (r feedRepository) GetFeed(ctx context.Context, feedId string) (*pb.Feed, error) {
	log := telemetry.GetLogger(ctx)

	var feed pb.Feed
	if err := r.conn.QueryRowContext(ctx, "SELECT id, type FROM feeds WHERE id = ?", feedId).Scan(&feed.Id, &feed.Type); err != nil {
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

func (r feedRepository) ListFeeds(ctx context.Context) ([]*pb.Feed, error) {
	log := telemetry.GetLogger(ctx)

	rows, err := r.conn.QueryContext(ctx, "SELECT id, type FROM feeds")
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
