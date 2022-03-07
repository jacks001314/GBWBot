// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

import (
	context "context"
	model "github.com/cbot/client/model"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LogStreamServiceClient is the client API for LogStreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogStreamServiceClient interface {
	Open(ctx context.Context, in *model.TargetNode, opts ...grpc.CallOption) (LogStreamService_OpenClient, error)
	Close(ctx context.Context, in *model.TargetNode, opts ...grpc.CallOption) (*model.OPStatus, error)
	CmdChannel(ctx context.Context, in *model.TargetNode, opts ...grpc.CallOption) (LogStreamService_CmdChannelClient, error)
	Channel(ctx context.Context, opts ...grpc.CallOption) (LogStreamService_ChannelClient, error)
}

type logStreamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogStreamServiceClient(cc grpc.ClientConnInterface) LogStreamServiceClient {
	return &logStreamServiceClient{cc}
}

func (c *logStreamServiceClient) Open(ctx context.Context, in *model.TargetNode, opts ...grpc.CallOption) (LogStreamService_OpenClient, error) {
	stream, err := c.cc.NewStream(ctx, &LogStreamService_ServiceDesc.Streams[0], "/sbot.proto.service.LogStreamService/Open", opts...)
	if err != nil {
		return nil, err
	}
	x := &logStreamServiceOpenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LogStreamService_OpenClient interface {
	Recv() (*model.LogStream, error)
	grpc.ClientStream
}

type logStreamServiceOpenClient struct {
	grpc.ClientStream
}

func (x *logStreamServiceOpenClient) Recv() (*model.LogStream, error) {
	m := new(model.LogStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *logStreamServiceClient) Close(ctx context.Context, in *model.TargetNode, opts ...grpc.CallOption) (*model.OPStatus, error) {
	out := new(model.OPStatus)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.LogStreamService/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logStreamServiceClient) CmdChannel(ctx context.Context, in *model.TargetNode, opts ...grpc.CallOption) (LogStreamService_CmdChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &LogStreamService_ServiceDesc.Streams[1], "/sbot.proto.service.LogStreamService/CmdChannel", opts...)
	if err != nil {
		return nil, err
	}
	x := &logStreamServiceCmdChannelClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LogStreamService_CmdChannelClient interface {
	Recv() (*model.LogCmd, error)
	grpc.ClientStream
}

type logStreamServiceCmdChannelClient struct {
	grpc.ClientStream
}

func (x *logStreamServiceCmdChannelClient) Recv() (*model.LogCmd, error) {
	m := new(model.LogCmd)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *logStreamServiceClient) Channel(ctx context.Context, opts ...grpc.CallOption) (LogStreamService_ChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &LogStreamService_ServiceDesc.Streams[2], "/sbot.proto.service.LogStreamService/Channel", opts...)
	if err != nil {
		return nil, err
	}
	x := &logStreamServiceChannelClient{stream}
	return x, nil
}

type LogStreamService_ChannelClient interface {
	Send(*model.LogStream) error
	CloseAndRecv() (*model.OPStatus, error)
	grpc.ClientStream
}

type logStreamServiceChannelClient struct {
	grpc.ClientStream
}

func (x *logStreamServiceChannelClient) Send(m *model.LogStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *logStreamServiceChannelClient) CloseAndRecv() (*model.OPStatus, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(model.OPStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LogStreamServiceServer is the server API for LogStreamService service.
// All implementations must embed UnimplementedLogStreamServiceServer
// for forward compatibility
type LogStreamServiceServer interface {
	Open(*model.TargetNode, LogStreamService_OpenServer) error
	Close(context.Context, *model.TargetNode) (*model.OPStatus, error)
	CmdChannel(*model.TargetNode, LogStreamService_CmdChannelServer) error
	Channel(LogStreamService_ChannelServer) error
	mustEmbedUnimplementedLogStreamServiceServer()
}

// UnimplementedLogStreamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLogStreamServiceServer struct {
}

func (UnimplementedLogStreamServiceServer) Open(*model.TargetNode, LogStreamService_OpenServer) error {
	return status.Errorf(codes.Unimplemented, "method Open not implemented")
}
func (UnimplementedLogStreamServiceServer) Close(context.Context, *model.TargetNode) (*model.OPStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedLogStreamServiceServer) CmdChannel(*model.TargetNode, LogStreamService_CmdChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method CmdChannel not implemented")
}
func (UnimplementedLogStreamServiceServer) Channel(LogStreamService_ChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method Channel not implemented")
}
func (UnimplementedLogStreamServiceServer) mustEmbedUnimplementedLogStreamServiceServer() {}

// UnsafeLogStreamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogStreamServiceServer will
// result in compilation errors.
type UnsafeLogStreamServiceServer interface {
	mustEmbedUnimplementedLogStreamServiceServer()
}

func RegisterLogStreamServiceServer(s grpc.ServiceRegistrar, srv LogStreamServiceServer) {
	s.RegisterService(&LogStreamService_ServiceDesc, srv)
}

func _LogStreamService_Open_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(model.TargetNode)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LogStreamServiceServer).Open(m, &logStreamServiceOpenServer{stream})
}

type LogStreamService_OpenServer interface {
	Send(*model.LogStream) error
	grpc.ServerStream
}

type logStreamServiceOpenServer struct {
	grpc.ServerStream
}

func (x *logStreamServiceOpenServer) Send(m *model.LogStream) error {
	return x.ServerStream.SendMsg(m)
}

func _LogStreamService_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.TargetNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogStreamServiceServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.LogStreamService/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogStreamServiceServer).Close(ctx, req.(*model.TargetNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogStreamService_CmdChannel_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(model.TargetNode)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LogStreamServiceServer).CmdChannel(m, &logStreamServiceCmdChannelServer{stream})
}

type LogStreamService_CmdChannelServer interface {
	Send(*model.LogCmd) error
	grpc.ServerStream
}

type logStreamServiceCmdChannelServer struct {
	grpc.ServerStream
}

func (x *logStreamServiceCmdChannelServer) Send(m *model.LogCmd) error {
	return x.ServerStream.SendMsg(m)
}

func _LogStreamService_Channel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LogStreamServiceServer).Channel(&logStreamServiceChannelServer{stream})
}

type LogStreamService_ChannelServer interface {
	SendAndClose(*model.OPStatus) error
	Recv() (*model.LogStream, error)
	grpc.ServerStream
}

type logStreamServiceChannelServer struct {
	grpc.ServerStream
}

func (x *logStreamServiceChannelServer) SendAndClose(m *model.OPStatus) error {
	return x.ServerStream.SendMsg(m)
}

func (x *logStreamServiceChannelServer) Recv() (*model.LogStream, error) {
	m := new(model.LogStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// LogStreamService_ServiceDesc is the grpc.ServiceDesc for LogStreamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogStreamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sbot.proto.service.LogStreamService",
	HandlerType: (*LogStreamServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Close",
			Handler:    _LogStreamService_Close_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Open",
			Handler:       _LogStreamService_Open_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "CmdChannel",
			Handler:       _LogStreamService_CmdChannel_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Channel",
			Handler:       _LogStreamService_Channel_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/service/logstream_service.proto",
}