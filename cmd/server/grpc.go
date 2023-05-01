package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"

	pb "github.com/dadn-dream-home/x/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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

func NewBackendGrpcService() (*BackendGrpcService, error) {
	database, err := NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	defer func() {
		if err != nil {
			closeErr := database.Close()
			if closeErr != nil {
				closeErr := fmt.Errorf("failed to close database: %w", closeErr)
				err = errors.Join(err, closeErr)
			}
		}
	}()

	service := &BackendGrpcService{
		MQTTClient: NewMQTTClient(),
		Database:   database,
		Server:     grpc.NewServer(),
	}
	pb.RegisterBackendServiceServer(service.Server, service)
	reflection.Register(service.Server)

	return service, nil
}

func (s *BackendGrpcService) Init() error {
	// get feeds from database and subscribe to them
	feeds, err := s.Database.ListFeeds()
	if err != nil {
		return fmt.Errorf("failed to list feeds: %w", err)
	}

	s.feedChannels = make(map[string]struct {
		tx chan<- []byte
		rx <-chan []byte
	}, len(feeds))
	for _, feed := range feeds {
		s.Subscribe(feed.Id)
	}

	return nil
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
	feed, err := s.Database.InsertFeed(r.Id, r.Type)
	if err != nil {
		log.WithError(err).Errorf("failed to insert feed")
		if errors.Is(err, ErrorFeedExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, err
	}

	s.Subscribe(feed.Id)

	return &pb.CreateFeedResponse{
		Id: feed.Id,
	}, nil
}
