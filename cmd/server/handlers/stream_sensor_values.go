package handlers

import (
	"strconv"

	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"

	"github.com/dadn-dream-home/x/server/state"
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

	ch, err := s.PubSubValues().Subscribe(ctx, req.Feed.Id)
	if err != nil {
		return err
	}

	log.Info("started streaming")

	for {
		select {
		case <-stream.Context().Done():
			log.Debug("cancelled streaming by client")
			err := s.PubSubValues().Unsubscribe(ctx, req.Feed.Id, ch)
			if err != nil {
				return err
			}

		case msg := <-ch:
			if msg == nil {
				log.Info("ended streaming")
				return nil
			}

			var value float64
			if value, err = strconv.ParseFloat(string(msg), 32); err != nil {
				log.Warn("error parsing float", zap.Error(err), zap.String("value", string(msg)))
				continue
			}

			if err = stream.Send(&pb.StreamSensorValuesResponse{
				Value: float32(value),
			}); err != nil {
				log.Warn("error sending message", zap.Error(err))
			}
		}
	}
}
