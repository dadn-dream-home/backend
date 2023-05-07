package services

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type pubSubFeeds struct {
	sync.RWMutex

	state.State

	subscribers map[string]pubSubFeedsSubscriber
}

type pubSubFeedsSubscriber struct {
	id   string
	ch   chan<- *pb.FeedsChange
	done chan struct{}
}

func NewPubSubFeeds(ctx context.Context, state state.State) state.PubSubFeeds {
	return &pubSubFeeds{
		State:       state,
		subscribers: make(map[string]pubSubFeedsSubscriber),
	}
}

func (p *pubSubFeeds) Subscribe(ctx context.Context, id string, ch chan<- *pb.FeedsChange) (<-chan struct{}, error) {
	p.Lock()
	defer p.Unlock()

	log := telemetry.GetLogger(ctx).WithField("subscriber_id", id)

	p.subscribers[id] = pubSubFeedsSubscriber{
		id:   id,
		ch:   ch,
		done: make(chan struct{}, 1),
	}

	// send initial list of feeds
	go func() {
		feeds, err := p.Repository().ListFeeds(ctx)
		if err != nil {
			log.WithError(err).Fatalf("failed to list feeds")
		}

		ch <- &pb.FeedsChange{Added: feeds}
		if err != nil {
			log.WithError(err).Fatalf("failed to send initial list of feeds")
		}
	}()

	log.Infof("subscribed to pub sub feeds")

	return p.subscribers[id].done, nil
}

func (p *pubSubFeeds) Unsubscribe(ctx context.Context, id string) error {
	log := telemetry.GetLogger(ctx).WithField("subscriber_id", id)

	p.Lock()
	p.subscribers[id].done <- struct{}{}
	delete(p.subscribers, id)
	p.Unlock()

	log.Infof("unsubscribed from pub sub feeds")

	return nil
}

func (p *pubSubFeeds) CreateFeed(ctx context.Context, feed *pb.Feed) error {
	err := p.Repository().CreateFeed(ctx, feed)
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	p.RLock()
	for _, subscriber := range p.subscribers {
		subscriber.ch <- &pb.FeedsChange{
			Added: []*pb.Feed{feed},
		}
	}
	p.RUnlock()

	return nil
}

func (p *pubSubFeeds) DeleteFeed(ctx context.Context, feed string) error {
	err := p.Repository().DeleteFeed(ctx, feed)
	if err != nil {
		return fmt.Errorf("failed to delete feed: %w", err)
	}

	p.RLock()
	for _, subscriber := range p.subscribers {
		subscriber.ch <- &pb.FeedsChange{
			Removed: []string{feed},
		}
	}
	p.RUnlock()

	return nil
}
