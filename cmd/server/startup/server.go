package startup

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"reflect"

	"github.com/dadn-dream-home/x/server/handlers"
	"github.com/dadn-dream-home/x/server/interceptors"
	"github.com/dadn-dream-home/x/server/services"
	"github.com/dadn-dream-home/x/server/state"
	"github.com/dadn-dream-home/x/server/telemetry"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/dadn-dream-home/x/protobuf"
)

type Server struct {
	pb.UnsafeBackendServiceServer

	handlers.CreateFeedHandler
	handlers.DeleteFeedHandler
	handlers.StreamActuatorStatesHandler
	handlers.StreamSensorValuesHandler
	handlers.StreamFeedsChangesHandler
	handlers.SetActuatorStateHandler
	handlers.StreamNotificationsHandler

	pubSubValues state.PubSubValues
	pubSubFeeds  state.PubSubFeeds
	repository   state.Repository
	notifier     state.Notifier
}

var _ state.State = (*Server)(nil)

func NewServer(ctx context.Context, db *sql.DB, mqtt mqtt.Client) *Server {
	s := &Server{}
	s.repository = services.NewRepository(ctx, s, db)
	s.pubSubFeeds = services.NewPubSubFeeds(ctx, s)
	s.pubSubValues = services.NewPubSubValues(ctx, s, mqtt)
	s.notifier = services.NewNotifier(ctx, s)

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
	} {
		serviceHandlerValue := serviceValue.Elem().FieldByName(handlerType.Name())
		handlerValue := reflect.New(handlerType)
		handlerValue.Elem().FieldByName("State").Set(serviceValue)
		serviceHandlerValue.Set(handlerValue.Elem())
	}

	return s
}

func (s *Server) PubSubFeeds() state.PubSubFeeds {
	return s.pubSubFeeds
}

func (s *Server) PubSubValues() state.PubSubValues {
	return s.pubSubValues
}

func (s *Server) Repository() state.Repository {
	return s.repository
}

func (s *Server) Notifier() state.Notifier {
	return s.notifier
}

func (s *Server) Serve(ctx context.Context, lis net.Listener) {
	log := telemetry.GetLogger(ctx)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptors.AddLoggerUnaryInterceptor,
		interceptors.RequestIdUnaryInterceptor,
	), grpc.ChainStreamInterceptor(
		interceptors.AddLoggerStreamInterceptor,
		interceptors.RequestIdStreamInterceptor,
	))
	pb.RegisterBackendServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
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
