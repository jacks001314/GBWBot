package node

import (
	"context"
	"fmt"
	"github.com/sbot/proto"
)

type Server struct {

	proto.UnimplementedNodeServer

}

func (s *Server) Ping(ctx context.Context, status *proto.Status) (*proto.PingReply, error) {

	fmt.Printf("{lip:%s,oip:%s,mac:%s,version:%s,os:%s}\n",status.LocalIP,status.OutIP,status.Mac,status.CbotVersion,status.Os)
	return &proto.PingReply{
		Status:  0,
		Message: "ok",
	},nil
}



