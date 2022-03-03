package rservice

import (
	"context"
	"fmt"
	"github.com/sbot/handler"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"time"
)

type NodeService struct {
	service.UnimplementedNodeServiceServer

	handle *handler.NodeHandler
}

func NewNodeService(handle *handler.NodeHandler) *NodeService {

	return &NodeService{
		handle: handle,
	}
}

func (s *NodeService) CreateNode(ctx context.Context, request *model.CreateNodeRequest) (*model.Node, error) {

	nodeId, err := s.handle.HandleCreateNode(request)

	if err != nil {

		return &model.Node{
			Status: -1,
			Id:     "",
		}, err
	}

	return &model.Node{
		Status: 0,
		Id:     nodeId,
	}, nil
}

func (s *NodeService) Ping(ctx context.Context, request *model.PingRequest) (*model.PingReply, error) {

	s.handle.HandlePing(request)

	return &model.PingReply{
		Status:  0,
		Message: "ok",
		Time:    uint64(time.Now().UnixNano() / (1000 * 1000)),
	}, nil

}

func (s *NodeService) SendAttackProcessRequest(ctx context.Context, request *model.AttackProcessRequest) (*model.AttackProcessReply, error) {

	if err := s.handle.HandleAttackProcess(request); err != nil {

		return &model.AttackProcessReply{
			Status:  -1,
			Message: fmt.Sprintf("%v", err),
		}, err
	}

	return &model.AttackProcessReply{
		Status:  0,
		Message: "ok",
	}, nil
}
