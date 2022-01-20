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


	r, err := c.FetchCmd(context.Background())

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	err = r.Send(&model.CmdReply{
		NodeId:   ip,
		Status:   0,
		Time:     0,
		Contents: []byte{},
	})

	if err!= nil {

		fmt.Printf("%v",err)
		return
	}

	for {

		// receive a cmd
		cmd,err:= r.Recv()

		if err!=nil {
			fmt.Printf("%v",err)
			return
		}

		fmt.Printf("receive a cmd:%v",cmd)

		err = r.Send(&model.CmdReply{
			NodeId:   ip,
			Status:   0,
			Time:     0,
			Contents: []byte(cmd.Name),
		})

		if err != nil {

			fmt.Printf("%v",err)
			return
		}

	}
}

func testChan(ich chan int,id string ){

	go func (){

		for {
			select {
			case i := <-ich:

				fmt.Printf("accept:%d,node:%s\n",i,id)
			}
		}
	}()

}


func main(){

	testRPCClient(os.Args[1],os.Args[2])

	/*
	ich := make(chan int )

	testChan(ich,"node1")
	testChan(ich,"node2")

	i := 0
	for {

		ich <- i

		time.Sleep(10*time.Second)
		i++
	}*/
	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//wg.Wait()
}


