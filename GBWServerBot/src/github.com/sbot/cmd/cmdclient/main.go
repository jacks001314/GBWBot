package main

import (
	"context"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func testRPCClient(addr string,ip string){

	conn, err := grpc.Dial(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := service.NewCmdServiceClient(conn)

	r, err := c.RunCmd(context.Background(),&model.CmdRequest{
		NodeIdS: []string{"192.168.1.151","192.168.1.152"},
		Cmd:     &model.Cmd{
			Code: 0,
			Name: "ls",
			Args: []string{"/var/tmp"},
		},
		Os:      "linux",
	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	for {

		// receive a reply
		reply,err:= r.Recv()

		if err!=nil {
			fmt.Printf("%v\n",err)
			return
		}

		fmt.Printf("receive a reply:%v\n",reply)

	}
}


func testLogStream(addr string,ip string,cmd string ) {

	conn, err := grpc.Dial(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	logStream := service.NewLogStreamServiceClient(conn)

	if cmd == "open" {

		openStream,err := logStream.Open(context.Background(),&model.TargetNode{NodeId:ip})

		if err!= nil {

			fmt.Printf("%v",err )
			return
		}

		for {

			log,err:= openStream.Recv()

			if err!=nil {
				fmt.Printf("%v",err)
				return
			}

			fmt.Printf("Receive a log from node:%s,size:%d,content:%s\n",log.NodeId,log.Dsize,string(log.Data))

		}
	}

	if cmd == "close" {

		logStream.Close(context.Background(),&model.TargetNode{NodeId:ip})

	}


}


func main(){

	//testRPCClient(os.Args[1],os.Args[2])

	testLogStream(os.Args[1],os.Args[2],os.Args[3])

}


