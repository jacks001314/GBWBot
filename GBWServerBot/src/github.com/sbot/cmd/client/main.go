package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"github.com/sbot/utils/fileutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func createAttackTask(conn *grpc.ClientConn, args []string) {

	if len(args)<7 {

		log.Fatalf("Invalid args for create attack task ,usage:createAttackTask args--->name:userId:host:port:user:passwd:privateKey")
	}

	client := service.NewAttackTaskServiceClient(conn)

	defer conn.Close()

	port,err:= strconv.ParseInt(args[3],10,32)

	if err!=nil {

		log.Fatalf("Cannot parse port:%s for create attack task args,failed:%v",args[3],err)

	}
	request := &model.CreateAttackTaskRequest{
		Name:       args[0],
		UserId:     args[1],
		Host:       args[2],
		Port:       int32(port),
		User:       args[4],
		Passwd:     args[5],
		PrivateKey: args[6],
		OsType:     model.OsType_Linux,
		Cbot:       "cbot",
	}

	if reply,err :=client.CreateAttackTask(context.Background(),request);err!=nil {

		log.Fatalf("create attack task failed:%v",err)

	}else {

		log.Printf("create attack task ok,reply.taskId:%s,reply.status:%d",reply.TaskId,reply.Status)

	}

}

func receiveCmdReply(cmdStream service.CmdService_RunCmdClient,ids int){

	var count int = 0

	timeOut := 2*time.Minute
	timeExp := time.Tick(timeOut)

	out:
		for {

		select {

		case <-timeExp:
			log.Fatal("receive cmd reply is timeout")

		default:

			if reply,err := cmdStream.Recv(); err!=nil {

				log.Fatalf("receive cmd reply failed:%v",err)

				goto out

			}else {

				rdata,_:= json.Marshal(reply)

				log.Printf("receive cmd reply,reply:%s",string(rdata))

				count++

				if ids >0 && count>=ids {
					log.Printf("Receive cmd is over")
					goto out
				}
			}
		}
	}

}

func addAttackScript(conn *grpc.ClientConn,args []string){

	var cmdName string
	var content []byte

	if len(args)<6 {

		log.Fatalf("invalid run add attack script cmd args number,args--->[nodeIds:scriptName:attackType:dport:dproto:fpath]")
	}

	client := service.NewCmdServiceClient(conn)

	nodeIds := strings.Split(args[0],",")
	scriptName := args[1]
	attackType := args[2]

	dport,err:= strconv.ParseInt(args[3],10,32)
	if err!=nil {

		log.Fatalf("Parse Default port for add attack script failed:%v,arg:%s",err,args[3])

	}

	dproto := args[4]
	fpath := args[5]

	if fileutils.FileIsExisted(fpath){

		cmdName = "notFromFile"
		if content,err= ioutil.ReadFile(fpath);err!=nil {

			log.Fatalf("Cannot read attack script content from file:%s,failed:%v",fpath,err)

		}

	}else {

		cmdName = "fromFile"
		content = []byte(fpath)
	}

	addRequest := model.AddAttackRequest{
		Name:         scriptName,
		AttackType:   attackType,
		DefaultPort:  int32(dport),
		DefaultProto: dproto,
		ContentLen:   uint64(len(content)),
		Content:      content,
	}

	addRequestJson,err:= json.Marshal(addRequest)

	if err!=nil {

		log.Fatalf("add attack script request to json failed:%v ",err)

	}

	addRequestBase64 := base64.StdEncoding.EncodeToString(addRequestJson)

	request := &model.CmdRequest{
		NodeIdS: nodeIds,
		Cmd:     &model.Cmd{
			Code: model.CmdCode_RunAddAttack,
			Name: cmdName,
			Args: []string{addRequestBase64},
		},
		Os:      "linux",
	}


	if cmdStream,err := client.RunCmd(context.Background(),request);err!=nil {

		log.Fatalf("run add attack script is failed:%v",err)

	}else {

		receiveCmdReply(cmdStream,len(nodeIds))
	}

}

