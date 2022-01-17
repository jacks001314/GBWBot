package main

import (
	"github.com/sbot/rpc"
	"os"
	"time"
)

func testRPC(host string,port string)  {

	cfg := rpc.Config{
		Host:     host,
		Port:     port,
		CertFlag: "",
		KeyFlag:  "",
	}

	service := rpc.NewGRPCService(&cfg)

	service.Start()

	for {

		time.Sleep(10*time.Second)


	}

}



func main(){

	testRPC(os.Args[1],os.Args[2])



}

