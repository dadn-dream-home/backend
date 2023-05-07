package handlers

import (
	"context"

	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type SetActuatorStateHandler struct {
	state.State
}

func (h SetActuatorStateHandler) SetActuatorState(ctx context.Context, req *pb.SetActuatorStateRequest) (*pb.SetActuatorStateResponse, error) {
	log := telemetry.GetLogger(ctx).WithField("feed_id", req.Feed.Id)
	ctx = telemetry.ContextWithLogger(ctx, log)
	rid := telemetry.GetRequestId(ctx)

	var payload []byte
	if req.State {
		payload = []byte("1")
	} else {
		payload = []byte("0")
	}

	err := h.PubSubValues().Publish(ctx, rid, req.Feed.Id, payload)
	if err != nil {
		log.WithError(err).Errorf("error setting actuator state")
		return nil, err
	}

	return &pb.SetActuatorStateResponse{}, nil
}
