package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
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

	log.Info("streaming feeds changes")

	// send initial batch of feeds

	feeds, err := h.Repository().ListFeeds(ctx)
	if err != nil {
		return err
	}

	if err := stream.Send(&pb.StreamFeedsChangesResponse{
		Change: &pb.FeedsChange{
			Addeds: feeds,
		},
	}); err != nil {
		return err
	}

	unsubFn, errCh := h.DatabaseListener().Subscribe(func(t topic.Topic, rowid int64) error {
		switch t {
		case topic.Insert("feeds"):
			feed, err := h.Repository().GetFeedByRowID(ctx, rowid)
			if err != nil {
				return err
			}

			return stream.Send(&pb.StreamFeedsChangesResponse{
				Change: &pb.FeedsChange{
					Addeds: []*pb.Feed{feed},
				},
			})
		case topic.Delete("feeds"):
			feed, err := h.Repository().GetFeedByRowID(ctx, rowid)
			if err != nil {
				return err
			}

			return stream.Send(&pb.StreamFeedsChangesResponse{
				Change: &pb.FeedsChange{
					RemovedIDs: []string{feed.Id},
				},
			})
		default:
			log.Error("unknown topic operation", zap.Int("topic.op", t.Op))
			return nil
		}
	}, topic.Insert("feeds"), topic.Delete("feeds"))

	select {
	case <-ctx.Done():
		unsubFn()
		return nil
	case err := <-errCh:
		log.Error("error while streaming feeds changes", zap.Error(err))
		unsubFn()
		return err
	}
}
