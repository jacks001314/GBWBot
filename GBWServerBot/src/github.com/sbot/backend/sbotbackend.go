package backend

import (
	"fmt"
	"github.com/sbot/handler"
	"github.com/sbot/rpc"
	"github.com/sbot/server"
	"github.com/sbot/store"
	redisstore "github.com/sbot/store/redis"
	"github.com/sbot/utils/jsonutils"
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

func makeNodeHandle(cfg *Config) (*handler.NodeHandler, error) {

	dbnode, err := openRedisDB(cfg, AttackedNodeDB, AttackedNodeDBTable)
	if err != nil {

		return nil, err
	}

	attackProcessDB, err := openRedisDB(cfg, AttackedNodeDB, NodeAttackProcessTable)
	if err != nil {

		return nil, err
	}

	return handler.NewNodeHandler(dbnode, attackProcessDB), nil

}

func makeAttackTaskHandle(cfg *Config) (*handler.AttackTaskHandler, error) {

	db, err := openRedisDB(cfg, AttackTaskDB, AttackTaskTable)
	if err != nil {

		return nil, err
	}

	return handler.NewAttackTaskHandler(cfg.CBotFileStoreDir,
		cfg.AttackFileServerDir,
		cfg.RHost, cfg.RPort,
		cfg.AttackFileServerPort, db), nil
}

func makeAttackFileDownloadHandle(cfg *Config) (*handler.AttackFileServerHandle, error) {

	db, err := openRedisDB(cfg, AttackFileDownloadDB, AttackFileDownloadTable)
	if err != nil {

		return nil, err
	}

	return handler.NewAttackFileServerHandle(db), nil
}

func NewSbotBacked(cfile string) (*SbotBackend, error) {

	var cfg Config

	if err := jsonutils.UNMarshalFromFile(&cfg, cfile); err != nil {

		log.Errorf("load config from file:%s is failed", cfile)
		return nil, err
	}

	nodeHandle, err := makeNodeHandle(&cfg)
	if err != nil {

		log.Error("Create attacked node handler failed:%v", err)
		return nil, err
	}

	attackTaskHandle, err := makeAttackTaskHandle(&cfg)
	if err != nil {

		log.Error("Create attack task handler failed:%v", err)
		return nil, err
	}

	attackFileDownloadHandle, err := makeAttackFileDownloadHandle(&cfg)
	if err != nil {

		log.Error("Create attack file download handler failed:%v", err)
		return nil, err
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
		rpcService:       rpc.NewGRPCService(rpcCfg, attackTaskHandle, nodeHandle),
		attackFileServer: server.NewAttackFileServer(attackFileDownloadHandle, cfg.AttackFileServerDir, "0.0.0.0", cfg.AttackFileServerPort),
	}, nil
}

func (sb *SbotBackend) Start() {

	//start rpc
	sb.rpcService.Start()

	go sb.attackFileServer.Start()

}
