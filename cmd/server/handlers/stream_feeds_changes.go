package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"

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
	rid := telemetry.GetRequestId(ctx)
	log := telemetry.GetLogger(ctx)

	log.Infof("streaming feed changes")

	ch := make(chan *pb.FeedsChange)
	done, err := h.PubSubFeeds().Subscribe(ctx, rid, ch)
	if err != nil {
		log.WithError(err).Errorf("failed to subscribe to pub sub feeds")
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			log.Tracef("<-stream.Context().Done()")
			err := h.PubSubFeeds().Unsubscribe(ctx, rid)
			if err != nil {
				log.WithError(err).Errorf("error unsubscribing from feeds")
			}
			return nil
		case <-done:
			log.Tracef("<-done")
			log.Infof("streaming feed changes done")
			return nil
		case change := <-ch:
			if err := stream.Send(&pb.StreamFeedsChangesResponse{
				Change: change,
			}); err != nil {
				log.WithError(err).Errorf("failed to send feed change")
				return err
			}
		}
	}
}
