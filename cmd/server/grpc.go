package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "github.com/dadn-dream-home/x/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type BackendGrpcService struct {
	pb.UnimplementedBackendServiceServer

	Database   Database
	MQTTClient MQTTClient
	Server     *grpc.Server

	feedChannels map[string]struct {
		tx chan<- []byte
		rx <-chan []byte
	}
}

func NewBackendGrpcService() *BackendGrpcService {
	service := &BackendGrpcService{
		MQTTClient: NewMQTTClient(),
		Database:   NewDatabase(),
		Server:     grpc.NewServer(),
	}
	pb.RegisterBackendServiceServer(service.Server, service)
	reflection.Register(service.Server)

	return service
}

func (s *BackendGrpcService) Init() {
	// get feeds from database and subscribe to them
	feeds := s.Database.ListFeeds()
	s.feedChannels = make(map[string]struct {
		tx chan<- []byte
		rx <-chan []byte
	}, len(feeds))
	for _, feed := range feeds {
		s.Subscribe(feed.Id)
	}
}

func (s *BackendGrpcService) Subscribe(topic string) {
	tx, rx := s.MQTTClient.Subscribe(topic)
	s.feedChannels[topic] = struct {
		tx chan<- []byte
		rx <-chan []byte
	}{tx, rx}
}

func (s *BackendGrpcService) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *BackendGrpcService) StreamFeedValues(r *pb.StreamFeedValuesRequest, stream pb.BackendService_StreamFeedValuesServer) error {
	// TODO: check if feed exists

	ch, ok := s.feedChannels[r.Id]
	if !ok {
		panic("channels not exist")
	}

outer:
	for {
		select {
		case payload, ok := <-ch.rx:
			if !ok {
				log.Printf("warning: channel for feed %v closed\n", r.Id)
				break outer
			}

			log.Printf("received message from feed %v: %s\n", r.Id, string(payload))

			value, err := strconv.ParseFloat(string(payload), 32)
			if err != nil {
				log.Printf("warning: failed to parse payload from feed %v: %s", r.Id, string(payload))
				continue
			}

			stream.Send(&pb.StreamFeedValuesResponse{
				Value: float32(value),
			})
		case <-stream.Context().Done():
			break outer
		}
	}

	log.Printf("finished streaming feed %v\n", r.Id)

	return stream.Context().Err()
}

func (s *BackendGrpcService) CreateFeed(ctx context.Context, r *pb.CreateFeedRequest) (*pb.CreateFeedResponse, error) {
	feed := s.Database.InsertFeed(r.Id, r.Type)

	s.Subscribe(feed.Id)

	return &pb.CreateFeedResponse{
		Id: feed.Id,
	}, nil
}
