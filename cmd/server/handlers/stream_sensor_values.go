package handlers

import (

	pb "github.com/dadn-dream-home/x/protobuf"

	"github.com/dadn-dream-home/x/server/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamSensorValuesHandler struct {
	state.State
}

func (s StreamSensorValuesHandler) StreamSensorValues(*pb.StreamFeedValuesRequest, pb.BackendService_StreamSensorValuesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSensorValues not implemented")
}
