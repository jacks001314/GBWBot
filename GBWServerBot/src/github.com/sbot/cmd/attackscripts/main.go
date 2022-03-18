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
	"strings"
)

func main(){


	rhost := flag.String("rhost","127.0.0.1","set the sbot rpc host")
	rport := flag.Int("rport",3333,"set the sbot rpc port")

	cfgIds := flag.String("cfgIds","","set the attack scripts config ids")

	nodeIds := flag.String("nodeIds","","set the nodeIds to accept this attack scripts")

	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d",*rhost,*rport),grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Cannot connect to rhost:%s:%d,failed:%v", *rhost,*rport,err)
	}

	client := service.NewAttackScriptsServiceClient(conn)

	reply,err:=client.AddAttackScripts(context.Background(),&model.AddAttackScriptsRequest{
		CfgIds:  strings.Split(*cfgIds,","),
		NodeIds: strings.Split(*nodeIds,","),
	})


	if err!=nil {

		log.Fatalf("add attack scripts failed:%v",err)
	}

	log.Printf("reply:%s\n",jsonutils.ToJsonString(reply,true))

}
