package main

import (
	"flag"
	"github.com/cbot/node"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	pnodeId := flag.String("pnode", "", "set parent nodeId")
	sbotHost := flag.String("rhost", "127.0.0.1", "set sbot host")
	sbotRpcPort := flag.Int("rport", 3333, "set sbot rpc port")
	sbotFileServerPort := flag.Int("fport", 6666, "set sbot file server port")
	maxThreads := flag.Int("threads", 100, "set max threads that attack")
	scapacity := flag.Int("scap", 100, "set source queue capacity")
	acapacity := flag.Int("acap", 10, "set attack queue capacity")
	taskId := flag.String("taskId", "", "set the task id")

	flag.Parse()

	cfg := &node.Config{
		TaskId:                *taskId,
		PNodeId:               *pnodeId,
		SbotHost:              *sbotHost,
		SbotRPCPort:           *sbotRpcPort,
		SbotFileServerPort:    *sbotFileServerPort,
		MaxThreads:            *maxThreads,
		SourceCapacity:        *scapacity,
		AttackProcessCapacity: *acapacity,
	}

	node := node.NewNode(cfg)

	for {

		if err := node.Start(); err != nil {

			log.Printf("Start node failed:%v", err)
			time.Sleep(2 * time.Minute)
			continue
		} else {

			//start ok
			break
		}
	}

	waitExit()

}

func waitExit() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			select {
			case <-sigs:

				wg.Done()
			}
		}
	}()
}
