package main

import (
	"log"
	"strconv"

	pb "github.com/dadn-dream-home/x/protobuf"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type BackendGrpcService struct {
	pb.UnimplementedBackendServiceServer
	mqtt.Client
}

func NewBackendGrpcService(mqttClient mqtt.Client) *grpc.Server {
	s := grpc.NewServer()
	backendGrpc := &BackendGrpcService{
		Client: mqttClient,
	}
	pb.RegisterBackendServiceServer(s, backendGrpc)
	reflection.Register(s)

	return s
}

func (s *BackendGrpcService) StreamFeedValues(r *pb.StreamFeedValuesRequest, stream pb.BackendService_StreamFeedValuesServer) error {
	// TODO: check if feed exists

	subscribeToken := s.Subscribe(r.FeedId, 0, func(c mqtt.Client, m mqtt.Message) {
		// print message
		log.Printf("received message from feed %v: %s\n", r.FeedId, string(m.Payload()))

		value, err := strconv.ParseFloat(string(m.Payload()), 32)
		if err != nil {
			log.Printf("warning: failed to parse payload from feed %v: %s", r.FeedId, string(m.Payload()))
		}

		stream.Send(&pb.StreamFeedValuesResponse{
			Value: float32(value),
		})

	});

	if subscribeToken.Wait() && subscribeToken.Error() != nil {
		return subscribeToken.Error()
	}

	<-stream.Context().Done()

	log.Printf("finished streaming feed %v\n", r.FeedId)

	return stream.Context().Err()
}