func addAttackSourceScript(conn *grpc.ClientConn, args []string){

	var cmdName string
	var content []byte
	var err error

	if len(args)<4 {

		log.Fatalf("invalid run add attack source script cmd args number,args--->[nodeIds:scriptName:attackTypes:fpath]")
	}

	client := service.NewCmdServiceClient(conn)

	nodeIds := strings.Split(args[0],",")
	scriptName := args[1]
	attackTypes := strings.Split(args[2],",")

	fpath := args[3]

	if fileutils.FileIsExisted(fpath){

		cmdName = "notFromFile"
		if content,err= ioutil.ReadFile(fpath);err!=nil {

			log.Fatalf("Cannot read attack script source content from file:%s,failed:%v",fpath,err)

		}

	}else {

		cmdName = "fromFile"
		content = []byte(fpath)
	}

	addRequest := model.AddAttackSourceRequest{
		Name:       scriptName,
		Types:      attackTypes,
		ContentLen: uint64(len(content)),
		Content:    content,
	}

	addRequestJson,err:= json.Marshal(addRequest)

	if err!=nil {

		log.Fatalf("add attack script source request to json failed:%v ",err)

	}

	addRequestBase64 := base64.StdEncoding.EncodeToString(addRequestJson)

	request := &model.CmdRequest{
		NodeIdS: nodeIds,
		Cmd:     &model.Cmd{
			Code: model.CmdCode_RunAddAttackSource,
			Name: cmdName,
			Args: []string{addRequestBase64},
		},
		Os:      "linux",
	}


	if cmdStream,err := client.RunCmd(context.Background(),request);err!=nil {

		log.Fatalf("run add attack source script is failed:%v",err)

	}else {

		receiveCmdReply(cmdStream,len(nodeIds))
	}
}

func addBruteForceDict(conn *grpc.ClientConn, args []string){

	if len(args)<4 {

		log.Fatalf("invalid run add bruteforce dict cmd args number,args--->[nodeIds:name:users:passwdContent]")
	}

	client := service.NewCmdServiceClient(conn)

	nodeIds := strings.Split(args[0],",")
	name := args[1]
	users := args[2]
	passwds := args[3]

	request := &model.CmdRequest{
		NodeIdS: nodeIds,
		Cmd:     &model.Cmd{
			Code: model.CmdCode_RunAddDict,
			Name: name,
			Args: []string{users,passwds},
		},
		Os:      "linux",
	}

	if cmdStream,err := client.RunCmd(context.Background(),request);err!=nil {

		log.Fatalf("run add bruteforce dictory for name:%s,failed:%v",name,err)

	}else {

		receiveCmdReply(cmdStream,len(nodeIds))
	}

}

func runCmd(conn *grpc.ClientConn, args []string){

	if len(args)<3 {

		log.Fatalf("invalid run cmd args number,args--->[nodeIds:name:args]")
	}
	client := service.NewCmdServiceClient(conn)

	nodeIds := strings.Split(args[0],",")
	name := args[1]
	cmdArgs := strings.Split(args[2],",")

	request := &model.CmdRequest{
		NodeIdS: nodeIds,
		Cmd:     &model.Cmd{
			Code: model.CmdCode_RunOSCmd,
			Name: name,
			Args: cmdArgs,
		},
		Os:      "linux",
	}


	if cmdStream,err := client.RunCmd(context.Background(),request);err!=nil {

		log.Fatalf("run cmd:%s is failed:%v",name,err)

	}else {

		receiveCmdReply(cmdStream,len(nodeIds))
	}

}


func main(){

	rhost := flag.String("rhost","127.0.0.1","set the sbot rpc host")
	rport := flag.Int("rport",3333,"set the sbot rpc port")

	op := flag.String("op","","set the operator to run,op[createAttackTask|addAttackScript|addAttackSourceScript|addBruteForceDict|runCmd]")

	args := flag.String("args","","set the operator's args")

	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d",*rhost,*rport),grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Cannot connect to rhost:%s:%d,failed:%v", *rhost,*rport,err)

	}

	opargs := strings.Split(*args,":")

	switch *op {

	case "createAttackTask":

		createAttackTask(conn,opargs)

	case "addAttackScript":
		addAttackScript(conn,opargs)

	case "addAttackSourceScript":
		addAttackSourceScript(conn,opargs)

	case "addBruteForceDict":
		addBruteForceDict(conn,opargs)

	case "runCmd":
		runCmd(conn,opargs)

	default:
		log.Fatalf("Unknown operator:%s",*op)

	}

}


