package services

import (
	"context"
	"fmt"

	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type mqttSubscriber struct {
	state.State

	mqtt.Client
}

func NewMQTTSubscriber(state state.State, mqtt mqtt.Client) state.MQTTSubscriber {
	return &mqttSubscriber{
		State:  state,
		Client: mqtt,
	}
}

func (m *mqttSubscriber) Serve(ctx context.Context) error {
	// initial feeds
	feeds, err := m.Repository().ListFeeds(ctx)
	if err != nil {
		return err
	}

	if len(feeds) != 0 {
		filters := make(map[string]byte, len(feeds))
		for _, feed := range feeds {
			filters[feed.Id] = 0
		}

		if token := m.SubscribeMultiple(filters, m.handleMessage); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}

	unsubFn, errCh := m.DatabaseHooker().Subscribe(func(t topic.Topic, rowid int64) error {

		switch t {
		case topic.Insert("feeds"):
			feed, err := m.Repository().GetFeedByRowID(ctx, rowid)
			if err != nil {
				return err
			}
			if token := m.Subscribe(feed.Id, 0, m.handleMessage); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		case topic.Insert("deleted_feeds"):
			feed, err := m.Repository().GetDeletedFeedByRowID(ctx, rowid)
			if err != nil {
				return err
			}
			if token := m.Unsubscribe(feed.Id); token.Wait() && token.Error() != nil {
				return token.Error()
			}
		}

		return nil
	}, topic.Insert("feeds"), topic.Insert("deleted_feeds"))

	select {
	case <-ctx.Done():
		unsubFn()
		return nil
	case err := <-errCh:
		unsubFn()
		return err
	}
}

func (m *mqttSubscriber) handleMessage(client mqtt.Client, msg mqtt.Message) {
	ctx := telemetry.InitLogger(context.Background())
	log := telemetry.GetLogger(ctx)

	err := m.Repository().InsertFeedValue(ctx, msg.Topic(), msg.Payload())
	if err != nil {
		log.Warn("failed to insert feed value", zap.Error(err))
	}
}

func (m *mqttSubscriber) Publish(ctx context.Context, feedID string, value []byte) error {
	log := telemetry.GetLogger(ctx)

	if token := m.Client.Publish(feedID, 0, false, value); token.Wait() && token.Error() != nil {
		return errutils.Internal(ctx, fmt.Errorf(
			"failed to publish to feed %s: %w", feedID, token.Error(),
		))
	}

	log.Info("published to feed")

	return nil
}
