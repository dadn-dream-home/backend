package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type feedValueRepository struct {
	baseRepository
}

func (r feedValueRepository) InsertFeedValue(ctx context.Context, feedId string, value []byte) error {
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
