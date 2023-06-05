package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type feedValueRepository struct {
	baseRepository
}

func (r feedValueRepository) InsertFeedValue(ctx context.Context, feedId string, value []byte) error {
	log := telemetry.GetLogger(ctx)

	if _, err := r.db.ExecContext(
		ctx,
		"INSERT INTO feed_values (feed_id, value) VALUES (?, ?)",
		feedId, value,
	); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting feed value into database: %w", err))
	}

	log.Info("inserted feed value into database successfully")

	return nil
}

func (r feedValueRepository) GetFeedLatestValue(ctx context.Context, feedId string) ([]byte, error) {
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

func (r feedValueRepository) GetFeedValueByRowID(ctx context.Context, rowID int64) (string, []byte, error) {
	log := telemetry.GetLogger(ctx)

	var feedID string
	var value []byte
	if err := r.db.QueryRow(
		"SELECT feed_id, value FROM feed_values WHERE rowid = ?",
		rowID,
	).Scan(&feedID, &value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, errutils.NotFound(ctx, &errdetails.ResourceInfo{
				ResourceType: "FeedValue",
				ResourceName: fmt.Sprintf("rowid %d", rowID),
				Description:  fmt.Sprintf("FeedValue with rowid %d not found", rowID),
			})
		}

		return "", nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying feed value from database: %w", err))
	}

	log.Info("got feed value from database successfully")

	return feedID, value, nil
}

func (r feedValueRepository) ListActivities(ctx context.Context) ([]*pb.Activity, error) {
	types := []int32{int32(pb.FeedType_LIGHT)}
	typesCommaSeparated := ""
	for i, t := range types {
		if i != 0 {
			typesCommaSeparated += ", "
		}
		typesCommaSeparated += fmt.Sprintf("%d", t)
	}

	rows, err := r.db.Query(fmt.Sprintf(`
		SELECT id, type, value, time
		FROM feeds, feed_values
		WHERE id = feed_id AND type in (%s)
		ORDER BY time DESC
	`, typesCommaSeparated))

	if err != nil {
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying feed values from database: %w", err))
	}

	defer rows.Close()

	var activities []*pb.Activity
	for rows.Next() {
		var activity pb.Activity
		var value []byte
		var t time.Time
		if err := rows.Scan(&activity.Feed.Id, &activity.Feed.Type, &value, &t); err != nil {
			return nil, errutils.Internal(ctx, fmt.Errorf(
				"error scanning feed value from database: %w", err))
		}
		activity.State = string(value) != "0"
		activity.Timestamp = timestamppb.New(t)
		activities = append(activities, &activity)
	}

	return activities, nil
}
