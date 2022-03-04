package node

import (
	"context"
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"google.golang.org/grpc"
	"time"
)

type NodeClient struct {
	nd *Node

	nodeId string

	grpcClient *grpc.ClientConn

	nodeClient service.NodeServiceClient
}

func NewNodeClient(nd *Node, grpcClient *grpc.ClientConn) *NodeClient {

	return &NodeClient{
		nd:         nd,
		nodeId:     "",
		grpcClient: grpcClient,
		nodeClient: service.NewNodeServiceClient(grpcClient),
	}
}

func (n *NodeClient) Start() error {

	return nil
}

func (n *NodeClient) Stop() {

}

func (n *NodeClient) CreateNode() (string, error) {

	nodeInfo := n.nd.GetNodeInfo()
	now := time.Now().UnixNano() / (1000 * 1000)

	ndreply, err := n.nodeClient.CreateNode(context.Background(), &model.CreateNodeRequest{
		TaskId:   n.nd.TaskId(),
		PnodeId:  n.nd.ParentNodeId(),
		Version:  "cbot-1.0",
		LocalIP:  nodeInfo.IP,
		OutIP:    nodeInfo.OutIP,
		Mac:      nodeInfo.Mac,
		Os:       nodeInfo.OS,
		Arch:     nodeInfo.Arch,
		User:     nodeInfo.User,
		HostName: nodeInfo.Hostname,
		Time:     uint64(now),
		LastTime: uint64(now),
	})

	if err != nil {

		return "", err
	}

	if ndreply.Status != 0 {

		return "", fmt.Errorf("create node failed:%d", ndreply.Status)
	}

	n.nodeId = ndreply.Id

	return n.nodeId, nil
}

func (n *NodeClient) Ping() error {

	_, err := n.nodeClient.Ping(context.Background(), &model.PingRequest{
		TaskId: n.nd.TaskId(),
		NodeId: n.nodeId,
		Time:   uint64(time.Now().UnixNano() / (1000 * 1000)),
	})

	return err
}

func (n *NodeClient) SendAttackProcess(process *attack.AttackProcess) error {

	_, err := n.nodeClient.SendAttackProcessRequest(context.Background(), &model.AttackProcessRequest{
		TaskId:     n.nd.TaskId(),
		NodeId:     n.nodeId,
		Time:       uint64(time.Now().UnixNano() / (1000 * 1000)),
		TargetIP:   process.IP,
		TargetHost: process.Host,
		TargetPort: int32(process.Port),
		Proto:      process.Proto,
		App:        process.App,
		Os:         process.OS,
		AttackName: process.Name,
		AttackType: process.Type,
		Status:     int32(process.Status),
		Payload:    process.Payload,
		Result:     process.Result,
		Details:    process.Details,
	})

	return err
}
