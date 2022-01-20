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

func main(){

	testRPCClient(os.Args[1],os.Args[2])

}


