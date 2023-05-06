package services

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"reflect"

	"github.com/dadn-dream-home/x/server/handlers"
	"github.com/dadn-dream-home/x/server/interceptors"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type BackendService struct {
	pb.UnsafeBackendServiceServer

	handlers.ListFeedsHandler
	handlers.CreateFeedHandler
	handlers.DeleteFeedHandler
	handlers.StreamActuatorStatesHandler
	handlers.StreamSensorValuesHandler

	mqtt mqtt.Client
	db   *sql.DB
}

func NewBackendService(ctx context.Context) (*BackendService, error) {
	service := BackendService{}

	mqtt, err := NewMQTTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mqtt broker: %w", err)
	}
	service.mqtt = mqtt

	db, err := NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	service.db = db

	// Inject dependencies into handlers, basically:
	// service.Handler = &handlers.Handler{State: &service}

	// TODO: could be great if the list is created from
	// pb.unimplementedBackendServiceServer

	serviceValue := reflect.ValueOf(&service)
	for _, handlerType := range []reflect.Type{
		reflect.TypeOf(service.ListFeedsHandler),
		reflect.TypeOf(service.CreateFeedHandler),
		reflect.TypeOf(service.DeleteFeedHandler),
		reflect.TypeOf(service.StreamActuatorStatesHandler),
		reflect.TypeOf(service.StreamSensorValuesHandler),
	} {
		serviceHandlerValue := serviceValue.Elem().FieldByName(handlerType.Name())
		handlerValue := reflect.New(handlerType)
		handlerValue.Elem().FieldByName("State").Set(serviceValue)
		serviceHandlerValue.Set(handlerValue.Elem())
	}

	return &service, nil
}

func (s *BackendService) MQTT() mqtt.Client {
	return s.mqtt
}

func (s *BackendService) DB() *sql.DB {
	return s.db
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
