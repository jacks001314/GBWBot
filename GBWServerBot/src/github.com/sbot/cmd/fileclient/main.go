package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"path/filepath"
)


func upload(fclient service.FileSerivceClient,fpath string) {

	fname := filepath.Base(fpath)

	upstream,err := fclient.UPload(context.Background())

	if err!=nil {

		fmt.Printf("upload failed:%v",err)
		return
	}

	f,err:= os.Open(fpath)

	if err!=nil {
		log.Fatal(err)
	}
	defer f.Close()

	finfo,_:= f.Stat()

	reader := bufio.NewReader(f)

	buf := make([]byte,4096)

	for {

		n,err := reader.Read(buf)

		if err!=nil {

			if err != io.EOF {
				log.Fatal(err)
			}

			err = upstream.Send(&model.FilePart{
				Fpath:       "/test/"+fname,
				Tbytes:      finfo.Size(),
				Bytes:       0,
				IsLastParts: true,
				Md5:         "",
				Contents:    []byte{},
			})

			if err!=nil {
				log.Fatal(err)
			}

			break
		}

		err = upstream.Send(&model.FilePart{
			Fpath:       "/test/"+fname,
			Tbytes:      finfo.Size(),
			Bytes:       int64(n),
			IsLastParts: false,
			Md5:         "",
			Contents:    buf[:n],
		})

		if err!=nil {

			log.Fatal(err)
		}

	}

	status ,err :=upstream.CloseAndRecv()

	fmt.Printf("upload status:%v",status)

}

func download(fclient service.FileSerivceClient,fname string,storeDir string) {

	fpath := filepath.Join(storeDir,filepath.Base(fname))


	file,err := os.Create(fpath)

	if err!=nil {
		log.Fatal(err)
	}
	defer file.Close()

	downloadRequest := &model.DownloadRequest{
		Node:  &model.Node{
			Status: 0,
			Id:     "192.168.1.151",
		},

		Fname: fname,
	}

	dstream,err := fclient.Download(context.Background(),downloadRequest)

	if err!=nil {
		log.Fatal(err)
	}

	for {

		fpart,err := dstream.Recv()

		if err!=nil {

			log.Fatal(err)
		}

		if fpart.IsLastParts {

			break
		}

		file.Write(fpart.Contents)
	}

	fmt.Printf("Download fname:%s to local file:%s",fname,fpath)
}

func fileCmd (fclient service.FileSerivceClient,op string,args []string) {

	var cmdReq *model.FileCmdRequest

	switch  {

	case op == "mkdir":

		cmdReq = &model.FileCmdRequest{
			Cmd:  model.FileCmd_MKDIR,
			Args: args,
		}

	case op == "list":
		cmdReq = &model.FileCmdRequest{
			Cmd:  model.FileCmd_LIST,
			Args: args,
		}

	case op == "del":
		cmdReq = &model.FileCmdRequest{
			Cmd:  model.FileCmd_DEL,
			Args: args,
		}
	case op == "rename":
		cmdReq = &model.FileCmdRequest{
			Cmd:  model.FileCmd_RENAME,
			Args: args,
		}
	}

	cmdRes,err := fclient.FileCmd(context.Background(),cmdReq)

	if err!=nil {
		log.Fatal(err )
	}

	fmt.Printf("Run cmd:%s ok,the result is:%v\n",op,cmdRes)

}


func main() {



	addr := os.Args[1]
	op := os.Args[2]

	conn, err := grpc.Dial(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	fclient := service.NewFileSerivceClient(conn)


	switch  {

	case op == "upload":
		upload(fclient,os.Args[3])
		fmt.Println("upload ok------------------")

	case op == "download":

		download(fclient,os.Args[3],os.Args[4])
		fmt.Println("download ok -------------------")

	default:
		fileCmd(fclient,op,os.Args[3:])

	}
}


