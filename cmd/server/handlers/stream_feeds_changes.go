package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
)

type StreamFeedsChangesHandler struct {
	state.State
}

func (h StreamFeedsChangesHandler) StreamFeedsChanges(
	req *pb.StreamFeedsChangesRequest,
	stream pb.BackendService_StreamFeedsChangesServer,
) error {
	ctx := stream.Context()
	log := telemetry.GetLogger(ctx)

	ch, err := h.PubSubFeeds().Subscribe(ctx)
	if err != nil {
		return err
	}

	log.Info("started streaming")

	for {
		select {
		case <-stream.Context().Done():
			log.Debug("cancelled streaming by client")
			h.PubSubFeeds().Unsubscribe(ctx, ch)
			for change := <-ch; change != nil; change = <-ch {
				// skip remaining messages
			}
			log.Info("ended streaming")
			return nil

		case change := <-ch:
			if change == nil {
				log.Debug("ended streaming")
				return nil
			}

			if err := stream.Send(&pb.StreamFeedsChangesResponse{
				Change: change,
			}); err != nil {
				log.Warn("failed to send feed change", zap.Error(err))
				continue
			}

			log.Info("sent changes to client")
		}
	}
}
