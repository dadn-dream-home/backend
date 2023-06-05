package handlers

import (
	"fmt"
	"strconv"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/errutils"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type StreamNotificationsHandler struct {
	state.State
}

func (h StreamNotificationsHandler) StreamNotifications(req *pb.StreamNotificationsRequest, stream pb.BackendService_StreamNotificationsServer) error {
	ctx := stream.Context()
	log := telemetry.GetLogger(ctx)

	unsubFn, errCh := h.DatabaseHooker().Subscribe(func(t topic.Topic, rowid int64) error {
		feedID, bytes, err := h.Repository().GetFeedValueByRowID(ctx, rowid)
		if err != nil {
			return err
		}

		feed, err := h.Repository().GetFeed(ctx, feedID)
		if err != nil {
			return err
		}

		feedConfig, err := h.Repository().GetFeedConfig(ctx, feed.Id)
		if err != nil {
			return err
		}

		log.Debug("feed config", zap.Any("config", feedConfig))

		if !feedConfig.GetSensorConfig().HasNotification {
			return nil
		}

		value, err := strconv.ParseFloat(string(bytes), 64)
		if err != nil {
			log.Warn("error parsing value", zap.Error(err), zap.String("value", string(bytes)))
			return nil
		}

		var notification *pb.Notification

		if value < float64(feedConfig.GetSensorConfig().LowerThreshold) {
			notification = &pb.Notification{
				Timestamp: timestamppb.Now(),
				Feed:      feed,
				Message:   fmt.Sprintf("%s (feed %s) is too low: %f", feed.Type, feed.Id, value),
			}
		} else if value > float64(feedConfig.GetSensorConfig().UpperThreshold) {
			notification = &pb.Notification{
				Timestamp: timestamppb.Now(),
				Feed:      feed,
				Message:   fmt.Sprintf("%s (feed %s) is too high: %f", feed.Type, feed.Id, value),
			}
		}

		if notification == nil {
			return nil
		}

		err = h.Repository().InsertNotification(ctx, notification)
		if err != nil {
			return errutils.Internal(ctx,
				fmt.Errorf("failed to insert notification: %w", err))
		}

		if err := stream.Send(&pb.StreamNotificationsResponse{
			Notification: notification,
		}); err != nil {
			return err
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
