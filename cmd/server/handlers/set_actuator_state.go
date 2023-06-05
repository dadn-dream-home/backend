package handlers

import (
	"context"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type SetActuatorStateHandler struct {
	state.State
}

func (h SetActuatorStateHandler) SetActuatorState(ctx context.Context, req *pb.SetActuatorStateRequest) (*pb.SetActuatorStateResponse, error) {
	log := telemetry.GetLogger(ctx)

	if req.Feed == nil {
		log.Error("feed is nil")
		return nil, status.Error(codes.InvalidArgument, "feed cannot be nil")
	}

	log = log.With(zap.String("feed.id", req.Feed.Id))
	ctx = telemetry.ContextWithLogger(ctx, log)

	var payload []byte
	if req.State {
		payload = []byte("1")
	} else {
		payload = []byte("0")
	}

	if err := h.MQTTSubscriber().Publish(ctx, req.Feed.Id, payload); err != nil {
		return nil, err
	}

	return &pb.SetActuatorStateResponse{}, nil
}
