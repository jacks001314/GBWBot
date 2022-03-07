package node

import (
	"context"
	"fmt"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"github.com/cbot/utils/fileutils"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
)

type AttackPayloadClient struct {

	nd *Node

	grpcClient *grpc.ClientConn

	attackPayloadClient service.AttackPayloadServiceClient
}


func NewAttackPayloadClient(nd *Node, grpcClient *grpc.ClientConn) *AttackPayloadClient {

	return &AttackPayloadClient{
		nd:         nd,
		grpcClient: grpcClient,
		attackPayloadClient: service.NewAttackPayloadServiceClient(grpcClient),
	}
}

func (apc *AttackPayloadClient) Start() error {

	return nil
}

func (apc *AttackPayloadClient) Stop() {

}

func (apc *AttackPayloadClient) MakeAttackJar(cmd string) (string,error){

	request := &model.MakeJarAttackPayloadRequest{
		TaskId: apc.nd.cfg.TaskId,
		NodeId: apc.nd.nodeId,
		Cmd:    cmd,
	}

	reply,err:=apc.attackPayloadClient.MakeJar(context.Background(),request)

	if err!=nil {

		errS := fmt.Sprintf("Make Attack Jar package from sbot failed:%v",err)

		log.Println(errS)
		return "",fmt.Errorf(errS)
	}

	jarFile := filepath.Join(os.TempDir(),"JarMain.jar")


	if err=fileutils.WriteFile(jarFile,reply.Content);err!=nil {

		errS := fmt.Sprintf("Write jar content into file:%s ,failed:%v",jarFile,err)
		log.Println(errS)
		return "",fmt.Errorf(errS)
	}

	return jarFile,nil
}