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

func (s StreamSensorValuesHandler) StreamSensorValues(req *pb.StreamFeedValuesRequest, stream pb.BackendService_StreamSensorValuesServer) error {
	ctx := stream.Context()
	rid := telemetry.GetRequestId(ctx)
	log := telemetry.GetLogger(ctx).WithField("feed_id", req.Id)

	ch := make(chan []byte)

	done, err := s.PubSub().Subscribe(ctx, req.Id, rid, ch)
	if err != nil {
		return err
	}

	log.Debugf("begin streaming")

	for {
		select {
		case <-stream.Context().Done():
			log.Tracef("<-stream.Context().Done()")
			err := s.PubSub().Unsubscribe(ctx, req.Id, rid)
			if err != nil {
				log.WithError(err).Warnf("error unsubscribing from feed")
			}
			// wait until done is closed

		case <-done:
			log.Tracef("<-done")
			return status.Errorf(codes.Canceled, "stream canceled")

		case msg := <-ch:
			value, err := strconv.ParseFloat(string(msg), 32)
			if err != nil {
				log.WithError(err).Warnf("error parsing float %s: %w", string(msg))
				continue
			}
			err = stream.Send(&pb.StreamFeedValuesResponse{
				Value: float32(value),
			})
			if err != nil {
				return err
			}
		}
	}
}
