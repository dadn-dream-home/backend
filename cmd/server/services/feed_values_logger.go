package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type logger struct {
	state.State
}

func NewFeedValuesLogger(state state.State) state.FeedValuesLogger {
	return logger{state}
}

func (l logger) Serve(ctx context.Context) error {
	log := telemetry.GetLogger(ctx)

	unsubFn, errCh := l.DatabaseHooker().Subscribe(func(t topic.Topic, rowid int64) error {
		feedID, bytes, err := l.Repository().GetFeedValueByRowID(ctx, rowid)
		if err != nil {
			return err
		}

		feed, err := l.Repository().GetFeed(ctx, feedID)
		if err != nil {
			return err
		}

		feedConfig, err := l.Repository().GetFeedConfig(ctx, feed.Id)
		if err != nil {
			return err
		}

		log.Debug("feed config", zap.Any("config", feedConfig))

		switch feed.Type {
		case pb.FeedType_TEMPERATURE, pb.FeedType_HUMIDITY:
			if feedConfig.GetSensorConfig().HasNotification {
				return l.logFeedValues(ctx, bytes, feed, feedConfig.GetSensorConfig())
			}
		}

		return nil
	}, topic.Insert("feed_values"))

	select {
	case <-ctx.Done():
		unsubFn()
		return ctx.Err()
	case err := <-errCh:
		unsubFn()
		return err
	}
}

func (l logger) logFeedValues(ctx context.Context, bytes []byte, feed *pb.Feed, config *pb.SensorConfig) error {
	log := telemetry.GetLogger(ctx)

	value, err := strconv.ParseFloat(string(bytes), 64)
	if err != nil {
		log.Warn("error parsing value", zap.Error(err), zap.String("value", string(bytes)))
		return nil
	}

	var notification *pb.Notification

	if value < float64(config.LowerThreshold) {
		notification = &pb.Notification{
			Timestamp: timestamppb.Now(),
			Feed:      feed,
			Message:   fmt.Sprintf("%s (feed %s) is too low: %f", feed.Type, feed.Id, value),
		}
	} else if value > float64(config.UpperThreshold) {
		notification = &pb.Notification{
			Timestamp: timestamppb.Now(),
			Feed:      feed,
			Message:   fmt.Sprintf("%s (feed %s) is too high: %f", feed.Type, feed.Id, value),
		}
	}

	if notification == nil {
		return nil
	}

	return l.Repository().InsertNotification(ctx, notification)
}
