package services

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/dadn-dream-home/x/server/handlers"
	"github.com/dadn-dream-home/x/server/interceptors"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type BackendService struct {
	pb.UnsafeBackendServiceServer

	handlers.CreateFeedHandler
	handlers.DeleteFeedHandler
	handlers.StreamActuatorStatesHandler
	handlers.StreamSensorValuesHandler
	handlers.StreamFeedsChangesHandler
	handlers.SetActuatorStateHandler

	pubSubValues state.PubSubValues
	pubSubFeeds  state.PubSubFeeds
	repository   state.Repository
}

var _ state.State = (*BackendService)(nil)

func NewBackendService(ctx context.Context) (service *BackendService, err error) {
	service = &BackendService{}

	service.repository, err = NewRepository(ctx, service)
	if err != nil {
		return nil, err
	}

	service.pubSubFeeds = NewPubSubFeeds(ctx, service)

	service.pubSubValues, err = NewPubSubValues(ctx, service)
	if err != nil {
		return nil, err
	}

	// Inject dependencies into handlers, basically:
	// service.Handler = &handlers.Handler{State: &service}

	// TODO: could be great if the list is created from
	// pb.unimplementedBackendServiceServer

	serviceValue := reflect.ValueOf(service)
	for _, handlerType := range []reflect.Type{
		reflect.TypeOf(service.CreateFeedHandler),
		reflect.TypeOf(service.DeleteFeedHandler),
		reflect.TypeOf(service.StreamActuatorStatesHandler),
		reflect.TypeOf(service.StreamSensorValuesHandler),
		reflect.TypeOf(service.StreamFeedsChangesHandler),
		reflect.TypeOf(service.SetActuatorStateHandler),
	} {
		serviceHandlerValue := serviceValue.Elem().FieldByName(handlerType.Name())
		handlerValue := reflect.New(handlerType)
		handlerValue.Elem().FieldByName("State").Set(serviceValue)
		serviceHandlerValue.Set(handlerValue.Elem())
	}

	return service, nil
}

func (s *BackendService) PubSubFeeds() state.PubSubFeeds {
	return s.pubSubFeeds
}

func (s *BackendService) PubSubValues() state.PubSubValues {
	return s.pubSubValues
}

func (s *BackendService) Repository() state.Repository {
	return s.repository
}

func (s *BackendService) Serve(ctx context.Context) {
	log := telemetry.GetLogger(ctx)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.AddLoggerUnaryInterceptor,
		interceptors.RequestIdUnaryInterceptor,
	), grpc.ChainStreamInterceptor(
		interceptors.AddLoggerStreamInterceptor,
		interceptors.RequestIdStreamInterceptor,
	))
	pb.RegisterBackendServiceServer(server, s)
	reflection.Register(server)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 50051))
	if err != nil {
		panic(err)
	}
	log.WithField("addr", lis.Addr()).Infof("server listening")

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
