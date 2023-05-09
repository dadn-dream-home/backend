// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: protobuf/backend.proto

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BackendServiceClient is the client API for BackendService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BackendServiceClient interface {
	StreamSensorValues(ctx context.Context, in *StreamSensorValuesRequest, opts ...grpc.CallOption) (BackendService_StreamSensorValuesClient, error)
	StreamActuatorStates(ctx context.Context, in *StreamActuatorStatesRequest, opts ...grpc.CallOption) (BackendService_StreamActuatorStatesClient, error)
	StreamFeedsChanges(ctx context.Context, in *StreamFeedsChangesRequest, opts ...grpc.CallOption) (BackendService_StreamFeedsChangesClient, error)
	CreateFeed(ctx context.Context, in *CreateFeedRequest, opts ...grpc.CallOption) (*CreateFeedResponse, error)
	DeleteFeed(ctx context.Context, in *DeleteFeedRequest, opts ...grpc.CallOption) (*DeleteFeedResponse, error)
	SetActuatorState(ctx context.Context, in *SetActuatorStateRequest, opts ...grpc.CallOption) (*SetActuatorStateResponse, error)
	StreamNotifications(ctx context.Context, in *StreamNotificationsRequest, opts ...grpc.CallOption) (BackendService_StreamNotificationsClient, error)
}

type backendServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBackendServiceClient(cc grpc.ClientConnInterface) BackendServiceClient {
	return &backendServiceClient{cc}
}

