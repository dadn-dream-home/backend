package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

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
		Server: grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				requestIdUnaryInterceptor,
			),
			grpc.ChainStreamInterceptor(
				requestIdStreamInterceptor,
			),
		),
	}
	pb.RegisterBackendServiceServer(service.Server, service)
	reflection.Register(service.Server)

	return service, nil
}

func (s *BackendGrpcService) Init(ctx context.Context) error {
	// get feeds from database and subscribe to them
	feeds, err := s.Database.ListFeeds(ctx)
	if err != nil {
		return fmt.Errorf("failed to list feeds: %w", err)
	}

	s.feedChannels = make(map[string]struct {
		tx chan<- []byte
		rx <-chan []byte
	}, len(feeds))
	for _, feed := range feeds {
		s.Subscribe(ctx, feed.Id)
	}

	return nil
}

func (s *BackendGrpcService) Subscribe(ctx context.Context, topic string) {
	tx, rx := s.MQTTClient.Subscribe(ctx, topic)
	s.feedChannels[topic] = struct {
		tx chan<- []byte
		rx <-chan []byte
	}{tx, rx}
}

func (s *BackendGrpcService) Unsubscribe(ctx context.Context, topic string) {
	channels := s.feedChannels[topic]
	close(channels.tx)
	delete(s.feedChannels, topic)

	log.WithField("feed", topic).Infof("unsubscribed\n")
}

func (s *BackendGrpcService) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		log.WithError(err).Fatalf("failed to listen")
	}

	log.Printf("server listening at %v\n", lis.Addr())
	if err := s.Server.Serve(lis); err != nil {
		log.WithError(err).Fatalf("failed to serve")
	}
}

func (s *BackendGrpcService) StreamFeedValues(r *pb.StreamFeedValuesRequest, stream pb.BackendService_StreamFeedValuesServer) error {
	log := LoggerFromContext(stream.Context())
	log = log.WithField("feed", r.Id)

	ch, ok := s.feedChannels[r.Id]
	if !ok {
		panic("feed not exist")
	}

	log.Infof("begin streaming\n")

outer:
	for {
		select {
		case payload, ok := <-ch.rx:
			if !ok {
				log.Warnln("channel closed")
				break outer
			}

			log = log.WithField("payload", string(payload))
			log.Infof("received payload\n")

			value, err := strconv.ParseFloat(string(payload), 32)
			if err != nil {
				log.Warnf("failed to parse payload\n")
				continue
			}

			stream.Send(&pb.StreamFeedValuesResponse{
				Value: float32(value),
			})
		case <-stream.Context().Done():
			break outer
		}
	}

	log.Infof("finished streaming\n")

	return stream.Context().Err()
}

func (s *BackendGrpcService) CreateFeed(ctx context.Context, r *pb.CreateFeedRequest) (*pb.CreateFeedResponse, error) {
	log := log.WithField("feed", r.Id)

	feed, err := s.Database.InsertFeed(ctx, r.Id, r.Type)
	if err != nil {
		log.WithError(err).Errorf("failed to insert feed")
		if errors.Is(err, ErrorFeedExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	s.Subscribe(ctx, feed.Id)

	return &pb.CreateFeedResponse{
		Id: feed.Id,
	}, nil
}

func (s *BackendGrpcService) ListFeeds(ctx context.Context, r *pb.ListFeedsRequest) (*pb.ListFeedsResponse, error) {
	feeds, err := s.Database.ListFeeds(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListFeedsResponse{Feeds: feeds}, nil
}

func (s *BackendGrpcService) DeleteFeed(ctx context.Context, r *pb.DeleteFeedRequest) (*pb.DeleteFeedResponse, error) {
	err := s.Database.DeleteFeed(ctx, r.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	s.Unsubscribe(ctx, r.Id)

	return &pb.DeleteFeedResponse{}, nil
}
