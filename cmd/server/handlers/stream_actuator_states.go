package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamActuatorStatesHandler struct {
	state.State
}

func (s StreamActuatorStatesHandler) StreamActuatorStates(req *pb.StreamActuatorStatesRequest, stream pb.BackendService_StreamActuatorStatesServer) error {
	ctx := stream.Context()
	rid := telemetry.GetRequestId(ctx)
	log := telemetry.GetLogger(ctx).WithField("feed_id", req.Id)

	feed, err := s.Repository().GetFeed(ctx, req.Id)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	switch feed.Type {
	case pb.FeedType_LIGHT:
		// ok
	default:
		return status.Errorf(codes.InvalidArgument, "feed type %s is not actuator", feed.Type)
	}

	ch := make(chan []byte)
	done, err := s.PubSubValues().Subscribe(ctx, req.Id, rid, ch)
	if err != nil {
		return err
	}

	log.Debugf("begin streaming")

	for {
		select {
		case <-stream.Context().Done():
			log.Tracef("<-stream.Context().Done()")
			err := s.PubSubValues().Unsubscribe(ctx, req.Id, rid)
			if err != nil {
				log.WithError(err).Warnf("error unsubscribing from feed")
			}
			return status.Errorf(codes.Canceled, "stream canceled")
			// wait until done is closed

		case <-done:
			log.Tracef("<-done")
			return status.Errorf(codes.Canceled, "stream canceled")

		case msg := <-ch:
			err = stream.Send(&pb.StreamActuatorStatesResponse{
				State: string(msg) != "0",
			})
			if err != nil {
				return err
			}
		}
	}
}
