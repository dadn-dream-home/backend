package handlers

import (
	"fmt"

	pb "github.com/dadn-dream-home/x/protobuf"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
)

type StreamActivitiesHandler struct {
	state.State
}

func (h StreamActivitiesHandler) StreamActivities(
	req *pb.StreamActivitiesRequest,
	stream pb.BackendService_StreamActivitiesServer,
) error {
	ctx := stream.Context()
	log := telemetry.GetLogger(ctx)

	// send initial activities

	activities, err := h.Repository().ListActivities(ctx)
	if err != nil {
		return err
	}

	for _, activity := range activities {
		if err := stream.Send(&pb.StreamActivitiesResponse{
			Activity: activity,
		}); err != nil {
			return err
		}
	}

	log.Debug("sent initial activities", zap.Int("activities.len", len(activities)))

	log.Info("streaming activities changes")

	unsubFn, errCh := h.DatabaseHooker().Subscribe(func(t topic.Topic, rowid int64) error {
		feedID, bytes, err := h.Repository().GetFeedValueByRowID(ctx, rowid)
		if err != nil {
			return err
		}

		feed, err := h.Repository().GetFeed(ctx, feedID)
		if err != nil {
			return err
		}

		switch feed.Type {
		case pb.FeedType_LIGHT:
			if err := stream.Send(&pb.StreamActivitiesResponse{
				Activity: &pb.Activity{
					Feed:  feed,
					State: string(bytes) != "0",
				}}); err != nil {
				return err
			}

			return nil
		case pb.FeedType_TEMPERATURE, pb.FeedType_HUMIDITY:
			return nil
		default:
			return fmt.Errorf("unknown feed type: %v", feed.Type)
		}
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
