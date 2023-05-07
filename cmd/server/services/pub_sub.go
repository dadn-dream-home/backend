package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type pubSub struct {
	sync.RWMutex
	state.State

	subscribers map[string][]subscriber
}

type subscriber struct {
	id   string
	ch   chan<- []byte
	done chan struct{}
}

var _ state.PubSub = (*pubSub)(nil)

func NewPubSub(ctx context.Context, state state.State) (state.PubSub, error) {
	ps := &pubSub{
		State:       state,
		subscribers: make(map[string][]subscriber),
	}

	// add initial feeds to pubsub
	rows, err := state.DB().Query("SELECT id from feeds")
	if err != nil {
		return nil, fmt.Errorf("error querying feeds: %w", err)
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		ps.AddFeed(ctx, id)
	}

	return ps, nil
}

func (p *pubSub) Subscribe(ctx context.Context, feed string, id string, ch chan<- []byte) (<-chan struct{}, error) {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if !p.HasFeed(ctx, feed) {
		return nil, fmt.Errorf("feed %s not found", feed)
	}

	p.Lock()
	subscriber := subscriber{
		id:   id,
		ch:   ch,
		done: make(chan struct{}, 1),
	}
	p.subscribers[feed] = append(p.subscribers[feed], subscriber)
	p.Unlock()

	log.Infof("subscribed to pubsub")

	return subscriber.done, nil
}

func (p *pubSub) HasFeed(ctx context.Context, feed string) bool {
	p.RLock()
	defer p.RUnlock()
	_, ok := p.subscribers[feed]
	return ok
}

func (p *pubSub) Unsubscribe(ctx context.Context, feed string, id string) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if !p.HasFeed(ctx, feed) {
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

			subscriber.done <- struct{}{}
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
func (p *pubSub) AddFeed(ctx context.Context, feed string) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if p.HasFeed(ctx, feed) {
		return fmt.Errorf("feed %s already added to pubsub", feed)
	}

	if token := p.MQTT().Subscribe(feed, 0, p.MessageHandler(ctx)); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to feed %s: %w", feed, token.Error())
	}

	p.Lock()
	p.subscribers[feed] = make([]subscriber, 0)
	p.Unlock()

	log.Infof("added feed to pubsub")

	return nil
}

func (p *pubSub) MessageHandler(ctx context.Context) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log := telemetry.GetLogger(ctx).WithField("feed", msg.Topic())

		log.Infof("received message")

		p.RLock()
		log.Tracef("acquired read lock")

		for _, subscriber := range p.subscribers[msg.Topic()] {

			log := log.WithField("subscriber_id", subscriber.id)
			log.Tracef("forwarding message to subscriber")

			subscriber.ch <- msg.Payload()

			log.Infof("forwarded message to subscriber")
		}
		p.RUnlock()

		log.WithField("subscribers", len(p.subscribers[msg.Topic()])).
			Infof("forwarded message to subscribers")
	}
}

func (p *pubSub) RemoveFeed(ctx context.Context, feed string) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if token := p.MQTT().Unsubscribe(feed); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from feed %s: %w", feed, token.Error())
	}

	p.Lock()
	for _, subscriber := range p.subscribers[feed] {
		subscriber.done <- struct{}{}
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
