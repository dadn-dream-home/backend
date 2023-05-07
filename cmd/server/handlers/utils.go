package handlers

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

func commitOrRollback(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error rolling back transaction: %w", rollbackErr)
		}
		return err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("error committing transaction: %w", commitErr)
	}
	return nil
}

func getFeedFromDatabase(ctx context.Context, s state.State, id string) (*pb.Feed, error) {
	log := telemetry.GetLogger(ctx).WithField("feed_id", id)

	log.Debugf("fetching feed type from database")

	var feed pb.Feed
	err := s.DB().QueryRow("SELECT id, type FROM feeds WHERE id = ?", id).Scan(&feed.Id, &feed.Type)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed from database: %w", err)
	}

	log.Infof("fetched feed type from database successfully")

	return &feed, nil
}
