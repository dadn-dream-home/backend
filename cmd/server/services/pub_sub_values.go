package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type pubSubValues struct {
	sync.RWMutex
	Ready chan struct{}

	state.State
	mqtt mqtt.Client

	// map feed -> rx -> tx
	subscribers map[string]map[<-chan []byte]chan<- []byte
}

func NewPubSubValues(ctx context.Context, state state.State, mqtt mqtt.Client) state.PubSubValues {
	ps := &pubSubValues{
		State:       state,
		mqtt:        mqtt,
		subscribers: make(map[string]map[<-chan []byte]chan<- []byte),
	}

	return ps
}

func (p *pubSubValues) Serve(ctx context.Context) error {
	log := telemetry.GetLogger(ctx)
	log = log.With(zap.String("service", "pubsub_values"))

	ch, err := p.PubSubFeeds().Subscribe(ctx)
	if err != nil {
		return err
	}

	log.Info("started streaming")

	for {
		select {
		case <-ctx.Done():
			p.PubSubFeeds().Unsubscribe(ctx, ch)

		case change := <-ch:
			if change == nil {
				log.Info("ended streaming")
				return nil
			}
			
			for _, feed := range change.Addeds {
				log := log.With(zap.String("feed.id", feed.Id))
				ctx = telemetry.ContextWithLogger(ctx, log)

				if err := p.addFeed(ctx, feed.Id); err != nil {
					return err
				}
				log.Info("added feed")
			}
			for _, feedID := range change.RemovedIDs {
				log := log.With(zap.String("feed.id", feedID))
				ctx = telemetry.ContextWithLogger(ctx, log)

				if err := p.removeFeed(ctx, feedID); err != nil {
					return err
				}
				log.Info("removed feed")
			}
		}
	}
}

func (p *pubSubValues) Subscribe(ctx context.Context, feedID string) (<-chan []byte, error) {
	p.Lock()
	defer p.Unlock()

	log := telemetry.GetLogger(ctx)

	if _, ok := p.subscribers[feedID]; !ok {
		return nil, errutils.NotFound(ctx, &errdetails.ResourceInfo{
			ResourceType: "Feed",
			ResourceName: feedID,
			Description:  fmt.Sprintf("Feed '%s' cannot be found", feedID),
		})
	}

	ch := make(chan []byte, 1)
	p.subscribers[feedID][ch] = ch

	log.Info("subscribed to pubsub")

	// send initial value
	value, _ := p.Repository().GetFeedLatestValue(ctx, feedID)
	if value != nil {
		ch <- value
		log.Info("sent latest value", zap.ByteString("value", value))
	}

	return ch, nil
}

func (p *pubSubValues) Unsubscribe(ctx context.Context, feedID string, ch <-chan []byte) {
	p.Lock()
	defer p.Unlock()

	log := telemetry.GetLogger(ctx)

	subscribers, ok := p.subscribers[feedID]
	if !ok {
		return
	}

	for rx, tx := range subscribers {
		if rx == ch {
			tx <- nil
			delete(subscribers, rx)
			close(tx)

			log.Info("unsubscribed from pubsub")
		}
	}
}

// AddFeed subscribes to the MQTT feed
func (p *pubSubValues) addFeed(ctx context.Context, feedID string) error {
	p.Lock()
	defer p.Unlock()

	log := telemetry.GetLogger(ctx)

	if _, ok := p.subscribers[feedID]; ok {
		return errutils.AlreadyExists(ctx, &errdetails.ResourceInfo{
			ResourceType: "Feed",
			ResourceName: feedID,
			Description:  fmt.Sprintf("Feed '%s' already exists", feedID),
		})
	}

	if token := p.mqtt.Subscribe(feedID, 0, func(c mqtt.Client, m mqtt.Message) {
		p.notify(ctx, feedID, m.Payload())
	}); token.Wait() && token.Error() != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"failed to subscribe to feed %s: %w", feedID, token.Error()))
	}

	p.subscribers[feedID] = make(map[<-chan []byte]chan<- []byte)

	log.Info("added feed to pubsub")

	return nil
}

func (p *pubSubValues) notify(ctx context.Context, feedID string, payload []byte) {
	log := telemetry.GetLogger(ctx)

	err := p.Repository().InsertFeedValue(ctx, feedID, payload)
	if err != nil {
		log.Warn("failed to insert feed value", zap.Error(err))
	}

	p.RLock()
	defer p.RUnlock()

	for _, tx := range p.subscribers[feedID] {
		tx <- payload

		log.Info("notified subscriber")
	}
}

func (p *pubSubValues) removeFeed(ctx context.Context, feed string) error {
	log := telemetry.GetLogger(ctx)

	if token := p.mqtt.Unsubscribe(feed); token.Wait() && token.Error() != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"failed to unsubscribe from feed %s: %w", feed, token.Error(),
		))
	}

	p.Lock()
	defer p.Unlock()
	for _, tx := range p.subscribers[feed] {
		tx <- nil
		close(tx)
		log.Info("unsubscribed subscriber")
	}
	delete(p.subscribers, feed)

	log.Info("removed feed from pubsub")

	return nil
}

func (p *pubSubValues) Publish(ctx context.Context, feedID string, value []byte) error {
	log := telemetry.GetLogger(ctx)

	if token := p.mqtt.Publish(feedID, 0, false, value); token.Wait() && token.Error() != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"failed to publish to feed %s: %w", feedID, token.Error(),
		))
	}

	// mqtt will notify this pubsub instance of the published message

	log.Info("published to feed")

	return nil
}
