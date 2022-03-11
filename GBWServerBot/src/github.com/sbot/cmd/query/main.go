package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"github.com/sbot/utils/jsonutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main(){

	sbotHost := flag.String("sbotHost","127.0.0.1","set the sbot rpc host")
	sbotPort := flag.Int("sbotPort",3333,"set the sbot rpc port")

	name := flag.String("name","","set the query name,valid name is :[task,node,process,download]")

	startTime := flag.Uint64("startTime",0,"set the start time to query")
	endTime := flag.Uint64("endTime",0,"set the end time to query")
	page := flag.Uint64("page",1,"set the page to query")
	pageSize := flag.Uint64("size",1,"set the page size to query")
	taskName := flag.String("taskName","","set the task name")
	taskId := flag.String("taskId","","set the taskId")
	nodeId := flag.String("nodeId","","set the nodeId")
	pnodeId := flag.String("pnodeId","","set the parent node id")
	mac := flag.String("mac","","set the node mac ")
	attackType := flag.String("attackType","","set the attack type")

	flag.Parse()

	if *name == "" {

		log.Fatalf("Please specify the name to query object,valid name is :[task,node,process,download]")

	}


	conn, err := grpc.Dial(fmt.Sprintf("%s:%d",*sbotHost,*sbotPort),grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Cannot connect to rhost:%s:%d,failed:%v", *sbotHost,*sbotPort,err)

	}

	defer conn.Close()

	client := service.NewSbotServiceClient(conn)

	switch *name {

	case "task":
		queryAttackTask(client,*startTime,*endTime,*page,*pageSize,*taskId,*taskName)

	case "node":
		queryAttackedNode(client,*startTime,*endTime,*page,*pageSize,*taskId,*nodeId,*pnodeId,*mac,*attackType)

	case "process":
		queryAttackProcess(client,*startTime,*endTime,*page,*pageSize,*taskId,*nodeId,*attackType)

	case "download":
		queryAttackedNodeDownloadFiles(client,*startTime,*endTime,*page,*pageSize,*taskId,*nodeId,*attackType)

	default:
		log.Fatalf("Unknown query object:%s\n",*name)

	}
}

func queryAttackedNodeDownloadFiles(client service.SbotServiceClient, startTime,endTime,page,size uint64, taskId,nodeId, attackType string) {

	query := &model.AttackedNodeDownloadFileQuery{
		TaskId:     taskId,
		NodeId:     nodeId,
		AttackType: attackType,
		Time:   &model.TimeRange{
			Start: startTime,
			End:   endTime,
		},
		Page:   &model.PageQuery{
			Page: page,
			Size: size,
		},
	}

	reply,err := client.QueryAttackedDownloadFiles(context.Background(),query)

	if err!=nil {
		log.Fatalf("Query AttackedNodeDownloaFiles Failed:%v\n",err)
	}

	log.Printf("Query AttackedNodeDownloaFiles Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.DownloadFiles {

		log.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}

}

func queryAttackProcess(client service.SbotServiceClient, startTime,endTime,page,size uint64, taskId,nodeId, attackType string) {

	query := &model.AttackProcessQuery{
		TaskId:     taskId,
		NodeId:     nodeId,
		AttackType: attackType,
		Time:   &model.TimeRange{
			Start: startTime,
			End:   endTime,
		},
		Page:   &model.PageQuery{
			Page: page,
			Size: size,
		},
	}

	reply,err:= client.QueryAttackProcess(context.Background(),query)

	if err!=nil {
		log.Fatalf("Query AttackedProcess Failed:%v\n",err)
	}

	log.Printf("Query AttackProcess Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.Aps {

		log.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}
}

func queryAttackedNode(client service.SbotServiceClient, startTime,endTime,page,size uint64, taskId,nodeId,pnodeId,mac,attackType string,) {

	query := &model.AttackedNodeQuery{
		TaskId:       taskId,
		ParentNodeId: pnodeId,
		NodeId:       nodeId,
		Mac:          mac,
		AttackType: attackType,
		Time:   &model.TimeRange{
			Start: startTime,
			End:   endTime,
		},
		Page:   &model.PageQuery{
			Page: page,
			Size: size,
		},
	}

	reply,err := client.QueryAttackedNodes(context.Background(),query)

	if err!=nil {
		log.Fatalf("Query AttackedNode Failed:%v\n",err)
	}

	log.Printf("Query Attacked Nodes Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.Nodes {

		log.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}
}


func queryAttackTask(client service.SbotServiceClient,startTime,endTime,page,size uint64,taskId,taskName string) {

	query := &model.AttackTaskQuery{
		TaskId: taskId,
		Name:   taskName,
		Time:   &model.TimeRange{
			Start: startTime,
			End:   endTime,
		},
		Page:   &model.PageQuery{
			Page: page,
			Size: size,
		},
	}

	reply,err:=client.QueryAttackTasks(context.Background(),query)

	if err!=nil {
		log.Fatalf("Query AttackTasks Failed:%v\n",err)
	}

	log.Printf("Query Attack Tasks Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.Messages {

		log.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}
}
