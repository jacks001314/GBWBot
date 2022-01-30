package rservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileService struct {

	// the dir that files store
	fdir string

	service.UnimplementedFileSerivceServer

}

type DirEntry struct {

	Name string `json:name`
	IsFile bool `json:"isFile"`
	Size  int64 `json:size`

}



func NewFileService(fdir string) *FileService {

	return &FileService{
		fdir:                           fdir ,
		UnimplementedFileSerivceServer: service.UnimplementedFileSerivceServer{},
	}

}

func (s *FileService) getPath(fpath string) string {


	if strings.Contains(fpath,"/") && (string(os.PathSeparator)!="/") {

		fpath = strings.ReplaceAll(fpath,"/",string(os.PathSeparator))

	}

	return filepath.Join(s.fdir,fpath)
}


func (s *FileService) Download(req *model.DownloadRequest, part service.FileSerivce_DownloadServer) error {

	var (
		buf     []byte
		n       int
		file    *os.File
	)

	fpath := s.getPath(req.Fname)

	file, err := os.Open(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			err = part.Send(&model.FilePart{
				Fpath:       fpath,
				Tbytes:      0,
				Bytes:       0,
				IsLastParts: true,
				Md5:         "",
				Contents:    []byte{},
			})
		}

		return err
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	buf = make([]byte, 1<<12)

	for  {

		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {

				err = part.Send(&model.FilePart{
					Fpath:       fpath,
					Tbytes:      size,
					Bytes:       0,
					IsLastParts: true,
					Md5:         "",
					Contents:   []byte{},
				})

				err = nil
			}

			return err
		}

		err = part.Send(&model.FilePart{
			Fpath:       fpath,
			Tbytes:      size,
			Bytes:       int64(n),
			IsLastParts: false,
			Md5:         "",
			Contents:   buf[:n],
		})

		if err != nil {
			return err
		}

	}
}


func (s *FileService) UPload(stream service.FileSerivce_UPloadServer) error {

	var (
		fpath string
		firstPart bool = true
		fd *os.File
		err error
	)

	for {
		fpart, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}

			return err
		}

		fpath =  s.getPath(fpart.Fpath)

		if fpart.Bytes != int64(len(fpart.Contents)) {
			return fmt.Errorf("%v == nk.SizeInBytes != int64(len(nk.Data)) == %v",fpart.Bytes,len(fpart.Contents))
		}

		if firstPart {
			if fpath != "" {
				fd, err = os.Create(fpath)
				if err != nil {
					return err
				}

				defer fd.Close()
			}

			firstPart = false
		}

		if fpart.IsLastParts {
			goto END
		}

		err = writeToFd(fd, fpart.Contents)
		if err != nil {
			return err
		}

	}

	END:
	err = stream.SendAndClose(&model.UPloadStatus{
		Status: 0,
		Fpath:  fpath,
	})

	return err
}

func writeToFd(fd *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {
		nw, err := fd.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}
}


func (s *FileService) listDir(args []string) (string,int) {

	if len(args)!= 1{

		return fmt.Sprintf(`{"message":"Invalid args number:%d"}`,len(args)),-1

	}

	entries := make([]*DirEntry,0)

	dir := s.getPath(args[0])

	files,err := ioutil.ReadDir(dir)

	if err!=nil {

		return fmt.Sprintf(`{"message":%v}`,err),-1
	}

	for _,f := range files {

		entries = append(entries,&DirEntry{
			Name:   f.Name(),
			IsFile: !f.IsDir(),
			Size:   f.Size(),
		})
	}


	data,_:= json.Marshal(entries)

	return fmt.Sprintf(`{"message:%s"}`,string(data)),0
}

func (s *FileService)mkdir(args []string) (string,int) {

	if len(args)!= 1{

		return fmt.Sprintf(`{"message":"Invalid args number:%d"}`,len(args)),-1
	}

	dir := s.getPath(args[0])

	if err := os.MkdirAll(dir,0755);err!=nil {

		return fmt.Sprintf(`{"message":%v}`,err),-1
	}

	return fmt.Sprintf(`{"message":"ok"}`),0
}

func (s *FileService)del(args []string) (string,int) {

	if len(args)!= 1{

		return fmt.Sprintf(`{"message":"Invalid args number:%d"}`,len(args)),-1
	}

	dir := s.getPath(args[0])

	if err := os.RemoveAll(dir);err!=nil {

		return fmt.Sprintf(`{"message":%v}`,err),-1
	}

	return fmt.Sprintf(`{"message":"ok"}`),0
}

func (s *FileService)rename(args []string) (string,int) {

	if len(args)!= 2 {

		return fmt.Sprintf(`{"message":"Invalid args number:%d"}`,len(args)),-1
	}

	oname := s.getPath(args[0])
	nname := s.getPath(args[1])

	if err := os.Rename(oname,nname);err!=nil {

		return fmt.Sprintf(`{"message":%v}`,err),-1
	}

	return fmt.Sprintf(`{"message":"ok"}`),0
}


func (s *FileService) FileCmd(ctx context.Context, req *model.FileCmdRequest) (*model.FileCmdResponse, error) {

	var message string
	var status int

	switch req.Cmd {

	case model.FileCmd_MKDIR:
		message,status = s.mkdir(req.Args)

	case model.FileCmd_LIST:
		message,status = s.listDir(req.Args)

	case model.FileCmd_DEL:
		message,status = s.del(req.Args)

	case model.FileCmd_RENAME:

		message,status = s.rename(req.Args)

	default:
		message = fmt.Sprintf(`{"message":"UnImplement File Cmd:%d"}`,req.Cmd)
		status = -1

	}

	return &model.FileCmdResponse{
		Status:   int32(status),
		Response: []byte(message),
	},nil

}
