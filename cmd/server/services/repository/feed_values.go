package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type feedValueRepository struct {
	baseRepository
}

func (r feedValueRepository) InsertFeedValue(ctx context.Context, feedId string, value []byte) error {
	log := telemetry.GetLogger(ctx)

	if _, err := r.conn.ExecContext(
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
	if err := r.conn.QueryRowContext(ctx,
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

func (r feedValueRepository) GetFeedValueByRowID(ctx context.Context, rowID int64) ([]byte, error) {
	log := telemetry.GetLogger(ctx)

	var value []byte
	if err := r.conn.QueryRowContext(ctx,
		"SELECT value FROM feed_values WHERE rowid = ?",
		rowID,
	).Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errutils.NotFound(ctx, &errdetails.ResourceInfo{
				ResourceType: "FeedValue",
				ResourceName: fmt.Sprintf("rowid %d", rowID),
				Description:  fmt.Sprintf("FeedValue with rowid %d not found", rowID),
			})
		}

		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying feed value from database: %w", err))
	}

	log.Info("got feed value from database successfully")

	return value, nil
}
