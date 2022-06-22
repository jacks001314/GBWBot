// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package service

import (
	context "context"
	model "github.com/sbot/proto/model"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SbotServiceClient is the client API for SbotService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SbotServiceClient interface {
	QueryAttackTasks(ctx context.Context, in *model.AttackTaskQuery, opts ...grpc.CallOption) (*model.AttackTaskReply, error)
	FacetAttackTasks(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error)
	CountAttackTasks(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error)
	QueryAttackedNodes(ctx context.Context, in *model.AttackedNodeQuery, opts ...grpc.CallOption) (*model.AttackedNodeReply, error)
	FacetAttackedNodes(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error)
	CountAttackedNodes(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error)
	QueryAttackProcess(ctx context.Context, in *model.AttackProcessQuery, opts ...grpc.CallOption) (*model.AttackProcessMessageReply, error)
	FacetAttackProcess(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error)
	CountAttackProcess(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error)
	QueryAttackedDownloadFiles(ctx context.Context, in *model.AttackedNodeDownloadFileQuery, opts ...grpc.CallOption) (*model.AttackedNodeDownloadFileReply, error)
	FacetAttackedDownloadFiles(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error)
	CountAttackedDownloadFiles(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error)
}

type sbotServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSbotServiceClient(cc grpc.ClientConnInterface) SbotServiceClient {
	return &sbotServiceClient{cc}
}

func (c *sbotServiceClient) QueryAttackTasks(ctx context.Context, in *model.AttackTaskQuery, opts ...grpc.CallOption) (*model.AttackTaskReply, error) {
	out := new(model.AttackTaskReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/QueryAttackTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) FacetAttackTasks(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error) {
	out := new(model.FacetReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/FacetAttackTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) CountAttackTasks(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error) {
	out := new(model.Count)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/CountAttackTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) QueryAttackedNodes(ctx context.Context, in *model.AttackedNodeQuery, opts ...grpc.CallOption) (*model.AttackedNodeReply, error) {
	out := new(model.AttackedNodeReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/QueryAttackedNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) FacetAttackedNodes(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error) {
	out := new(model.FacetReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/FacetAttackedNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) CountAttackedNodes(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error) {
	out := new(model.Count)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/CountAttackedNodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) QueryAttackProcess(ctx context.Context, in *model.AttackProcessQuery, opts ...grpc.CallOption) (*model.AttackProcessMessageReply, error) {
	out := new(model.AttackProcessMessageReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/QueryAttackProcess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) FacetAttackProcess(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error) {
	out := new(model.FacetReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/FacetAttackProcess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) CountAttackProcess(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error) {
	out := new(model.Count)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/CountAttackProcess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) QueryAttackedDownloadFiles(ctx context.Context, in *model.AttackedNodeDownloadFileQuery, opts ...grpc.CallOption) (*model.AttackedNodeDownloadFileReply, error) {
	out := new(model.AttackedNodeDownloadFileReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/QueryAttackedDownloadFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) FacetAttackedDownloadFiles(ctx context.Context, in *model.FacetRequest, opts ...grpc.CallOption) (*model.FacetReply, error) {
	out := new(model.FacetReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/FacetAttackedDownloadFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sbotServiceClient) CountAttackedDownloadFiles(ctx context.Context, in *model.CountRequest, opts ...grpc.CallOption) (*model.Count, error) {
	out := new(model.Count)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.SbotService/CountAttackedDownloadFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SbotServiceServer is the server API for SbotService service.
// All implementations must embed UnimplementedSbotServiceServer
// for forward compatibility
type SbotServiceServer interface {
	QueryAttackTasks(context.Context, *model.AttackTaskQuery) (*model.AttackTaskReply, error)
	FacetAttackTasks(context.Context, *model.FacetRequest) (*model.FacetReply, error)
	CountAttackTasks(context.Context, *model.CountRequest) (*model.Count, error)
	QueryAttackedNodes(context.Context, *model.AttackedNodeQuery) (*model.AttackedNodeReply, error)
	FacetAttackedNodes(context.Context, *model.FacetRequest) (*model.FacetReply, error)
	CountAttackedNodes(context.Context, *model.CountRequest) (*model.Count, error)
	QueryAttackProcess(context.Context, *model.AttackProcessQuery) (*model.AttackProcessMessageReply, error)
	FacetAttackProcess(context.Context, *model.FacetRequest) (*model.FacetReply, error)
	CountAttackProcess(context.Context, *model.CountRequest) (*model.Count, error)
	QueryAttackedDownloadFiles(context.Context, *model.AttackedNodeDownloadFileQuery) (*model.AttackedNodeDownloadFileReply, error)
	FacetAttackedDownloadFiles(context.Context, *model.FacetRequest) (*model.FacetReply, error)
	CountAttackedDownloadFiles(context.Context, *model.CountRequest) (*model.Count, error)
	mustEmbedUnimplementedSbotServiceServer()
}

// UnimplementedSbotServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSbotServiceServer struct {
}

func (UnimplementedSbotServiceServer) QueryAttackTasks(context.Context, *model.AttackTaskQuery) (*model.AttackTaskReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryAttackTasks not implemented")
}
func (UnimplementedSbotServiceServer) FacetAttackTasks(context.Context, *model.FacetRequest) (*model.FacetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FacetAttackTasks not implemented")
}
func (UnimplementedSbotServiceServer) CountAttackTasks(context.Context, *model.CountRequest) (*model.Count, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountAttackTasks not implemented")
}
func (UnimplementedSbotServiceServer) QueryAttackedNodes(context.Context, *model.AttackedNodeQuery) (*model.AttackedNodeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryAttackedNodes not implemented")
}
func (UnimplementedSbotServiceServer) FacetAttackedNodes(context.Context, *model.FacetRequest) (*model.FacetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FacetAttackedNodes not implemented")
}
func (UnimplementedSbotServiceServer) CountAttackedNodes(context.Context, *model.CountRequest) (*model.Count, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountAttackedNodes not implemented")
}
func (UnimplementedSbotServiceServer) QueryAttackProcess(context.Context, *model.AttackProcessQuery) (*model.AttackProcessMessageReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryAttackProcess not implemented")
}
func (UnimplementedSbotServiceServer) FacetAttackProcess(context.Context, *model.FacetRequest) (*model.FacetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FacetAttackProcess not implemented")
}
func (UnimplementedSbotServiceServer) CountAttackProcess(context.Context, *model.CountRequest) (*model.Count, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountAttackProcess not implemented")
}
func (UnimplementedSbotServiceServer) QueryAttackedDownloadFiles(context.Context, *model.AttackedNodeDownloadFileQuery) (*model.AttackedNodeDownloadFileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryAttackedDownloadFiles not implemented")
}
func (UnimplementedSbotServiceServer) FacetAttackedDownloadFiles(context.Context, *model.FacetRequest) (*model.FacetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FacetAttackedDownloadFiles not implemented")
}
func (UnimplementedSbotServiceServer) CountAttackedDownloadFiles(context.Context, *model.CountRequest) (*model.Count, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountAttackedDownloadFiles not implemented")
}
func (UnimplementedSbotServiceServer) mustEmbedUnimplementedSbotServiceServer() {}

// UnsafeSbotServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SbotServiceServer will
// result in compilation errors.
type UnsafeSbotServiceServer interface {
	mustEmbedUnimplementedSbotServiceServer()
}

func RegisterSbotServiceServer(s grpc.ServiceRegistrar, srv SbotServiceServer) {
	s.RegisterService(&SbotService_ServiceDesc, srv)
}

func _SbotService_QueryAttackTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.AttackTaskQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).QueryAttackTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/QueryAttackTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).QueryAttackTasks(ctx, req.(*model.AttackTaskQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_FacetAttackTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.FacetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).FacetAttackTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/FacetAttackTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).FacetAttackTasks(ctx, req.(*model.FacetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_CountAttackTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.CountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).CountAttackTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/CountAttackTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).CountAttackTasks(ctx, req.(*model.CountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_QueryAttackedNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.AttackedNodeQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).QueryAttackedNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/QueryAttackedNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).QueryAttackedNodes(ctx, req.(*model.AttackedNodeQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_FacetAttackedNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.FacetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).FacetAttackedNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/FacetAttackedNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).FacetAttackedNodes(ctx, req.(*model.FacetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_CountAttackedNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.CountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).CountAttackedNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/CountAttackedNodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).CountAttackedNodes(ctx, req.(*model.CountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_QueryAttackProcess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.AttackProcessQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).QueryAttackProcess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/QueryAttackProcess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).QueryAttackProcess(ctx, req.(*model.AttackProcessQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_FacetAttackProcess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.FacetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).FacetAttackProcess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/FacetAttackProcess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).FacetAttackProcess(ctx, req.(*model.FacetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_CountAttackProcess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.CountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).CountAttackProcess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/CountAttackProcess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).CountAttackProcess(ctx, req.(*model.CountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_QueryAttackedDownloadFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.AttackedNodeDownloadFileQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).QueryAttackedDownloadFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/QueryAttackedDownloadFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).QueryAttackedDownloadFiles(ctx, req.(*model.AttackedNodeDownloadFileQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_FacetAttackedDownloadFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.FacetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).FacetAttackedDownloadFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/FacetAttackedDownloadFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).FacetAttackedDownloadFiles(ctx, req.(*model.FacetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SbotService_CountAttackedDownloadFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.CountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SbotServiceServer).CountAttackedDownloadFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.SbotService/CountAttackedDownloadFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SbotServiceServer).CountAttackedDownloadFiles(ctx, req.(*model.CountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SbotService_ServiceDesc is the grpc.ServiceDesc for SbotService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SbotService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sbot.proto.service.SbotService",
	HandlerType: (*SbotServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryAttackTasks",
			Handler:    _SbotService_QueryAttackTasks_Handler,
		},
		{
			MethodName: "FacetAttackTasks",
			Handler:    _SbotService_FacetAttackTasks_Handler,
		},
		{
			MethodName: "CountAttackTasks",
			Handler:    _SbotService_CountAttackTasks_Handler,
		},
		{
			MethodName: "QueryAttackedNodes",
			Handler:    _SbotService_QueryAttackedNodes_Handler,
		},
		{
			MethodName: "FacetAttackedNodes",
			Handler:    _SbotService_FacetAttackedNodes_Handler,
		},
		{
			MethodName: "CountAttackedNodes",
			Handler:    _SbotService_CountAttackedNodes_Handler,
		},
		{
			MethodName: "QueryAttackProcess",
			Handler:    _SbotService_QueryAttackProcess_Handler,
		},
		{
			MethodName: "FacetAttackProcess",
			Handler:    _SbotService_FacetAttackProcess_Handler,
		},
		{
			MethodName: "CountAttackProcess",
			Handler:    _SbotService_CountAttackProcess_Handler,
		},
		{
			MethodName: "QueryAttackedDownloadFiles",
			Handler:    _SbotService_QueryAttackedDownloadFiles_Handler,
		},
		{
			MethodName: "FacetAttackedDownloadFiles",
			Handler:    _SbotService_FacetAttackedDownloadFiles_Handler,
		},
		{
			MethodName: "CountAttackedDownloadFiles",
			Handler:    _SbotService_CountAttackedDownloadFiles_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service/sbot_service.proto",
}