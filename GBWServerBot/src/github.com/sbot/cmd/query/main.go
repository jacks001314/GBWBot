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
	"strconv"
	"strings"
)

func main(){

	sbotHost := flag.String("sbotHost","127.0.0.1","set the sbot rpc host")
	sbotPort := flag.Int("sbotPort",3333,"set the sbot rpc port")

	name := flag.String("name","","set the query name,valid name is :[task,node,process,download,facet,count]")

	startTime := flag.Uint64("startTime",0,"set the start time to query")
	endTime := flag.Uint64("endTime",0,"set the end time to query")
	page := flag.Uint64("page",1,"set the page to query")
	pageSize := flag.Uint64("size",1,"set the page size to query")
	taskName := flag.String("taskName","","set the task name")
	taskId := flag.String("taskId","","set the taskId")
	nodeId := flag.String("nodeId","","set the nodeId")
	userId := flag.String("userId","","set the userId")
	pnodeId := flag.String("pnodeId","","set the parent node id")
	mac := flag.String("mac","","set the node mac ")
	nodeIP := flag.String("nodeIP","","set the node ip ")
	attackType := flag.String("attackType","","set the attack type")
	args := flag.String("args","","specify facet/count args")

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
		queryAttackTask(client,*startTime,*endTime,*page,*pageSize,*taskId,*taskName,*userId)

	case "node":
		queryAttackedNode(client,*startTime,*endTime,*page,*pageSize,*taskId,*nodeId,*pnodeId,*mac,*attackType,*nodeIP,*userId)

	case "process":
		queryAttackProcess(client,*startTime,*endTime,*page,*pageSize,*taskId,*nodeId,*attackType,*userId)

	case "download":
		queryAttackedNodeDownloadFiles(client,*startTime,*endTime,*page,*pageSize,*taskId,*nodeId,*attackType,*userId)

	case "facet":
		facet(client,*args,*userId)

	case "count":
		count(client,*args,*userId)

	default:
		log.Fatalf("Unknown query object:%s\n",*name)

	}
}

func facet(client service.SbotServiceClient,args string,userId string) {

	var reply *model.FacetReply
	var err error

	arr := strings.Split(args,":")
	if len(arr) != 4 {

		log.Fatalf("Please specify valid facet args:[name:term:topN:isDec]\n")
	}

	topN,_ := strconv.ParseInt(arr[2],10,32)
	isDec,_:= strconv.ParseBool(arr[3])

	request := &model.FacetRequest{
		UserId: userId,
		Term:  arr[1],
		TopN:  int32(topN),
		IsDec: isDec,
	}

	switch arr[0] {

	case "task":
		reply,err = client.FacetAttackTasks(context.Background(),request)

	case "node":
		reply,err = client.FacetAttackedNodes(context.Background(),request)

	case "process":
		reply,err = client.FacetAttackProcess(context.Background(),request)

	case "download":
		reply,err = client.FacetAttackedDownloadFiles(context.Background(),request)

	default:
		log.Fatalf("unknown facet dbname:%s",arr[0])
	}

	if err!=nil {
		log.Fatalf("facet dbname:%s failed:%v",arr[0],err)
	}

	fmt.Printf("facets:%s",jsonutils.ToJsonString(reply,true))


}

func count(client service.SbotServiceClient,name string,userId string)  {

	var count *model.Count
	var err error
	request := &model.CountRequest{UserId:userId}

	switch name {

	case "task":
		count,err = client.CountAttackTasks(context.Background(),request)

	case "node":
		count,err = client.CountAttackedNodes(context.Background(),request)

	case "process":
		count,err = client.CountAttackProcess(context.Background(),request)

	case "download":
		count,err = client.CountAttackedDownloadFiles(context.Background(),request)

	default:
		log.Fatalf("unknown count dbname:%s",name)
	}

	if err!=nil {
		log.Fatalf("count dbname:%s failed:%v",name,err)
	}

	fmt.Printf("%s.count:%d\n",name,count.C)

}

func queryAttackedNodeDownloadFiles(client service.SbotServiceClient, startTime,endTime,page,size uint64, taskId,nodeId, attackType string,userId string ) {

	query := &model.AttackedNodeDownloadFileQuery{
		UserId: userId,
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

	fmt.Printf("Query AttackedNodeDownloaFiles Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.DownloadFiles {

		fmt.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}

}

func queryAttackProcess(client service.SbotServiceClient, startTime,endTime,page,size uint64, taskId,nodeId, attackType string,userId string ) {

	query := &model.AttackProcessQuery{
		UserId: userId,
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

	fmt.Printf("Query AttackProcess Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.Aps {

		fmt.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}
}

func queryAttackedNode(client service.SbotServiceClient, startTime, endTime, page, size uint64, taskId, nodeId, pnodeId, mac, attackType, nodeIP,userId string) {

	query := &model.AttackedNodeQuery{
		UserId: userId,
		TaskId:       taskId,
		ParentNodeId: pnodeId,
		NodeId:       nodeId,
		NodeIP: 	  nodeIP,
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

	fmt.Printf("Query Attacked Nodes Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.Nodes {

		fmt.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}
}


func queryAttackTask(client service.SbotServiceClient,startTime,endTime,page,size uint64,taskId,taskName,userId string) {

	query := &model.AttackTaskQuery{
		TaskId: taskId,
		UserId: userId,
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

	fmt.Printf("Query Attack Tasks Results:%s\n",jsonutils.ToJsonString(reply.Page,true))

	for _,entry:= range reply.Messages {

		fmt.Printf("%s\n",jsonutils.ToJsonString(entry,true))

	}
}
