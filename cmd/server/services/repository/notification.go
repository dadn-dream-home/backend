package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/telemetry"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type notificationRepository struct {
	baseRepository
}

func (r notificationRepository) InsertNotification(ctx context.Context, notification *pb.Notification) error {
	log := telemetry.GetLogger(ctx)

	if _, err := r.conn.ExecContext(ctx,
		"INSERT INTO notifications (feed_id, message) VALUES (?, ?)",
		notification.Feed.Id, notification.Message,
	); err != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"error inserting notification into database: %w", err))
	}

	log.Info("inserted notification into database successfully")

	return nil
}

func (r notificationRepository) GetLatestNotification(ctx context.Context, feedId string) (*pb.Notification, error) {
	log := telemetry.GetLogger(ctx)

	var notification pb.Notification
	if err := r.conn.QueryRowContext(ctx,
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
