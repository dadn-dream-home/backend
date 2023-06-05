package state

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state/topic"
)

type State interface {
	DatabaseHooker() DatabaseHooker
	MQTTSubscriber() MQTTSubscriber
	FeedValuesLogger() FeedValuesLogger
	Repository() Repository
}

type DatabaseHooker interface {
	Subscribe(
		callback func(t topic.Topic, rowid int64) error,
		topics ...topic.Topic,
	) (unsubscribe func(), errCh <-chan error)
}

type MQTTSubscriber interface {
	Serve(ctx context.Context) error
	Publish(ctx context.Context, feedID string, value []byte) error
}

type FeedValuesLogger interface {
	Serve(ctx context.Context) error
}

type Repository interface {
	FeedRepository
	FeedValueRepository
	NotificationRepository
	ConfigRepository
}

type FeedRepository interface {
	CreateFeed(ctx context.Context, feed *pb.Feed) error
	GetFeed(ctx context.Context, feedID string) (*pb.Feed, error)
	GetFeedByRowID(ctx context.Context, rowID int64) (*pb.Feed, error)
	GetDeletedFeedByRowID(ctx context.Context, rowID int64) (*pb.Feed, error)
	ListFeeds(ctx context.Context) ([]*pb.Feed, error)
	DeleteFeed(ctx context.Context, feedID string) error
}

type FeedValueRepository interface {
	InsertFeedValue(ctx context.Context, feedID string, value []byte) error
	GetFeedLatestValue(ctx context.Context, feedID string) ([]byte, error)
	GetFeedValueByRowID(ctx context.Context, rowID int64) (string, []byte, error)
}

type NotificationRepository interface {
	InsertNotification(ctx context.Context, notification *pb.Notification) error
	GetLatestNotification(ctx context.Context, feedID string) (*pb.Notification, error)
	GetNotificationByRowID(ctx context.Context, rowID int64) (*pb.Notification, error)
}

type ConfigRepository interface {
	GetFeedConfig(ctx context.Context, feedID string) (*pb.Config, error)
	UpdateFeedConfig(ctx context.Context, config *pb.Config) error
}
