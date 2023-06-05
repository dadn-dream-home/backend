package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
)

type StreamNotificationsHandler struct {
	state.State
}

func (h StreamNotificationsHandler) StreamNotifications(req *pb.StreamNotificationsRequest, stream pb.BackendService_StreamNotificationsServer) error {
	ctx := stream.Context()
	unsubFn, errCh := h.DatabaseHooker().Subscribe(func(t topic.Topic, rowid int64) error {
		switch t {
		case topic.Insert("notifications"):
			notification, err := h.Repository().GetNotificationByRowID(ctx, rowid)
			if err != nil {
				return err
			}

			return stream.Send(&pb.StreamNotificationsResponse{
				Notification: notification,
			})
		}

		return nil
	}, topic.Insert("notifications"))

	select {
	case <-stream.Context().Done():
		unsubFn()
		return nil
	case err := <-errCh:
		unsubFn()
		return err
	}
}
