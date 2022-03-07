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

// AttackPayloadServiceClient is the client API for AttackPayloadService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AttackPayloadServiceClient interface {
	//generate a java jar attack payload
	MakeJar(ctx context.Context, in *model.MakeJarAttackPayloadRequest, opts ...grpc.CallOption) (*model.MakeJarAttackPayloadReply, error)
}

type attackPayloadServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAttackPayloadServiceClient(cc grpc.ClientConnInterface) AttackPayloadServiceClient {
	return &attackPayloadServiceClient{cc}
}

func (c *attackPayloadServiceClient) MakeJar(ctx context.Context, in *model.MakeJarAttackPayloadRequest, opts ...grpc.CallOption) (*model.MakeJarAttackPayloadReply, error) {
	out := new(model.MakeJarAttackPayloadReply)
	err := c.cc.Invoke(ctx, "/sbot.proto.service.AttackPayloadService/MakeJar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AttackPayloadServiceServer is the server API for AttackPayloadService service.
// All implementations must embed UnimplementedAttackPayloadServiceServer
// for forward compatibility
type AttackPayloadServiceServer interface {
	//generate a java jar attack payload
	MakeJar(context.Context, *model.MakeJarAttackPayloadRequest) (*model.MakeJarAttackPayloadReply, error)
	mustEmbedUnimplementedAttackPayloadServiceServer()
}

// UnimplementedAttackPayloadServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAttackPayloadServiceServer struct {
}

func (UnimplementedAttackPayloadServiceServer) MakeJar(context.Context, *model.MakeJarAttackPayloadRequest) (*model.MakeJarAttackPayloadReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeJar not implemented")
}
func (UnimplementedAttackPayloadServiceServer) mustEmbedUnimplementedAttackPayloadServiceServer() {}

// UnsafeAttackPayloadServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AttackPayloadServiceServer will
// result in compilation errors.
type UnsafeAttackPayloadServiceServer interface {
	mustEmbedUnimplementedAttackPayloadServiceServer()
}

func RegisterAttackPayloadServiceServer(s grpc.ServiceRegistrar, srv AttackPayloadServiceServer) {
	s.RegisterService(&AttackPayloadService_ServiceDesc, srv)
}

func _AttackPayloadService_MakeJar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(model.MakeJarAttackPayloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttackPayloadServiceServer).MakeJar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sbot.proto.service.AttackPayloadService/MakeJar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttackPayloadServiceServer).MakeJar(ctx, req.(*model.MakeJarAttackPayloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AttackPayloadService_ServiceDesc is the grpc.ServiceDesc for AttackPayloadService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AttackPayloadService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sbot.proto.service.AttackPayloadService",
	HandlerType: (*AttackPayloadServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MakeJar",
			Handler:    _AttackPayloadService_MakeJar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service/generate_attack_payload.proto",
}
