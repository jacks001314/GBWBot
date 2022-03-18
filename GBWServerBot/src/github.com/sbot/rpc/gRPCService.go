package rpc

import (
	"fmt"
	"github.com/sbot/handler"
	"github.com/sbot/proto/service"
	"github.com/sbot/rpc/rservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
)

type GRPCService struct {

	cfg *Config

	listener net.Listener

	grpcServer *grpc.Server

	attackTaskHandle *handler.AttackTaskHandler

	nodeHandle *handler.NodeHandler

	attackJarPayloadHandler *handler.AttackJarPayloadHandle

	sbotQueryHandle *handler.SbotQueryHandler

	attackTargetsHandle *handler.AttackTargetsHandler

	attackScriptsHandle *handler.AttackScriptsHandler

}

type Config struct {
	Host     string
	Port     int
	CertFlag string
	KeyFlag  string

	FDir string
}

func NewGRPCService(cfg *Config,
	attackTaskHandle *handler.AttackTaskHandler,
	nodeHandler *handler.NodeHandler,
	attackJarPayloadHandler *handler.AttackJarPayloadHandle,
	sbotQueryHandle *handler.SbotQueryHandler,
	attackTargetsHandle *handler.AttackTargetsHandler,
	attackScriptsHandle *handler.AttackScriptsHandler) *GRPCService {

	return &GRPCService{
		cfg:              cfg,
		attackTaskHandle: attackTaskHandle,
		nodeHandle:       nodeHandler,
		attackJarPayloadHandler: attackJarPayloadHandler,
		sbotQueryHandle: sbotQueryHandle,
		attackTargetsHandle:attackTargetsHandle,
		attackScriptsHandle:attackScriptsHandle,
	}
}

func (s *GRPCService) Start() {

	address := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {

		log.Errorf("Could not listen to port in Start() %s: %v", address, err)
	}

	s.listener = lis
	log.WithField("address", address).Info("gRPC server listening on port")

	opts := []grpc.ServerOption{}

	if s.cfg.CertFlag != "" && s.cfg.KeyFlag != "" {
		creds, err := credentials.NewServerTLSFromFile(s.cfg.CertFlag, s.cfg.KeyFlag)
		if err != nil {
			log.WithError(err).Fatal("Could not load TLS keys")
		}
		opts = append(opts, grpc.Creds(creds))
	} else {
		log.Warn("You are using an insecure gRPC server. If you are running your beacon node and " +
			"validator on the same machines, you can ignore this message. If you want to know " +
			"how to enable secure connections, see: https://docs.prylabs.network/docs/prysm-usage/secure-grpc")
	}

	s.grpcServer = grpc.NewServer(opts...)

	service.RegisterFileSerivceServer(s.grpcServer, rservice.NewFileService(s.cfg.FDir))
	service.RegisterNodeServiceServer(s.grpcServer, rservice.NewNodeService(s.nodeHandle))
	service.RegisterCmdServiceServer(s.grpcServer, rservice.NewCmdService())
	service.RegisterLogStreamServiceServer(s.grpcServer, rservice.NewLogStreamService())
	service.RegisterAttackTaskServiceServer(s.grpcServer, rservice.NewAttackTaskService(s.attackTaskHandle))
	service.RegisterAttackPayloadServiceServer(s.grpcServer,rservice.NewAttackPayloadService(s.attackJarPayloadHandler))
	service.RegisterSbotServiceServer(s.grpcServer,rservice.NewSbotQueryservice(s.sbotQueryHandle))
	service.RegisterAttackTargetsServiceServer(s.grpcServer,rservice.NewAttackTargetsService(s.attackTargetsHandle))
	service.RegisterAttackScriptsServiceServer(s.grpcServer,rservice.NewAttackScriptsService(s.attackScriptsHandle))

	// Register reflection service on gRPC server.
	reflection.Register(s.grpcServer)

	go func() {
		if s.listener != nil {
			if err := s.grpcServer.Serve(s.listener); err != nil {
				log.Errorf("Could not serve gRPC: %v", err)
			}
		}
	}()

}

func (s *GRPCService) Stop() {

}

func (s *GRPCService) Name() string {

	return "GRPCService"
}
