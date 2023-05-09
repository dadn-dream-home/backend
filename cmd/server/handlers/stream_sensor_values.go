package handlers

import (
	"strconv"

	pb "github.com/dadn-dream-home/x/protobuf"

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
	rid := telemetry.GetRequestId(ctx)
	log := telemetry.GetLogger(ctx).WithField("feed_id", req.Id)

	feed, err := s.Repository().GetFeed(ctx, req.Id)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	switch feed.Type {
	case pb.FeedType_TEMPERATURE, pb.FeedType_HUMIDITY:
		// ok
	default:
		return status.Errorf(codes.InvalidArgument, "feed type %s is not sensor", feed.Type)
	}

	ch, err := s.PubSubValues().Subscribe(ctx, rid, req.Id)
	if err != nil {
		return err
	}

	log.Debugf("begin streaming")

	for {
		select {
		case <-stream.Context().Done():
			log.Tracef("<-stream.Context().Done()")
			err := s.PubSubValues().Unsubscribe(ctx, rid, req.Id)
			if err != nil {
				log.WithError(err).Warnf("error unsubscribing from feed")
			}

		case msg := <-ch:
			if msg == nil {
				log.Infof("streaming finished")
				return nil
			}

			value, err := strconv.ParseFloat(string(msg), 32)
			if err != nil {
				log.WithError(err).Warnf("error parsing float %s: %w", string(msg))
				continue
			}
			err = stream.Send(&pb.StreamSensorValuesResponse{
				Value: float32(value),
			})
			if err != nil {
				return err
			}
		}
	}
}
