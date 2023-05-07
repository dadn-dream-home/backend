package state

import (
	"context"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type State interface {
	PubSubValues() PubSubValues
	PubSubFeeds() PubSubFeeds
	Repository() Repository
}

type PubSubValues interface {
	Subscribe(ctx context.Context, sid string, feed string, ch chan<- []byte) (done <-chan struct{}, err error)
	Unsubscribe(ctx context.Context, sid string, feed string) error
	Publish(ctx context.Context, sid string, feed string, value []byte) error
}

type PubSubFeeds interface {
	Subscribe(ctx context.Context, sid string, ch chan<- *pb.FeedsChange) (done <-chan struct{}, err error)
	Unsubscribe(ctx context.Context, sid string) error
	CreateFeed(ctx context.Context, sid string, feed *pb.Feed) error
	DeleteFeed(ctx context.Context, sid string, feed string) error
}

type Repository interface {
	CreateFeed(ctx context.Context, feed *pb.Feed) error
	GetFeed(ctx context.Context, id string) (*pb.Feed, error)
	ListFeeds(ctx context.Context) ([]*pb.Feed, error)
	DeleteFeed(ctx context.Context, id string) error
}
