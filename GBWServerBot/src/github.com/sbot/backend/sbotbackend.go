package backend

import (
	"github.com/sbot/rpc"
	"github.com/sbot/server"
	"github.com/sbot/store"
	"github.com/sbot/utils/jsonutils"
)

type SbotBackend struct {
	cfg *Config

	db store.Store

	rpcService *rpc.GRPCService

	dnslog *server.DNSServer

	fserver *server.FileServer
}

func NewSbotBacked(cfile string) (*SbotBackend, error) {

	var cfg Config

	if err := jsonutils.UNMarshalFromFile(&cfg, cfile); err != nil {

		log.Errorf("load config from file:%s is failed", cfile)
		return nil, err
	}

	return &SbotBackend{
		cfg:        &cfg,
		db:         nil,
		rpcService: nil,
		dnslog:     nil,
		fserver:    nil,
	}, nil
}

func (sb *SbotBackend) Start() error {

	return nil
}
