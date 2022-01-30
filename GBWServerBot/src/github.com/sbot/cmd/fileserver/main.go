package main

import (
	"github.com/sbot/rpc"
	"time"
)

func main() {


	cfg := &rpc.Config{
		Host:     "0.0.0.0",
		Port:     "8080",
		CertFlag: "",
		KeyFlag:  "",
	}

	s := rpc.NewGRPCService(cfg)

	s.Start()

	for {

		time.Sleep(10*time.Second)
	}
}
