package startup

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/dadn-dream-home/x/server/handlers"
	"github.com/dadn-dream-home/x/server/interceptors"
	"github.com/dadn-dream-home/x/server/services"
	"github.com/dadn-dream-home/x/server/services/repository"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type Server struct {
	pb.UnsafeBackendServiceServer
	grpcServer *grpc.Server

	handlers.CreateFeedHandler
	handlers.DeleteFeedHandler
	handlers.StreamActuatorStatesHandler
	handlers.StreamSensorValuesHandler
	handlers.StreamFeedsChangesHandler
	handlers.SetActuatorStateHandler
	handlers.StreamNotificationsHandler
	handlers.GetFeedConfigHandler
	handlers.UpdateFeedConfigHandler
	handlers.StreamActivitiesHandler

	repository         state.Repository
	databaseHooker     state.DatabaseHooker
	mqttSubscriber     state.MQTTSubscriber
	feedValuesLogger   state.FeedValuesLogger
	automaticScheduler state.AutomaticScheduler
}

var _ state.State = (*Server)(nil)

func NewServer(ctx context.Context, db *sql.DB, mqtt mqtt.Client, hooker state.DatabaseHooker) *Server {
	s := &Server{}
	s.databaseHooker = hooker
	s.repository = repository.NewRepository(ctx, s, db)
	s.mqttSubscriber = services.NewMQTTSubscriber(s, mqtt)
	s.feedValuesLogger = services.NewFeedValuesLogger(s)
	s.automaticScheduler = services.NewAutomaticScheduler(s)

	// Inject dependencies into handlers, basically:
	// service.Handler = &handlers.Handler{State: &service}

	// TODO: could be great if the list is created from
	// pb.unimplementedBackendServiceServer

	serviceValue := reflect.ValueOf(s)
	for _, handlerType := range []reflect.Type{
		reflect.TypeOf(s.CreateFeedHandler),
		reflect.TypeOf(s.DeleteFeedHandler),
		reflect.TypeOf(s.StreamActuatorStatesHandler),
		reflect.TypeOf(s.StreamSensorValuesHandler),
		reflect.TypeOf(s.StreamFeedsChangesHandler),
		reflect.TypeOf(s.SetActuatorStateHandler),
		reflect.TypeOf(s.StreamNotificationsHandler),
		reflect.TypeOf(s.GetFeedConfigHandler),
		reflect.TypeOf(s.UpdateFeedConfigHandler),
		reflect.TypeOf(s.StreamActivitiesHandler),
	} {
		serviceHandlerValue := serviceValue.Elem().FieldByName(handlerType.Name())
		handlerValue := reflect.New(handlerType)
		handlerValue.Elem().FieldByName("State").Set(serviceValue)
		serviceHandlerValue.Set(handlerValue.Elem())
	}

	return s
}

func (s *Server) DatabaseHooker() state.DatabaseHooker {
	return s.databaseHooker
}

func (s *Server) MQTTSubscriber() state.MQTTSubscriber {
	return s.mqttSubscriber
}

func (s *Server) FeedValuesLogger() state.FeedValuesLogger {
	return s.feedValuesLogger
}

func (s *Server) AutomaticScheduler() state.AutomaticScheduler {
	return s.automaticScheduler
}

func (s *Server) Repository() state.Repository {
	return s.repository
}

func (s *Server) Serve(ctx context.Context, lis net.Listener) {
	log := telemetry.GetLogger(ctx)

	s.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.AddLoggerUnaryInterceptor,
		interceptors.RequestIdUnaryInterceptor,
	), grpc.ChainStreamInterceptor(
		interceptors.AddLoggerStreamInterceptor,
		interceptors.RequestIdStreamInterceptor,
	))
	pb.RegisterBackendServiceServer(s.grpcServer, s)
	reflection.Register(s.grpcServer)

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.mqttSubscriber.Serve(ctx)
	})
	group.Go(func() error {
		return s.feedValuesLogger.Serve(ctx)
	})
	group.Go(func() error {
		return s.automaticScheduler.Serve(ctx)
	})

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			group.Go(func() error { return err })
		} else {
			group.Go(func() error { return context.Canceled })
		}
	}()

	if err := group.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		s.grpcServer.Stop()
		log.Fatal("server error", zap.Error(err))
	}
}

func (s *Server) Stop() {
	s.grpcServer.Stop()
}

func Listen(ctx context.Context, config ServerConfig) net.Listener {
	log := telemetry.GetLogger(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", config.Port))
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	log.Info("server listening", zap.String("address", lis.Addr().String()))

	return lis
}
