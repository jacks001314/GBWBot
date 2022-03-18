package node

import (
	"context"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (

	FetchAttackScriptsTimeInterval = 10*time.Second
)

type AttackScriptsClient struct {

	nd *Node

	grpcClient *grpc.ClientConn

	attackScriptsClient service.AttackScriptsServiceClient

}

func NewAttackScriptsClient(nd *Node, grpcClient *grpc.ClientConn) *AttackScriptsClient {

	return &AttackScriptsClient{
		nd:                  nd,
		grpcClient:          grpcClient,
		attackScriptsClient: service.NewAttackScriptsServiceClient(grpcClient),
	}
}

func (asc *AttackScriptsClient) fetch(fstream service.AttackScriptsService_FetchAttackScriptsClient){

	for {
		scripts, err := fstream.Recv()

		if err == nil {

			if scripts.Content != nil && len(scripts.Content) > 0 {

				asc.nd.AddAttack(scripts.Name, scripts.AttackType, scripts.DefaultProto, int(scripts.DefaultPort), scripts.Content)
			}

			if !scripts.HasNext {
				//close this stream
				fstream.CloseSend()
				break
			}

		} else {
			log.Printf("fetch attack scripts failed:%v\n", err)
			fstream.CloseSend()
			break
		}
	}
}

func (asc *AttackScriptsClient)Start() {

	tchan := time.Tick(FetchAttackScriptsTimeInterval)

	go func() {

		for {

			select {
			case <-tchan:

				//try to fetch attack scripts from sbot

				fstream,err := asc.attackScriptsClient.FetchAttackScripts(context.Background(),&model.FetchAttackScriptsRequest{NodeId:asc.nd.nodeId})

				if err==nil {

					//ok
					asc.fetch(fstream)
				}
			}
		}
	}()

}