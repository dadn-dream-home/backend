package services

import (
	"context"
	"sync"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
)

type pubSubFeeds struct {
	sync.RWMutex

	state.State

	subscribers map[<-chan *pb.FeedsChange]chan<- *pb.FeedsChange
}

func NewPubSubFeeds(ctx context.Context, state state.State) state.PubSubFeeds {
	return &pubSubFeeds{
		State:       state,
		subscribers: make(map[<-chan *pb.FeedsChange]chan<- *pb.FeedsChange),
	}
}

func (p *pubSubFeeds) Subscribe(ctx context.Context) (<-chan *pb.FeedsChange, error) {
	p.Lock()
	defer p.Unlock()

	log := telemetry.GetLogger(ctx)

	ch := make(chan *pb.FeedsChange)

	p.subscribers[ch] = ch

	// send initial list of feeds
	go func() {
		feeds, err := p.Repository().ListFeeds(ctx)
		if err != nil {
			log.Fatal("failed to list feeds", zap.Error(err))
		}

		ch <- &pb.FeedsChange{Addeds: feeds}

		if err != nil {
			log.Fatal("failed to send initial list of feeds", zap.Error(err))
		}
	}()

	log.Info("subscribed to pub sub feeds")

	return ch, nil
}

func (p *pubSubFeeds) Unsubscribe(ctx context.Context, ch <-chan *pb.FeedsChange) error {
	p.Lock()
	defer p.Unlock()

	log := telemetry.GetLogger(ctx)

	p.subscribers[ch] <- nil
	delete(p.subscribers, ch)

	log.Info("unsubscribed from pub sub feeds")

	return nil
}

func (p *pubSubFeeds) CreateFeed(ctx context.Context, feed *pb.Feed) error {
	p.RLock()
	defer p.RUnlock()

	log := telemetry.GetLogger(ctx)

	if err := p.Repository().CreateFeed(ctx, feed); err != nil {
		return err
	}

	for _, subscriber := range p.subscribers {
		subscriber <- &pb.FeedsChange{
			Addeds: []*pb.Feed{feed},
		}
	}

	log.Info("created feed")

	return nil
}

func (p *pubSubFeeds) DeleteFeed(ctx context.Context, feedID string) error {
	p.RLock()
	defer p.RUnlock()

	log := telemetry.GetLogger(ctx)

	if err := p.Repository().DeleteFeed(ctx, feedID); err != nil {
		return err
	}

	for _, subscriber := range p.subscribers {
		subscriber <- &pb.FeedsChange{
			RemovedIDs: []string{feedID},
		}
	}

	log.Info("deleted feed")

	return nil
}
