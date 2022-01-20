package rservice

import (
	"context"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type LogStreamService struct {

	lock sync.Mutex

	nodes map[string]*nodeContext

	service.UnimplementedLogStreamServiceServer

}

type nodeContext struct {

	cmdCh chan *model.LogCmd
	logCh chan *model.LogStream
	notifyCh chan int
	isOpen bool
}

func NewLogStreamService() *LogStreamService {

	return &LogStreamService{
		lock:                                sync.Mutex{},
		nodes:                               make(map[string]*nodeContext),
		UnimplementedLogStreamServiceServer: service.UnimplementedLogStreamServiceServer{},
	}

}


func (s *LogStreamService) addNode(target *model.TargetNode) *nodeContext {

	s.lock.Lock()
	defer s.lock.Unlock()

	nodeCtx := &nodeContext{
		cmdCh:    make(chan *model.LogCmd),
		logCh:    make(chan *model.LogStream),
		notifyCh: make(chan int),
		isOpen:false,
	}

	s.nodes[target.NodeId] = nodeCtx

	return nodeCtx
}

func (s *LogStreamService) removeNode(target *model.TargetNode) {

	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.nodes,target.NodeId)

}

func (s *LogStreamService) findNode(target *model.TargetNode) *nodeContext {

	s.lock.Lock()
	defer s.lock.Unlock()


	if nodeCtx,ok:= s.nodes[target.NodeId];ok {

		return nodeCtx
	}

	return nil
}

func (s *LogStreamService) Open(target *model.TargetNode, stream service.LogStreamService_OpenServer) error {

	nodeCtx := s.findNode(target)

	if nodeCtx == nil {

		return status.Errorf(codes.NotFound,"Cannot find node:%s",target.NodeId)

	}

	if nodeCtx.isOpen {

		return status.Errorf(codes.AlreadyExists,"Node:%s has been openned,first close it")

	}

	//send a open cmd to channel
	logCmd := &model.LogCmd{
		Op:     model.CmdOP_OPEN,
		NodeId: target.NodeId,
	}

	nodeCtx.cmdCh<- logCmd

	//wait cbot call Channel Method to open stream
	<-nodeCtx.notifyCh

	nodeCtx.isOpen = true

	//can listen logstream
	for {
		select {

		case logStream := <- nodeCtx.logCh:

			log.WithField("NodeId",target.NodeId).Info("Receive a log ,logsize:%d",logStream.Dsize)

			//send to client

			if err := stream.Send(logStream);err!=nil {

				//s.removeNode(target)

				s.Close(context.Background(),target)

				log.WithField("NodeId",target.NodeId).Info("Cannot Send log to client,Remove it")

				return status.Errorf(codes.Unavailable,"Cannot Send log data to client")

			}
		}
	}
}

func (s *LogStreamService) Close(ctx context.Context, target *model.TargetNode) (*model.OPStatus, error) {

	nodeCtx := s.findNode(target)

	if nodeCtx == nil {

		return &model.OPStatus{
			Status:  -1,
			Message: "Cannot Find Node",
			NodeId:  target.NodeId,
		},status.Errorf(codes.NotFound,"Cannot find node:%s",target.NodeId)

	}

	if !nodeCtx.isOpen {

		return &model.OPStatus{
			Status:  -1,
			Message: "Node has been closed",
			NodeId:  target.NodeId,
		},status.Errorf(codes.Unavailable,"Node:%s has been closed")

	}

	//send a open cmd to channel
	logCmd := &model.LogCmd{
		Op:     model.CmdOP_CLOSE,
		NodeId: target.NodeId,
	}

	nodeCtx.cmdCh<- logCmd

	//wait cbot call Channel Method to close stream
	<-nodeCtx.notifyCh

	log.WithField("NodeId",target.NodeId).Info("Node log stream has been closed")

	//s.removeNode(target)

	nodeCtx.isOpen = false

	return &model.OPStatus{
		Status:  0,
		Message: "Ok",
		NodeId:  target.NodeId,
	},nil

}


func (s *LogStreamService) CmdChannel(target *model.TargetNode, logCmd service.LogStreamService_CmdChannelServer) error {

	//when cbot call this method ,add node to context
	nodeCtx := s.addNode(target)

	for {

		select {
		case cmd :=<- nodeCtx.cmdCh:

			if err := logCmd.Send(cmd);err!=nil {

				log.WithField("NodeId",target.NodeId).Info("Cannot Send Cmd to node")
				return status.Errorf(codes.Unavailable,"Cannot Send cmd to node:%s",target.NodeId)
			}

			//ok

			if cmd.Op == model.CmdOP_CLOSE{

				nodeCtx.notifyCh<- 0

			}

		}
	}
}

func (s *LogStreamService) Channel(stream service.LogStreamService_ChannelServer) error {

	//cbot first send a empty logstream
	flogStream,err:= stream.Recv()

	if err !=nil {
		return err
	}


	nodeCtx := s.findNode(&model.TargetNode{NodeId:flogStream.NodeId})

	if nodeCtx ==nil {

		return status.Errorf(codes.Canceled,"Cannot find a node context,maybe has been closed")

	}

	// notify open ,now open ok
	nodeCtx.notifyCh<- 0

	for {

		logstream,err := stream.Recv()

		if err!=nil {

			return err
		}

		//ok
		log.WithField("NodeId",flogStream.NodeId).Info("Receive a log from cbot node")

		nodeCtx.logCh <- logstream

	}
}
