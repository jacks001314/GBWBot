package main

import (
	"context"
	"github.com/sbot/proto"
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
	c := proto.NewNodeClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	os,_:= os.Hostname()

	r, err := c.Ping(ctx, &proto.Status{
		LocalIP:     ip,
		OutIP:       ip,
		Mac:         "fff",
		Os:          os,
		CbotVersion: "1.1.0",
	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}

func main(){

	for {

		time.Sleep(10*time.Second)
		testRPCClient(os.Args[1],os.Args[2])

	}

}


