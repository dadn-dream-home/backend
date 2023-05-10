package state

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type State interface {
	PubSubValues() PubSubValues
	PubSubFeeds() PubSubFeeds
	Repository() Repository
	Notifier() Notifier
}

type PubSubValues interface {
	Subscribe(ctx context.Context, feedID string) (<-chan []byte, error)
	Unsubscribe(ctx context.Context, feedID string, ch <-chan []byte) error
	Publish(ctx context.Context, feedID string, value []byte) error
}

type PubSubFeeds interface {
	Subscribe(ctx context.Context) (ch <-chan *pb.FeedsChange, err error)
	Unsubscribe(ctx context.Context, ch <-chan *pb.FeedsChange) error
	CreateFeed(ctx context.Context, feed *pb.Feed) error
	DeleteFeed(ctx context.Context, feedID string) error
}

type Repository interface {
	CreateFeed(ctx context.Context, feed *pb.Feed) error
	GetFeed(ctx context.Context, feedID string) (*pb.Feed, error)
	ListFeeds(ctx context.Context) ([]*pb.Feed, error)
	DeleteFeed(ctx context.Context, feedID string) error
	InsertFeedValue(ctx context.Context, feedID string, value []byte) error
	GetFeedLatestValue(ctx context.Context, feedID string) ([]byte, error)
	InsertNotification(ctx context.Context, notification *pb.Notification) error
	GetLatestNotification(ctx context.Context, feedID string) (*pb.Notification, error)
}

type Notifier interface {
	Subscribe(ctx context.Context) (<-chan *pb.Notification, error)
	Unsubscribe(ctx context.Context, ch <-chan *pb.Notification) error
}
