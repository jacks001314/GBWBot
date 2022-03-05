package node

import (
	"context"
	"fmt"
	"github.com/cbot/client/model"
	"github.com/cbot/client/service"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
)

type FServerClient struct {

	nd *Node

	grpcClient *grpc.ClientConn

	downloadStorePath string

	fsClient service.FileSerivceClient
}

func NewFServerClient(nd *Node,grpcClient *grpc.ClientConn,downloadStorePath string) *FServerClient {


	return &FServerClient{
		nd:                nd,
		grpcClient:        grpcClient,
		downloadStorePath: downloadStorePath,
		fsClient: service.NewFileSerivceClient(grpcClient),
	}

}

func (fsc *FServerClient)Download(fname string) (string,error) {

	fpath := filepath.Join(fsc.downloadStorePath,filepath.Base(fname))

	file,err := os.Create(fpath)

	if err!=nil {

		errS := fmt.Sprintf("Cannot open file:%s to store download file:%s",fpath,fname)

		log.Println(errS)

		return "",fmt.Errorf(errS)
	}

	defer file.Close()

	downloadRequest := &model.DownloadRequest{
		Node:  &model.Node{
			Status: 0,
			Id:    fsc.nd.nodeId,
		},

		Fname: fname,
	}

	dstream,err := fsc.fsClient.Download(context.Background(),downloadRequest)

	if err!=nil {

		errS := fmt.Sprintf("Download file:%s failed:%v",fpath,err)

		log.Println(errS)

		return "",fmt.Errorf(errS)

	}

	for {

		fpart,err := dstream.Recv()

		if err!=nil {

			errS := fmt.Sprintf("Receive file part failed:%v for file:%s",err,fpath)
			log.Println(errS)

			return "",fmt.Errorf(errS)

		}

		if fpart.IsLastParts {

			break
		}

		file.Write(fpart.Contents)
	}

	log.Printf("Download fname:%s to local file:%s ok",fname,fpath)

	return fpath,nil
}




