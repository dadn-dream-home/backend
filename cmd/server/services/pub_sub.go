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

	subscribers map[string][]chan<- []byte
}

var _ state.PubSub = (*pubSub)(nil)

func NewPubSub(ctx context.Context, state state.State) (state.PubSub, error) {
	ps := &pubSub{
		State:       state,
		subscribers: make(map[string][]chan<- []byte),
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

func (p *pubSub) Subscribe(ctx context.Context, feed string, ch chan<- []byte) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if !p.HasFeed(ctx, feed) {
		return fmt.Errorf("feed %s not found", feed)
	}

	p.Lock()
	p.subscribers[feed] = append(p.subscribers[feed], ch)
	p.Unlock()

	log.Infof("subscribed to pubsub")

	return nil
}

func (p *pubSub) HasFeed(ctx context.Context, feed string) bool {
	p.RLock()
	defer p.RUnlock()
	_, ok := p.subscribers[feed]
	return ok
}

func (p *pubSub) Unsubscribe(ctx context.Context, feed string, ch chan<- []byte) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if !p.HasFeed(ctx, feed) {
		return fmt.Errorf("feed %s not found", feed)
	}

	p.RLock()
	for i, subscriber := range p.subscribers[feed] {
		if subscriber == ch {
			p.RUnlock()
			p.Lock()
			p.subscribers[feed] = append(p.subscribers[feed][:i], p.subscribers[feed][i+1:]...)
			close(subscriber)
			p.Unlock()

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
	p.subscribers[feed] = make([]chan<- []byte, 0)
	p.Unlock()

	log.Infof("added feed to pubsub")

	return nil
}

func (p *pubSub) MessageHandler(ctx context.Context) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log := telemetry.GetLogger(ctx).WithField("feed", msg.Topic())

		log.Infof("received message")

		p.RLock()
		for _, subscriber := range p.subscribers[msg.Topic()] {
			subscriber <- msg.Payload()
		}
		p.RUnlock()

		log.WithField("subscribers", len(p.subscribers[msg.Topic()])).Infof("forwarded message to subscribers")
	}
}

func (p *pubSub) RemoveFeed(ctx context.Context, feed string) error {
	log := telemetry.GetLogger(ctx).WithField("feed", feed)

	if token := p.MQTT().Unsubscribe(feed); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from feed %s: %w", feed, token.Error())
	}

	p.Lock()
	for _, subscriber := range p.subscribers[feed] {
		close(subscriber)
	}
	delete(p.subscribers, feed)
	p.Unlock()

	if len := len(p.subscribers[feed]); len > 0 {
		log.WithField("subscribers", len).Warnf("closed subscribers")
	}

	log.Infof("removed feed from pubsub")

	return nil
}
