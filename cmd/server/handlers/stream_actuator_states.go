package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamActuatorStatesHandler struct {
	state.State
}

func (s StreamActuatorStatesHandler) StreamActuatorStates(
	req *pb.StreamActuatorStatesRequest,
	stream pb.BackendService_StreamActuatorStatesServer,
) error {
	ctx := stream.Context()
	log := telemetry.GetLogger(ctx)

	// validate request

	if req.Feed == nil {
		log.Error("feed is nil")
		return status.Error(codes.InvalidArgument, "feed cannot be nil")
	}

	log = log.With(zap.String("feed.id", req.Feed.Id))
	ctx = telemetry.ContextWithLogger(ctx, log)

	feed, err := s.Repository().GetFeed(ctx, req.Feed.Id)
	if err != nil {
		return err
	}

	switch feed.Type {
	case pb.FeedType_LIGHT:
		// ok
	default:
		log.Error("feed type is not actuator", zap.String("feed.type", feed.Type.String()))
		return status.Errorf(codes.InvalidArgument, "feed type %s is not actuator", feed.Type)
	}

	// send latest value

	bytes, err := s.Repository().GetFeedLatestValue(ctx, req.Feed.Id)
	if err != nil {
		return err
	}

	if bytes != nil {
		if err := stream.Send(&pb.StreamActuatorStatesResponse{
			State: string(bytes) != "0",
		}); err != nil {
			return err
		}
	}

	// subscribe to feed

	unsubFn, errCh := s.DatabaseHooker().Subscribe(func(t topic.Topic, rowid int64) error {
		feedID, bytes, err := s.Repository().GetFeedValueByRowID(ctx, rowid)
		if err != nil {
			return err
		}

		if feedID != req.Feed.Id {
			return nil
		}

		if bytes != nil {
			if err := stream.Send(&pb.StreamActuatorStatesResponse{
				State: string(bytes) != "0",
			}); err != nil {
				return err
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
