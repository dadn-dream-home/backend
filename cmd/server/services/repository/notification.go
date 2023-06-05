package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type notificationRepository struct {
	baseRepository
}

func (r notificationRepository) InsertNotification(ctx context.Context, notification *pb.Notification) error {
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

func (r notificationRepository) GetLatestNotification(ctx context.Context, feedId string) (*pb.Notification, error) {
	log := telemetry.GetLogger(ctx)

	notification := pb.Notification{
		Feed: &pb.Feed{},
	}
	var t time.Time
	if err := r.db.QueryRow(
		"SELECT feed_id, message, time FROM notifications WHERE feed_id = ? ORDER BY time DESC LIMIT 1",
		feedId,
	).Scan(&notification.Feed.Id, &notification.Message, &t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying latest notification from database: %w", err))
	}
	notification.Timestamp = timestamppb.New(t)

	log.Info("got latest notification from database successfully")

	return &notification, nil
}

func (r notificationRepository) GetNotificationByRowID(ctx context.Context, rowID int64) (*pb.Notification, error) {
	log := telemetry.GetLogger(ctx)

	notification := pb.Notification{
		Feed: &pb.Feed{},
	}
	var t time.Time
	if err := r.db.QueryRow(
		"SELECT feed_id, message, time FROM notifications WHERE rowid = ?",
		rowID,
	).Scan(&notification.Feed.Id, &notification.Message, &t); err != nil {
		return nil, errutils.Internal(ctx, fmt.Errorf(
			"error querying notification from database: %w", err))
	}
	notification.Timestamp = timestamppb.New(t)

	log.Info("got notification from database successfully")

	return &notification, nil
}
