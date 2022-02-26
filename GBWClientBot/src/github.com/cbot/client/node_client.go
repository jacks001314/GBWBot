package client

import (
	"context"
	"fmt"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"github.com/cbot/node"
	"google.golang.org/grpc"
	"time"
)

type NodeClient struct {
	nd *node.Node

	nodeId string

	grpcClient *grpc.ClientConn

	nodeClient service.NodeServiceClient
}

func NewNodeClient(nd *node.Node, grpcClient *grpc.ClientConn) *NodeClient {

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

	ndreply, err := n.nodeClient.CreateNode(context.Background(), &model.CreateNodeRequest{
		Version:  "cbot-1.0",
		LocalIP:  nodeInfo.IP,
		OutIP:    nodeInfo.OutIP,
		Mac:      nodeInfo.Mac,
		Os:       nodeInfo.OS,
		Arch:     nodeInfo.Arch,
		User:     nodeInfo.User,
		HostName: nodeInfo.Hostname,
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
		NodeId: n.nodeId,
		Time:   uint64(time.Now().UnixNano() / (1000 * 1000)),
	})

	if err != nil {

		return err
	}

	return nil
}
