package rservice

import (
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type CmdService struct {

	lock sync.Mutex

	service.UnimplementedCmdServiceServer

	cbots map[string]*cbotContext

}

type cbotContext struct {

	cmdCh chan *cmdContext

}

type cmdContext struct {


	curCmd *model.CmdRequest
	replyCh chan *model.CmdReply

	acceptedNum int
}

func NewCmdService() *CmdService {

	return &CmdService{
		lock:                          sync.Mutex{},
		UnimplementedCmdServiceServer: service.UnimplementedCmdServiceServer{},
		cbots:                         make(map[string]*cbotContext),
	}

}

func (s *CmdService) makeCbotContext(nodeId string) *cbotContext{

	cbotCtx := &cbotContext{cmdCh:make(chan *cmdContext)}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.cbots[nodeId] = cbotCtx

	return cbotCtx
}

func (s *CmdService) removeCbotContext(nodeId string) {

	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.cbots,nodeId)

	log.WithField("NodeId",nodeId).Info("Remove cbot context ")

}

func canRunCmd(cmdReq *model.CmdRequest,nodeId string) bool {

	if cmdReq.NodeIdS == nil || len(cmdReq.NodeIdS) == 0 {

		//broadcast this cmd
		return true
	}

	for _,nid := range cmdReq.NodeIdS {

		if nid == nodeId {

			return true
		}
	}

	return false
}

func (s *CmdService) sendCmds(req *model.CmdRequest) *cmdContext{

	s.lock.Lock()
	defer s.lock.Unlock()

	cmdCtx := &cmdContext{
		curCmd:      req,
		replyCh: make(chan *model.CmdReply),
		acceptedNum: 0,
	}

	for nodeId,cbotCtx := range s.cbots {

		if canRunCmd(req,nodeId) {

			cbotCtx.cmdCh<-cmdCtx
			cmdCtx.acceptedNum ++
		}
	}

	return cmdCtx
}

func (s *CmdService) RunCmd(req *model.CmdRequest, reply service.CmdService_RunCmdServer) error {

	//write to cmd channel to wait cmd reply

	cmdCtx := s.sendCmds(req)

	if cmdCtx.acceptedNum == 0 {

		//no cbots can accept this cmd
		return status.Errorf(codes.Aborted,"no cbots can accept this cmd")

	}

	for {
		select {
		case res := <-cmdCtx.replyCh :

			log.WithField("NodeId",res.NodeId).Info("Accept a reply from cbot")

			cmdCtx.acceptedNum--

			if err := reply.Send(res); err!=nil {

				return status.Errorf(codes.Unavailable, "Could not send over stream: %v", err)
			}

			if cmdCtx.acceptedNum == 0 {

				return status.Errorf(codes.OK,"Reply Over.......................")
			}
		}
	}

	return nil
}

func (s *CmdService) FetchCmd(stream service.CmdService_FetchCmdServer) error {


	var nodeId string

	// the client --cbot should send a empty reply contains nodeID
	fisrtReply,err := stream.Recv()
	if err != nil {

		return err
	}

	nodeId = fisrtReply.NodeId

	cbotCtx := s.makeCbotContext(nodeId)

	log.WithField("nodeId",nodeId).Info("accept a cbot client")

	for {

		select {

		case cmdCtx := <-cbotCtx.cmdCh:

			cmd := cmdCtx.curCmd

			log.WithField("cmd", cmd.Cmd.Name).
					WithField("cmdArgs", cmd.Cmd.Args).
					WithField("nodeIds", cmd.NodeIdS).
					Info("Accept cmd")

			//send this cmd to client---cbot
			if err = stream.Send(cmd.Cmd); err != nil {

				s.removeCbotContext(nodeId)
				cmdCtx.replyCh <- &model.CmdReply{
					NodeId:   nodeId,
					Status:   0,
					Time:     0,
					Contents: []byte(fmt.Sprintf("%v",err)),
				}

				return status.Errorf(codes.Unavailable, "Could not send over stream: %v", err)
			}

			//wait to reply from cbot
			reply, err := stream.Recv()
			if err != nil {
				s.removeCbotContext(nodeId)
				cmdCtx.replyCh <- &model.CmdReply{
					NodeId:   nodeId,
					Status:   0,
					Time:     0,
					Contents: []byte(fmt.Sprintf("%v",err)),
				}

				return err
			}

			log.Infof("Receive a reply:%s from cbot:%s", string(reply.Contents), reply.NodeId)
			cmdCtx.replyCh <- reply

		}
	}
}

