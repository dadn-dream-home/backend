package handlers

import (
	"context"

	"github.com/dadn-dream-home/x/server/state"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type GetFeedConfigHandler struct {
	state.State
}

func (h GetFeedConfigHandler) GetFeedConfig(ctx context.Context, req *pb.GetFeedConfigRequest) (*pb.GetFeedConfigResponse, error) {
	config, err := h.Repository().GetFeedConfig(ctx, req.Feed.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetFeedConfigResponse{Config: config}, err
}
