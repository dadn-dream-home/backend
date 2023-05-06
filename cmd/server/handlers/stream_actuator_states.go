package handlers

import (

	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/dadn-dream-home/x/server/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamActuatorStatesHandler struct {
	state.State
}

func (s StreamActuatorStatesHandler) StreamActuatorStates(*pb.StreamActuatorStatesRequest, pb.BackendService_StreamActuatorStatesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamActuatorStates not implemented")
}
