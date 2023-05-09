package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type pubSubValues struct {
	sync.RWMutex

	state.State
	mqtt mqtt.Client

	subscribers    map[string][]subscriber
	subscribersAll []subscriberAll
}

type subscriber struct {
	id string
	ch chan<- []byte
}

type subscriberAll struct {
	id string
	ch chan<- map[string][]byte
}

func NewPubSubValues(ctx context.Context, state state.State, mqtt mqtt.Client) state.PubSubValues {
	ps := &pubSubValues{
		State:       state,
		mqtt:        mqtt,
		subscribers: make(map[string][]subscriber),
	}

	return ps
}

func (p *pubSubValues) ListenToFeedsChanges(ctx context.Context) {
	log := telemetry.GetLogger(ctx)

	ch, err := p.PubSubFeeds().Subscribe(ctx, "pub-sub-values")
	if err != nil {
		log.Fatalf("failed to subscribe to pubsub feeds: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Infof("context cancelled")
			p.PubSubFeeds().Unsubscribe(ctx, "pub-sub-values")

		case change := <-ch:
			if change == nil {
				log.Infof("stream ended")
				return
			}

			for _, feed := range change.Added {
				if err := p.AddFeed(ctx, feed.Id); err != nil {
					log.Fatalf("failed to add feed %s: %v", feed, err)
				}
			}
			for _, feed := range change.Removed {
				if err := p.RemoveFeed(ctx, feed); err != nil {
					log.Fatalf("failed to remove feed %s: %v", feed, err)
				}
			}
		}
	}
}

func (p *pubSubValues) Subscribe(ctx context.Context, id string, feed string) (<-chan []byte, error) {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if !p.hasFeed(ctx, feed) {
		return nil, fmt.Errorf("feed %s not found", feed)
	}

	p.Lock()
	ch := make(chan []byte, 1)
	subscriber := subscriber{
		id: id,
		ch: ch,
	}
	p.subscribers[feed] = append(p.subscribers[feed], subscriber)
	p.Unlock()

	log.Infof("subscribed to pubsub")

	// send initial value
	go func() {
		value, _ := p.Repository().GetFeedLatestValue(ctx, feed)
		ch <- value
		log.Infof("sent latest value")
	}()

	return ch, nil
}

func (p *pubSubValues) hasFeed(ctx context.Context, feed string) bool {
	p.RLock()
	defer p.RUnlock()
	_, ok := p.subscribers[feed]
	return ok
}

func (p *pubSubValues) Unsubscribe(ctx context.Context, id string, feed string) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if !p.hasFeed(ctx, feed) {
		return fmt.Errorf("feed %s not found", feed)
	}

	p.RLock()
	log.Tracef("acquired read lock")

	for i, subscriber := range p.subscribers[feed] {
		if subscriber.id == id {
			p.RUnlock()
			log.Tracef("released read lock")

			p.Lock()
			log.Tracef("upgraded to write lock")

			subscriber.ch <- nil
			p.subscribers[feed] = append(p.subscribers[feed][:i], p.subscribers[feed][i+1:]...)
			close(subscriber.ch)

			p.Unlock()
			log.Tracef("released write lock")

			log.Infof("unsubscribed from pubsub")

			return nil
		}
	}
	p.RUnlock()

	return fmt.Errorf("subscriber not found")
}

// AddFeed subscribes to the MQTT feed
func (p *pubSubValues) AddFeed(ctx context.Context, feed string) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if p.hasFeed(ctx, feed) {
		return fmt.Errorf("feed %s already added to pubsub", feed)
	}

	if token := p.mqtt.Subscribe(feed, 0, func(c mqtt.Client, m mqtt.Message) {
		p.Notify(ctx, feed, m.Payload())
	}); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to feed %s: %w", feed, token.Error())
	}

	p.Lock()
	p.subscribers[feed] = make([]subscriber, 0)
	p.Unlock()

	log.Infof("added feed to pubsub")

	return nil
}

func (p *pubSubValues) Notify(ctx context.Context, feed string, payload []byte) {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	log.Infof("received message")

	err := p.Repository().InsertFeedValue(ctx, feed, payload)
	if err != nil {
		log.Errorf("failed to insert feed value: %w", err)
	}

	p.RLock()
	log.Tracef("acquired read lock")

	for _, subscriber := range p.subscribers[feed] {

		log := log.WithField("subscriber_id", subscriber.id)
		log.Tracef("forwarding message to subscriber")

		subscriber.ch <- payload

		log.Infof("notified subscriber")
	}

	for _, subscriber := range p.subscribersAll {
		log := log.WithField("subscriber_id", subscriber.id)
		log.Tracef("forwarding message to all-subscriber")

		data := make(map[string][]byte, 1)
		data[feed] = payload
		subscriber.ch <- data

		log.Infof("notified all-subscriber")
	}

	p.RUnlock()

	log.WithField("subscribers", len(p.subscribers[feed])).
		Infof("notified subscribers")
}

func (p *pubSubValues) RemoveFeed(ctx context.Context, feed string) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if token := p.mqtt.Unsubscribe(feed); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from feed %s: %w", feed, token.Error())
	}

	p.Lock()
	for _, subscriber := range p.subscribers[feed] {
		subscriber.ch <- nil
		close(subscriber.ch)
	}
	delete(p.subscribers, feed)
	p.Unlock()

	if len := len(p.subscribers[feed]); len > 0 {
		log.WithField("subscribers", len).Warnf("closed subscribers")
	}

	log.Infof("removed feed from pubsub")

	return nil
}

func (p *pubSubValues) Publish(ctx context.Context, id string, feed string, value []byte) error {
	log := telemetry.GetLogger(ctx).WithField("publisher_id", id).WithField("feed", feed)

	log.Debugf("publishing to feed")

	if token := p.mqtt.Publish(feed, 0, false, value); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish to feed %s: %w", feed, token.Error())
	}

	// mqtt will notify this pubsub instance of the published message

	log.Infof("published to feed")

	return nil
}

func (p *pubSubValues) SubscribeAll(ctx context.Context, id string) (<-chan map[string][]byte, error) {
	log := telemetry.GetLogger(ctx)

	p.Lock()
	ch := make(chan map[string][]byte, 1)
	subscriber := subscriberAll{
		id: id,
		ch: ch,
	}
	p.subscribersAll = append(p.subscribersAll, subscriber)
	p.Unlock()

	log.Infof("subscribed to all pubsub")

	return ch, nil
}
