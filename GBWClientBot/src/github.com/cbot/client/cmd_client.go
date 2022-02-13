package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"github.com/cbot/node"
	"google.golang.org/grpc"
	"os/exec"
	"time"
)

type CmdClient struct {

	nd *node.Node

	grpcClient *grpc.ClientConn

	cmdClient  service.CmdService_FetchCmdClient

	stop bool

}

func NewCmdClient(nd *node.Node,grpcClient *grpc.ClientConn) (*CmdClient,error) {

	cmdClient,err := service.NewCmdServiceClient(grpcClient).FetchCmd(context.Background())

	if err!=nil {

		return nil,err
	}


	return &CmdClient{
		nd:         nd,
		grpcClient: grpcClient,
		cmdClient:  cmdClient,
		stop: false,
	},nil


}

func (cc *CmdClient) Start() error {


	err := cc.cmdClient.Send(&model.CmdReply{
		NodeId:   cc.nd.NodeId(),
		Status:   0,
		Time:     0,
		Contents: []byte{},
	})

	if err!= nil {

		return err
	}

	for {

		if cc.stop {
			break
		}

		// receive a cmd
		cmd,err:= cc.cmdClient.Recv()

		if err!=nil {
			continue
		}

		content,err := cc.handle(cmd)
		status := 0
		if err!=nil {

			status = -1
		}

		cc.cmdClient.Send(&model.CmdReply{
			NodeId:   cc.nd.NodeId(),
			Status:   int32(status),
			Time:     0,
			Contents: content,
		})

	}

	return nil
}

func (cc *CmdClient) Stop(){

	cc.stop = true

}


func (cc *CmdClient) handle(cmd *model.Cmd) ([]byte,error){

	switch cmd.Code {

	case model.CmdCode_RunAddAttackSource:

		err := cc.addAttackSource(cmd.Args[0])

		if err!=nil {

			return []byte(fmt.Sprintf("Add Attack Source is failed:%v",err)),err
		}

	case model.CmdCode_RunAddAttack:

		err := cc.addAttack(cmd.Args[0])

		if err!=nil {

			return []byte(fmt.Sprintf("Add Attack  is failed:%v",err)),err
		}

	case model.CmdCode_RunOSCmd:
		return cc.runOsCmd(cmd.Name,cmd.Args)

	}

	return []byte(fmt.Sprintf("Unkown cmd:%s",cmd.Name)),fmt.Errorf("Unkown cmd:%s",cmd.Name)
}


//add a attack source script

func (cc *CmdClient) addAttackSource(content string) error {

	var addSourceRequest model.AddAttackSourceRequest

	data,err := base64.StdEncoding.DecodeString(content)

	if err!=nil {

		return err
	}


	if err = json.Unmarshal(data,&addSourceRequest); err!=nil {

		return err
	}

	if addSourceRequest.ContentLen!=uint64(len(addSourceRequest.Content)) {

		return fmt.Errorf("Invalid attack source script content")
	}

	return cc.nd.AddAttackSource(addSourceRequest.Name,addSourceRequest.Types,addSourceRequest.Content)

}

//add a attack script
func (cc *CmdClient) addAttack(content string) error {

	var addAttackRequest model.AddAttackRequest

	data,err := base64.StdEncoding.DecodeString(content)

	if err!=nil {

		return err
	}


	if err = json.Unmarshal(data,&addAttackRequest); err!=nil {

		return err
	}

	if addAttackRequest.ContentLen!=uint64(len(addAttackRequest.Content)) {

		return fmt.Errorf("Invalid attack script content")
	}

	return cc.nd.AddAttack(addAttackRequest.Name,addAttackRequest.AttackType,addAttackRequest.DefaultProto,
		int(addAttackRequest.DefaultPort),addAttackRequest.Content)

}

//run os cmd
func (cc *CmdClient) runOsCmd(name string,args []string) ([]byte,error) {

	ctx, cancel := context.WithTimeout(context.Background(),30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name,args ...)

	res,err := cmd.CombinedOutput()

	if err!=nil {

		return []byte(fmt.Sprintf("%v",err)),err
	}

	return []byte(res),nil
}