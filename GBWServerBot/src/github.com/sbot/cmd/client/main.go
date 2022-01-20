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
	"time"
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

	for {

		select {
		case <-ich :
			fmt.Printf("accept a .....\n")

		default:
			fmt.Println("no ..............")
			time.Sleep(10*time.Second)

		}
	}
}

func openLogStream( logClient service.LogStreamServiceClient,ip string,closeCh chan model.CmdOP) error {

	logStream ,err :=logClient.Channel(context.Background())

	if err!=nil {

		return err
	}

	err = logStream.Send(&model.LogStream{
		NodeId: ip,
		Dsize:  0,
		Data:   []byte{},
	})

	if err !=nil {
		return err
	}

	for {

		select {
		case <-closeCh:

			fmt.Println("close log stream")
			//logStream.CloseSend()
			return nil

		default:

			logb := []byte(fmt.Sprintf("This a log at time:%d",time.Now().Second()))

			err = logStream.Send(&model.LogStream{
				NodeId: ip,
				Dsize:  uint64(len(logb)),
				Data:   logb,
			})

			if err!=nil {

				return err
			}

			time.Sleep(10*time.Second)
		}

	}
}

func testLogStream(addr string,ip string)  {


	ch := make(chan model.CmdOP)

	conn, err := grpc.Dial(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	logClient := service.NewLogStreamServiceClient(conn)

	cmdMethod,err := logClient.CmdChannel(context.Background(),&model.TargetNode{NodeId:ip})

	if err!=nil {

		fmt.Printf("%v",err)
		return
	}


	for {

		cmd,err:= cmdMethod.Recv()

		fmt.Printf("Reveive a cmd:%v from nodeID:%s\n",cmd,ip)

		if err!=nil {

			fmt.Printf("%v",err)
			return
		}

		if cmd.NodeId != ip {

			fmt.Printf("thie Cmd is not my ignore it!\n")
			continue
		}


		if cmd.Op == model.CmdOP_OPEN {

			go openLogStream(logClient, ip, ch)
		}

		if cmd.Op == model.CmdOP_CLOSE {

			ch<- model.CmdOP_CLOSE
		}
	}

}

func main(){

	//testRPCClient(os.Args[1],os.Args[2])

	testLogStream(os.Args[1],os.Args[2])

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