func (c *backendServiceClient) StreamSensorValues(ctx context.Context, in *StreamSensorValuesRequest, opts ...grpc.CallOption) (BackendService_StreamSensorValuesClient, error) {
	stream, err := c.cc.NewStream(ctx, &BackendService_ServiceDesc.Streams[0], "/protobuf.BackendService/StreamSensorValues", opts...)
	if err != nil {
		return nil, err
	}
	x := &backendServiceStreamSensorValuesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BackendService_StreamSensorValuesClient interface {
	Recv() (*StreamSensorValuesResponse, error)
	grpc.ClientStream
}

type backendServiceStreamSensorValuesClient struct {
	grpc.ClientStream
}

func (x *backendServiceStreamSensorValuesClient) Recv() (*StreamSensorValuesResponse, error) {
	m := new(StreamSensorValuesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backendServiceClient) StreamActuatorStates(ctx context.Context, in *StreamActuatorStatesRequest, opts ...grpc.CallOption) (BackendService_StreamActuatorStatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &BackendService_ServiceDesc.Streams[1], "/protobuf.BackendService/StreamActuatorStates", opts...)
	if err != nil {
		return nil, err
	}
	x := &backendServiceStreamActuatorStatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BackendService_StreamActuatorStatesClient interface {
	Recv() (*StreamActuatorStatesResponse, error)
	grpc.ClientStream
}

type backendServiceStreamActuatorStatesClient struct {
	grpc.ClientStream
}

func (x *backendServiceStreamActuatorStatesClient) Recv() (*StreamActuatorStatesResponse, error) {
	m := new(StreamActuatorStatesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backendServiceClient) StreamFeedsChanges(ctx context.Context, in *StreamFeedsChangesRequest, opts ...grpc.CallOption) (BackendService_StreamFeedsChangesClient, error) {
	stream, err := c.cc.NewStream(ctx, &BackendService_ServiceDesc.Streams[2], "/protobuf.BackendService/StreamFeedsChanges", opts...)
	if err != nil {
		return nil, err
	}
	x := &backendServiceStreamFeedsChangesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BackendService_StreamFeedsChangesClient interface {
	Recv() (*StreamFeedsChangesResponse, error)
	grpc.ClientStream
}

type backendServiceStreamFeedsChangesClient struct {
	grpc.ClientStream
}

func (x *backendServiceStreamFeedsChangesClient) Recv() (*StreamFeedsChangesResponse, error) {
	m := new(StreamFeedsChangesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backendServiceClient) CreateFeed(ctx context.Context, in *CreateFeedRequest, opts ...grpc.CallOption) (*CreateFeedResponse, error) {
	out := new(CreateFeedResponse)
	err := c.cc.Invoke(ctx, "/protobuf.BackendService/CreateFeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendServiceClient) DeleteFeed(ctx context.Context, in *DeleteFeedRequest, opts ...grpc.CallOption) (*DeleteFeedResponse, error) {
	out := new(DeleteFeedResponse)
	err := c.cc.Invoke(ctx, "/protobuf.BackendService/DeleteFeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendServiceClient) SetActuatorState(ctx context.Context, in *SetActuatorStateRequest, opts ...grpc.CallOption) (*SetActuatorStateResponse, error) {
	out := new(SetActuatorStateResponse)
	err := c.cc.Invoke(ctx, "/protobuf.BackendService/SetActuatorState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendServiceClient) StreamNotifications(ctx context.Context, in *StreamNotificationsRequest, opts ...grpc.CallOption) (BackendService_StreamNotificationsClient, error) {
	stream, err := c.cc.NewStream(ctx, &BackendService_ServiceDesc.Streams[3], "/protobuf.BackendService/StreamNotifications", opts...)
	if err != nil {
		return nil, err
	}
	x := &backendServiceStreamNotificationsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BackendService_StreamNotificationsClient interface {
	Recv() (*StreamNotificationsResponse, error)
	grpc.ClientStream
}

type backendServiceStreamNotificationsClient struct {
	grpc.ClientStream
}

func (x *backendServiceStreamNotificationsClient) Recv() (*StreamNotificationsResponse, error) {
	m := new(StreamNotificationsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BackendServiceServer is the server API for BackendService service.
// All implementations must embed UnimplementedBackendServiceServer
// for forward compatibility
type BackendServiceServer interface {
	StreamSensorValues(*StreamSensorValuesRequest, BackendService_StreamSensorValuesServer) error
	StreamActuatorStates(*StreamActuatorStatesRequest, BackendService_StreamActuatorStatesServer) error
	StreamFeedsChanges(*StreamFeedsChangesRequest, BackendService_StreamFeedsChangesServer) error
	CreateFeed(context.Context, *CreateFeedRequest) (*CreateFeedResponse, error)
	DeleteFeed(context.Context, *DeleteFeedRequest) (*DeleteFeedResponse, error)
	SetActuatorState(context.Context, *SetActuatorStateRequest) (*SetActuatorStateResponse, error)
	StreamNotifications(*StreamNotificationsRequest, BackendService_StreamNotificationsServer) error
	mustEmbedUnimplementedBackendServiceServer()
}

// UnimplementedBackendServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBackendServiceServer struct {
}

func (UnimplementedBackendServiceServer) StreamSensorValues(*StreamSensorValuesRequest, BackendService_StreamSensorValuesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSensorValues not implemented")
}
func (UnimplementedBackendServiceServer) StreamActuatorStates(*StreamActuatorStatesRequest, BackendService_StreamActuatorStatesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamActuatorStates not implemented")
}
func (UnimplementedBackendServiceServer) StreamFeedsChanges(*StreamFeedsChangesRequest, BackendService_StreamFeedsChangesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamFeedsChanges not implemented")
}
func (UnimplementedBackendServiceServer) CreateFeed(context.Context, *CreateFeedRequest) (*CreateFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFeed not implemented")
}
func (UnimplementedBackendServiceServer) DeleteFeed(context.Context, *DeleteFeedRequest) (*DeleteFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFeed not implemented")
}
func (UnimplementedBackendServiceServer) SetActuatorState(context.Context, *SetActuatorStateRequest) (*SetActuatorStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetActuatorState not implemented")
}
func (UnimplementedBackendServiceServer) StreamNotifications(*StreamNotificationsRequest, BackendService_StreamNotificationsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamNotifications not implemented")
}
func (UnimplementedBackendServiceServer) mustEmbedUnimplementedBackendServiceServer() {}

// UnsafeBackendServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BackendServiceServer will
// result in compilation errors.
type UnsafeBackendServiceServer interface {
	mustEmbedUnimplementedBackendServiceServer()
}

func RegisterBackendServiceServer(s grpc.ServiceRegistrar, srv BackendServiceServer) {
	s.RegisterService(&BackendService_ServiceDesc, srv)
}

func _BackendService_StreamSensorValues_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamSensorValuesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BackendServiceServer).StreamSensorValues(m, &backendServiceStreamSensorValuesServer{stream})
}

type BackendService_StreamSensorValuesServer interface {
	Send(*StreamSensorValuesResponse) error
	grpc.ServerStream
}

type backendServiceStreamSensorValuesServer struct {
	grpc.ServerStream
}

func (x *backendServiceStreamSensorValuesServer) Send(m *StreamSensorValuesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _BackendService_StreamActuatorStates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamActuatorStatesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BackendServiceServer).StreamActuatorStates(m, &backendServiceStreamActuatorStatesServer{stream})
}

type BackendService_StreamActuatorStatesServer interface {
	Send(*StreamActuatorStatesResponse) error
	grpc.ServerStream
}

type backendServiceStreamActuatorStatesServer struct {
	grpc.ServerStream
}

func (x *backendServiceStreamActuatorStatesServer) Send(m *StreamActuatorStatesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _BackendService_StreamFeedsChanges_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamFeedsChangesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BackendServiceServer).StreamFeedsChanges(m, &backendServiceStreamFeedsChangesServer{stream})
}

type BackendService_StreamFeedsChangesServer interface {
	Send(*StreamFeedsChangesResponse) error
	grpc.ServerStream
}

type backendServiceStreamFeedsChangesServer struct {
	grpc.ServerStream
}

func (x *backendServiceStreamFeedsChangesServer) Send(m *StreamFeedsChangesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _BackendService_CreateFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServiceServer).CreateFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.BackendService/CreateFeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServiceServer).CreateFeed(ctx, req.(*CreateFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BackendService_DeleteFeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServiceServer).DeleteFeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.BackendService/DeleteFeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServiceServer).DeleteFeed(ctx, req.(*DeleteFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BackendService_SetActuatorState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetActuatorStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServiceServer).SetActuatorState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.BackendService/SetActuatorState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServiceServer).SetActuatorState(ctx, req.(*SetActuatorStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BackendService_StreamNotifications_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamNotificationsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BackendServiceServer).StreamNotifications(m, &backendServiceStreamNotificationsServer{stream})
}

type BackendService_StreamNotificationsServer interface {
	Send(*StreamNotificationsResponse) error
	grpc.ServerStream
}

type backendServiceStreamNotificationsServer struct {
	grpc.ServerStream
}

func (x *backendServiceStreamNotificationsServer) Send(m *StreamNotificationsResponse) error {
	return x.ServerStream.SendMsg(m)
}

// BackendService_ServiceDesc is the grpc.ServiceDesc for BackendService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BackendService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.BackendService",
	HandlerType: (*BackendServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFeed",
			Handler:    _BackendService_CreateFeed_Handler,
		},
		{
			MethodName: "DeleteFeed",
			Handler:    _BackendService_DeleteFeed_Handler,
		},
		{
			MethodName: "SetActuatorState",
			Handler:    _BackendService_SetActuatorState_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSensorValues",
			Handler:       _BackendService_StreamSensorValues_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamActuatorStates",
			Handler:       _BackendService_StreamActuatorStates_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamFeedsChanges",
			Handler:       _BackendService_StreamFeedsChanges_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamNotifications",
			Handler:       _BackendService_StreamNotifications_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protobuf/backend.proto",
}
