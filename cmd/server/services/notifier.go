package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type notifier struct {
	state.State

	subscribers map[string]chan<- *pb.Notification
}

func NewNotifier(ctx context.Context, state state.State) state.Notifier {
	n := &notifier{state, make(map[string]chan<- *pb.Notification)}

	go n.subscribeToValues(ctx)

	return n
}

func (n *notifier) Subscribe(ctx context.Context, sid string) (<-chan *pb.Notification, error) {
	ch := make(chan *pb.Notification, 1)
	n.subscribers[sid] = ch
	return ch, nil
}

func (n *notifier) Unsubscribe(ctx context.Context, sid string) error {
	n.subscribers[sid] <- nil
	delete(n.subscribers, sid)
	return nil
}

func (n *notifier) subscribeToValues(ctx context.Context) {
	ch := make(chan state.PubSubValuesAllData, 1)
	n.State.PubSubValues().SubscribeAll(ctx, "notifier", ch)
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-ch:
			n.notify(ctx, data)
		}
	}
}

func (n *notifier) notify(ctx context.Context, data state.PubSubValuesAllData) {
	log := telemetry.GetLogger(ctx).WithField("feed_id", data.Feed)

	feed, err := n.Repository().GetFeed(ctx, data.Feed)
	if err != nil {
		log.WithError(err).Errorf("error getting feed")
		return
	}

	if feed.Type != pb.FeedType_TEMPERATURE && feed.Type != pb.FeedType_HUMIDITY {
		return
	}

	lower := 10
	upper := 30

	value, err := strconv.ParseFloat(string(data.Value), 64)
	if err != nil {
		log.WithError(err).Warnf("error parsing value")
		return
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
		return
	}

	err = n.Repository().InsertNotification(ctx, notification)
	if err != nil {
		log.WithError(err).Errorf("error inserting notification")
		return
	}

	for _, ch := range n.subscribers {
		select {
		case ch <- notification:
		default:
		}
	}
}
