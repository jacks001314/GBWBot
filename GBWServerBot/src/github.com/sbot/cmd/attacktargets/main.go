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

func main()  {

	rhost := flag.String("rhost","127.0.0.1","set the sbot rpc host")
	rport := flag.Int("rport",3333,"set the sbot rpc port")

	cfgId := flag.String("cfgId","","set the attack targets config id")

	name := flag.String("name","name","set the attack targets name")

	proto := flag.String("proto","","set attack targets protocol")

	app := flag.String("app","","set attack targets application")

	key := flag.String("key","","set attack targets api key for shodan/zoomeye.....")

	query := flag.String("query","","set attack targets api query for shodan/zoomeye.....")

	useDefaultPort := flag.Bool("udport",true,"set attack targets use Default port or not")

	nodeIds := flag.String("nodeIds","","set the nodeIds to accept this attack targets")

	attackTypes := flag.String("attackTypes","","set the attack types for this attack targets")

	whiteList := flag.String("wlist","","set the white list for iprange attack targets")

	blacklist :=  flag.String("blist","","set the black list for iprange attack targets")

	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d",*rhost,*rport),grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Cannot connect to rhost:%s:%d,failed:%v", *rhost,*rport,err)
	}

	client := service.NewAttackTargetsServiceClient(conn)

	reply,err:=client.AddAttackTargets(context.Background(),&model.AddAttackTargetsRequest{
		CfgId:          *cfgId,
		Name:           *name,
		Proto:          *proto,
		App:            *app,
		Key:            *key,
		Query:          *query,
		UseDefaultPort: *useDefaultPort,
		NodeIds:        strings.Split(*nodeIds,","),
		AttackTypes:    strings.Split(*attackTypes,","),
		WhiteLists:     strings.Split(*whiteList,","),
		BlackLists:     strings.Split(*blacklist,","),
	})

	if err!=nil {

		log.Fatalf("add attack targets failed:%v",err)
	}

	log.Printf("reply:%s\n",jsonutils.ToJsonString(reply,true))


}
