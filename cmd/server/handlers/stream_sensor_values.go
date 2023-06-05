package handlers

import (
	"strconv"

	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/state/topic"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamSensorValuesHandler struct {
	state.State
}

func (s StreamSensorValuesHandler) StreamSensorValues(req *pb.StreamSensorValuesRequest, stream pb.BackendService_StreamSensorValuesServer) error {
	ctx := stream.Context()
	log := telemetry.GetLogger(ctx)

	if req.Feed == nil {
		log.Error("feed is nil")
		return status.Error(codes.InvalidArgument, "feed cannot be nil")
	}

	log = log.With(zap.String("feed.id", req.Feed.Id))

	feed, err := s.Repository().GetFeed(ctx, req.Feed.Id)
	if err != nil {
		return err
	}

	switch feed.Type {
	case pb.FeedType_TEMPERATURE, pb.FeedType_HUMIDITY:
		// ok
	default:
		log.Error("feed type is not sensor", zap.String("feed.type", feed.Type.String()))
		return status.Errorf(codes.InvalidArgument, "feed type %s is not sensor", feed.Type)
	}

	// send latest value
	bytes, err := s.Repository().GetFeedLatestValue(ctx, req.Feed.Id)

	if err != nil {
		return err
	}

	if bytes != nil {
		value, err := strconv.ParseFloat(string(bytes), 32)
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.StreamSensorValuesResponse{
			Value: float32(value),
		}); err != nil {
			return err
		}
	}

	// subscribe to new values

	unsubFn, errCh := s.DatabaseHooker().Subscribe(func(_ topic.Topic, rowid int64) error {
		feedID, bytes, err := s.Repository().GetFeedValueByRowID(ctx, rowid)
		if err != nil {
			return err
		}

		if feedID != req.Feed.Id {
			return nil
		}

		value, err := strconv.ParseFloat(string(bytes), 32)
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.StreamSensorValuesResponse{
			Value: float32(value),
		}); err != nil {
			return err
		}

		return nil
	}, topic.Insert("feed_values"))

	for {
		select {
		case <-ctx.Done():
			log.Debug("cancelled streaming by client")
			unsubFn()
			return nil
		case err := <-errCh:
			log.Error("failed to subscribe to new values", zap.Error(err))
		}
	}
}
