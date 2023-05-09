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
	Subscribe(ctx context.Context, sid string, feed string) (<-chan []byte, error)
	SubscribeAll(ctx context.Context, sid string) (<-chan map[string][]byte, error)
	Unsubscribe(ctx context.Context, sid string, feed string) error
	Publish(ctx context.Context, sid string, feed string, value []byte) error
}

type PubSubFeeds interface {
	Subscribe(ctx context.Context, sid string) (ch <-chan *pb.FeedsChange, err error)
	Unsubscribe(ctx context.Context, sid string) error
	CreateFeed(ctx context.Context, sid string, feed *pb.Feed) error
	DeleteFeed(ctx context.Context, sid string, feed string) error
}

type Repository interface {
	CreateFeed(ctx context.Context, feed *pb.Feed) error
	GetFeed(ctx context.Context, id string) (*pb.Feed, error)
	ListFeeds(ctx context.Context) ([]*pb.Feed, error)
	DeleteFeed(ctx context.Context, id string) error
	InsertFeedValue(ctx context.Context, feed string, value []byte) error
	GetFeedLatestValue(ctx context.Context, feed string) ([]byte, error)
	InsertNotification(ctx context.Context, notification *pb.Notification) error
	GetLatestNotification(ctx context.Context, feed string) (*pb.Notification, error)
}

type Notifier interface {
	Subscribe(ctx context.Context, sid string) (<-chan *pb.Notification, error)
	Unsubscribe(ctx context.Context, sid string) error
}
