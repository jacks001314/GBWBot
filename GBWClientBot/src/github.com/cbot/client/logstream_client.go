package client

import (
	"context"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"github.com/cbot/logstream"
	"github.com/cbot/node"
	"google.golang.org/grpc"
)

type LogStreamClient struct {
	nd *node.Node

	grpcClient *grpc.ClientConn

	logSub *logstream.LogSub

	logClient service.LogStreamServiceClient

	logStream *logstream.LogStream
	closeCh   chan int
}

func NewLogStreamClient(nd *node.Node, grpcClient *grpc.ClientConn, logStream *logstream.LogStream) *LogStreamClient {

	return &LogStreamClient{
		nd:         nd,
		grpcClient: grpcClient,
		logSub:     logStream.Sub("grpclog"),
		logClient:  service.NewLogStreamServiceClient(grpcClient),
		logStream:  logStream,
		closeCh:    make(chan int),
	}

}

func (ls *LogStreamClient) sendLog() error {

	logClient := ls.logClient

	logChannel, err := logClient.Channel(context.Background())

	if err != nil {

		ls.logStream.UnSub(ls.logSub)
		return err
	}

	err = logChannel.Send(&model.LogStream{
		NodeId: ls.nd.NodeId(),
		Dsize:  0,
		Data:   []byte{},
	})

	if err != nil {
		ls.logStream.UnSub(ls.logSub)
		return err
	}

	for {

		select {

		case <-ls.closeCh:
			ls.logStream.UnSub(ls.logSub)
			return nil

		case log := <-ls.logSub.LogChan:

			err = logChannel.Send(&model.LogStream{
				NodeId: ls.nd.NodeId(),
				Dsize:  uint64(len(log)),
				Data:   log,
			})

			if err != nil {
				ls.logStream.UnSub(ls.logSub)
				return err
			}

		}

	}

}

func (ls *LogStreamClient) Start() error {

	logClient := ls.logClient

	cmdMethod, err := logClient.CmdChannel(context.Background(), &model.TargetNode{NodeId: ls.nd.NodeId()})

	if err != nil {

		return err
	}

	for {

		cmd, err := cmdMethod.Recv()

		if err != nil || cmd.NodeId != ls.nd.NodeId() {

			continue
		}

		switch cmd.Op {

		case model.CmdOP_OPEN:

			go ls.sendLog()

		case model.CmdOP_CLOSE:

			ls.closeCh <- 1
		}

	}
}

func (ls *LogStreamClient) Stop() {

}
