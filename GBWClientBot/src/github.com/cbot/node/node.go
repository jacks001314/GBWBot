package node

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/attack/ascript"
	"github.com/cbot/client"
	"github.com/cbot/targets/local"
	"github.com/cbot/targets/source"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Node struct {
	nodeId string

	cfg *Config

	grpcClient *grpc.ClientConn

	nodeClient *client.NodeClient

	cmdClient *client.CmdClient

	logStreamClient *client.LogStreamClient

	spool *source.SourcePool

	attackTasks *attack.AttackTasks

	nodeInfo *local.NodeInfo
}

func NewNode(cfg *Config) *Node {

	return &Node{
		cfg:         nil,
		spool:       nil,
		attackTasks: nil,
		nodeInfo:    local.GetNodeInfo(),
	}

}

func (n *Node) GetNodeInfo() *local.NodeInfo {

	return n.nodeInfo
}

func (n *Node) Start() error {

	var err error

	//connect to sbot
	n.grpcClient, err = grpc.Dial(fmt.Sprintf("%s:%d", n.cfg.sbotHost, n.cfg.sbotRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		return err
	}

	//connect to sbot,and create node
	n.nodeClient = client.NewNodeClient(n, n.grpcClient)
	n.nodeId, err = n.nodeClient.CreateNode()

	if err != nil {

		return err
	}

	return nil
}

func (n *Node) Stop() {

	n.grpcClient.Close()
}

func (n *Node) NodeId() string {

	return n.nodeId
}

func (n *Node) AddAttackSource(name string, types []string, content []byte) error {

	s, err := source.NewScriptSourceFromContent(n.spool, name, types, content)

	if err != nil {

		return err
	}

	n.spool.StartSource(s)

	return nil
}

func (n *Node) AddAttack(name string, attackType string, defaultProto string, defaultPort int, content []byte) error {

	att, err := ascript.NewAttackScriptFromContent(n.attackTasks, name, attackType, defaultPort, defaultProto, content)

	if err != nil {

		return err
	}

	n.attackTasks.AddAttack(att)

	return nil
}
