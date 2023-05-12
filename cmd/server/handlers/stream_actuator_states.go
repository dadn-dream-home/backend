package handlers

import (
	pb "github.com/dadn-dream-home/x/protobuf"
	"go.uber.org/zap"

	"github.com/dadn-dream-home/x/server/state"
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

	ch, err := s.PubSubValues().Subscribe(ctx, req.Feed.Id)
	if err != nil {
		return err
	}

	log.Info("begin streaming")

	for {
		select {
		case <-ctx.Done():
			log.Debug("cancelled streaming by client")
			s.PubSubValues().Unsubscribe(ctx, req.Feed.Id, ch)
			for msg := <-ch; msg != nil; msg = <-ch {
				// skip remaining messages
			}
			log.Info("end streaming")
			return nil

		case msg := <-ch:
			if msg == nil {
				log.Info("end streaming")
				return nil
			}

			state := string(msg) != "0"

			if err = stream.Send(&pb.StreamActuatorStatesResponse{
				State: state,
			}); err != nil {
				log.Warn("error sending response", zap.Error(err))
			}

			log.Info("sent state to client", zap.Bool("state", state))
		}
	}
}
