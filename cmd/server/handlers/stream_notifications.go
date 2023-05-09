package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type StreamNotificationsHandler struct {
	state.State
}

func (h StreamNotificationsHandler) StreamNotifications(req *pb.StreamNotificationsRequest, stream pb.BackendService_StreamNotificationsServer) error {
	ctx := stream.Context()
	log := telemetry.GetLogger(ctx)
	rid := telemetry.GetRequestId(ctx)

	ch, err := h.Notifier().Subscribe(stream.Context(), rid)
	if err != nil {
		return err
	}

	log.Debugf("begin streaming")

	for {
		select {
		case <-stream.Context().Done():
			log.Tracef("<-stream.Context().Done()")
			err := h.Notifier().Unsubscribe(ctx, rid)
			if err != nil {
				log.WithError(err).Warnf("error unsubscribing from feed")
			}
			
		case msg := <-ch:
			if msg == nil {
				return nil
			}

			err = stream.Send(&pb.StreamNotificationsResponse{
				Notification: msg,
			})

			if err != nil {
				return err
			}
		}
	}
}
