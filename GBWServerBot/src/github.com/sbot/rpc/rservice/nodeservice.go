package rservice

import (
	"context"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NodeService struct {

	service.UnimplementedNodeServiceServer
}

func (s *NodeService) CreateNode(context.Context, *model.CreateNodeRequest) (*model.Node, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNode not implemented")
}

func (s *NodeService) Ping(context.Context, *model.PingRequest) (*model.PingReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}


