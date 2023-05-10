package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type notifier struct {
	sync.RWMutex
	state.State

	listeners map[string]<-chan []byte

	// map from rx to tx
	subscribers map[<-chan *pb.Notification]chan<- *pb.Notification
}

func NewNotifier(ctx context.Context, state state.State) state.Notifier {
	n := &notifier{
		State:       state,
		subscribers: make(map[<-chan *pb.Notification]chan<- *pb.Notification),
	}

	return n
}

func (n *notifier) Subscribe(ctx context.Context) (<-chan *pb.Notification, error) {
	n.Lock()
	defer n.Unlock()

	ch := make(chan *pb.Notification)
	n.subscribers[ch] = ch
	return ch, nil
}

func (n *notifier) Unsubscribe(ctx context.Context, ch <-chan *pb.Notification) {
	n.Lock()
	defer n.Unlock()

	n.subscribers[ch] <- nil
	delete(n.subscribers, ch)
}

func (n *notifier) Serve(ctx context.Context) (err error) {
	log := telemetry.GetLogger(ctx)
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		if err == nil {
			err = recover().(error)
		}
	}()

	ch := n.State.PubSubFeeds().Subscribe(ctx)

	for {
		select {
		case <-ctx.Done():
			n.PubSubFeeds().Unsubscribe(ctx, ch)

		case changes := <-ch:
			if changes == nil {
				return nil
			}

			for _, feed := range changes.Addeds {
				log = log.With(zap.String("feed.id", feed.Id))

				// ignore feeds that are not temperature or humidity
				if feed.Type != pb.FeedType_TEMPERATURE && feed.Type != pb.FeedType_HUMIDITY {
					continue
				}

				rx, err := n.subscribeToFeed(ctx, feed.Id)
				if err != nil {
					return errutils.Internal(
						fmt.Errorf("failed to subscribe to feed: %w", err))
				}

				go func(feed *pb.Feed) {
					if err := n.listenToFeed(ctx, feed, rx); err != nil {
						panic(err)
					}
				}(feed)
			}

			for _, feedID := range changes.RemovedIDs {
				log = log.With(zap.String("feed.id", feedID))

				if _, ok := n.listeners[feedID]; !ok {
					continue
				}
				n.unsubscribeFromFeed(ctx, feedID)
			}
		}
	}
}

func (n *notifier) subscribeToFeed(ctx context.Context, feedID string) (<-chan []byte, error) {
	n.Lock()
	defer n.Unlock()

	rx, err := n.PubSubValues().Subscribe(ctx, feedID)
	if err != nil {
		return nil, err
	}

	n.listeners[feedID] = rx

	return rx, nil
}

func (n *notifier) listenToFeed(ctx context.Context, feed *pb.Feed, rx <-chan []byte) error {
	log := telemetry.GetLogger(ctx)

	for {
		select {
		case <-ctx.Done():
			n.PubSubValues().Unsubscribe(ctx, feed.Id, rx)

		case payload := <-rx:
			if payload == nil {
				log.Info("stopped listener by feed removed")
				return nil
			}

			lower := 10
			upper := 30

			value, err := strconv.ParseFloat(string(payload), 64)
			if err != nil {
				log.Warn("error parsing value", zap.Error(err), zap.String("value", string(payload)))
				continue
			}

			var notification *pb.Notification

			if value < float64(lower) {
				notification = &pb.Notification{
					Timestamp: timestamppb.Now(),
					Feed:      feed,
					Message:   fmt.Sprintf("%s (feed %s) is too low: %f", feed.Type, feed.Id, value),
				}
			} else if value > float64(upper) {
				notification = &pb.Notification{
					Timestamp: timestamppb.Now(),
					Feed:      feed,
					Message:   fmt.Sprintf("%s (feed %s) is too high: %f", feed.Type, feed.Id, value),
				}
			}

			if notification == nil {
				continue
			}

			err = n.Repository().InsertNotification(ctx, notification)
			if err != nil {
				return errutils.Internal(
					fmt.Errorf("failed to insert notification: %w", err))
			}

			for _, ch := range n.subscribers {
				ch <- notification
			}
		}
	}
}

func (n *notifier) unsubscribeFromFeed(ctx context.Context, feedID string) {
	n.Lock()
	defer n.Unlock()

	n.PubSubValues().Unsubscribe(ctx, feedID, n.listeners[feedID])
	delete(n.listeners, feedID)
}
