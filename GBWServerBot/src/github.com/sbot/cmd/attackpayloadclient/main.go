package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"github.com/sbot/utils/fileutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"path/filepath"
)


func  MakeAttackJar(client service.AttackPayloadServiceClient,cmd string) (string,error){

	request := &model.MakeJarAttackPayloadRequest{
		TaskId: "test",
		NodeId: "test",
		Cmd:    cmd,
	}

	reply,err:=client.MakeJar(context.Background(),request)

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

func main(){


	addr := flag.String("addr","127.0.0.1:3333","set the remote sbot address")
	cmd := flag.String("cmd","","set the attack payload the cmd to been run")

	flag.Parse()

	conn, err := grpc.Dial(*addr,grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Cannot connect to sbot:%s,failed:%v", addr,err)
	}

	defer conn.Close()

	client := service.NewAttackPayloadServiceClient(conn)

	file,err:=MakeAttackJar(client,*cmd)

	if err!=nil {
		log.Fatalf("%v",err)
	}

	log.Printf("Make Jar Attack Payload ok,file:%s",file)

}
