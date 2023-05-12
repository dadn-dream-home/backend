package handlers

import (
	"context"

	"github.com/dadn-dream-home/x/server/state"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type UpdateFeedConfigHandler struct {
	state.State
}

func (h UpdateFeedConfigHandler) UpdateFeedConfig(ctx context.Context, req *pb.UpdateFeedConfigRequest) (*pb.UpdateFeedConfigResponse, error) {
	if err := h.Repository().UpdateFeedConfig(ctx, req.Config); err != nil {
		return nil, err
	}
	return &pb.UpdateFeedConfigResponse{}, nil
}
