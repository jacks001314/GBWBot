package main

import (
	"flag"
	"fmt"
	"github.com/sbot/backend"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfile := flag.String("cfile", "/etc/sbot.json", "set config file path")

	flag.Parse()

	bend, err := backend.NewSbotBacked(*cfile)

	if err != nil {

		fmt.Printf("setup backend failed:%v", err)
		return
	}

	bend.Start()

	waitExit()

}

func waitExit() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
