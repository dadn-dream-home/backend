package state

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state/topic"
)

type State interface {
	DatabaseListener() DatabaseListener
	PubSubValues() PubSubValues
	PubSubFeeds() PubSubFeeds
	Repository() Repository
	Notifier() Notifier
}

type DatabaseListener interface {
	Subscribe(
		topic topic.Topic,
		callback func(t topic.Topic, rowid int64) error,
	) (unsubscribe func(), errCh <-chan error)
}

type PubSubValues interface {
	Serve(ctx context.Context) error

	// If err != nil, feedID doesn't exist
	Subscribe(ctx context.Context, feedID string) (<-chan []byte, error)

	Unsubscribe(ctx context.Context, feedID string, ch <-chan []byte)
	Publish(ctx context.Context, feedID string, value []byte) error
}

type PubSubFeeds interface {
	Subscribe(ctx context.Context) (<-chan *pb.FeedsChange, error)
	Unsubscribe(ctx context.Context, ch <-chan *pb.FeedsChange)
	CreateFeed(ctx context.Context, feed *pb.Feed) error
	DeleteFeed(ctx context.Context, feedID string) error
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
	ListFeeds(ctx context.Context) ([]*pb.Feed, error)
	DeleteFeed(ctx context.Context, feedID string) error
}

type FeedValueRepository interface {
	InsertFeedValue(ctx context.Context, feedID string, value []byte) error
	GetFeedLatestValue(ctx context.Context, feedID string) ([]byte, error)
	GetFeedValueByRowID(ctx context.Context, rowID int64) ([]byte, error)
}

type NotificationRepository interface {
	InsertNotification(ctx context.Context, notification *pb.Notification) error
	GetLatestNotification(ctx context.Context, feedID string) (*pb.Notification, error)
}

type ConfigRepository interface {
	GetFeedConfig(ctx context.Context, feedID string) (*pb.Config, error)
	UpdateFeedConfig(ctx context.Context, config *pb.Config) error
}

type Notifier interface {
	Serve(ctx context.Context) error
	Subscribe(ctx context.Context) (<-chan *pb.Notification, error)
	Unsubscribe(ctx context.Context, ch <-chan *pb.Notification)
}
