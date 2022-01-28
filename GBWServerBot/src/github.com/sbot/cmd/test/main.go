package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sbot/server"
	"os"
)


func getReq(fname string) string {

	/*
	Fname:        argsMap["fname"],
		AttackType:   argsMap["atype"],
			AttackIP:     argsMap["pip"],
			TargetIP:     argsMap["tip"],
			TargetPort:   int(port),
			TargetOutIP:  r.RemoteAddr,
			DownloadTool: argsMap["dt"],
	*/

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("fname=%s&atype=hadoop&pip=192.168.1.151&tip=192.168.1.152&tPort=8081&dt=wget",fname)))

}



func main(){


	cfg := &server.Config{
		RootDir: os.Args[1],
		Host:    "0.0.0.0",
		Port:    8080,
	}

	f := server.NewFileServer(cfg)

	go f.Start()

	sub := f.NewFileDownloadRequestSub()

	for {

		select {

		case fr := <- sub.Sub():

			data,_:= json.Marshal(fr)

			fmt.Println(string(data))

		}
	}

	fmt.Println(getReq("init.sh"))
}

