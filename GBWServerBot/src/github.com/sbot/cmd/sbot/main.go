package main

import (
	"flag"
	"github.com/sbot/backend"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	cfile := flag.String("cfile", "/etc/sbot.json", "set config file path")

	flag.Parse()

	bend, err := backend.NewSbotBacked(*cfile)

	if err != nil {

		log.Errorf("setup backend failed:%v", err)
		return
	}

	bend.Start()

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
