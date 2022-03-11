package node

import (
	"context"
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"github.com/cbot/utils/jsonutils"
	"google.golang.org/grpc"
	"log"
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

	request := &model.CreateNodeRequest{
		TaskId:   n.nd.TaskId(),
		PnodeId:  n.nd.ParentNodeId(),
		AttackType: n.nd.AttackType(),
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
	}

	log.Printf("Prepare to Send a create node request to sbot,details:%s",jsonutils.ToJsonString(request,true))

	ndreply, err := n.nodeClient.CreateNode(context.Background(), request)

	if err != nil {

		errS := fmt.Sprintf("Create node failed:%v",err)

		log.Println(errS)

		return "", fmt.Errorf(errS)
	}

	n.nodeId = ndreply.Id

	log.Printf("Create node ok,nodeId:%s,taskId:%s",n.nodeId,n.nd.TaskId())

	return n.nodeId, nil
}

func (n *NodeClient) Ping() error {

	request :=  &model.PingRequest{
		TaskId: n.nd.TaskId(),
		NodeId: n.nodeId,
		Time:   uint64(time.Now().UnixNano() / (1000 * 1000)),
	}

	log.Printf("Prepare to Send a Ping Message to sbot,details:%s",jsonutils.ToJsonString(request,true))

	reply, err := n.nodeClient.Ping(context.Background(),request)

	if err!=nil {

		errS := fmt.Sprintf("Send Ping Message to Sbot failed:%v",err)
		log.Println(errS)
		return fmt.Errorf(errS)
	}

	log.Printf("Send Ping Message to sbot ok,reply.status:%d,reply.message:%s",reply.Status,reply.Message)

	return nil
}

func (n *NodeClient) SendAttackProcess(process *attack.AttackProcess) error {

	request := &model.AttackProcessRequest{
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
	}

	log.Printf("Prepare to Send an attack process result to sbot,details:%s",jsonutils.ToJsonString(request,true))

	reply, err := n.nodeClient.SendAttackProcessRequest(context.Background(),request)

	if err!=nil {

		errS := fmt.Sprintf("Send Attack Process to sbot failed:%v",err)
		log.Println(errS)
		return fmt.Errorf(errS)
	}

	log.Printf("Send Attack Process to sbot is ok,reply.status:%d,reply.message:%s",reply.Status,reply.Message)

	return err
}
