package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
)

type StreamNotificationsHandler struct {
	state.State
}

func (h StreamNotificationsHandler) StreamNotifications(req *pb.StreamNotificationsRequest, stream pb.BackendService_StreamNotificationsServer) error {
	ctx := stream.Context()
	log := telemetry.GetLogger(ctx)

	ch, err := h.Notifier().Subscribe(stream.Context())
	if err != nil {
		return err
	}

	log.Info("started streaming")

	for {
		select {
		case <-ctx.Done():
			log.Debug("cancelled streaming by client")

			err := h.Notifier().Unsubscribe(ctx, ch)
			if err != nil {
				return err
			}

		case notification := <-ch:
			if notification == nil {
				log.Info("ended streaming")
				return nil
			}

			log = log.With(zap.String("feed.id", notification.Feed.Id))
			log.Debug("got notification", zap.String("notification.message", notification.Message))

			if err = stream.Send(&pb.StreamNotificationsResponse{
				Notification: notification,
			}); err != nil {
				log.Warn("failed to send notification", zap.Error(err))
			}

			log.Info("sent notification to client")
		}
	}
}
