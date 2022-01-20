package rservice

import (
	"fmt"
	"github.com/sbot/proto/model"
	"github.com/sbot/proto/service"
	"io"
	"os"
	"path/filepath"
)

type FileService struct {

	// the dir that files store
	fdir string

	service.UnimplementedFileSerivceServer

}

func NewFileService(fdir string) *FileService {

	return &FileService{
		fdir:                           fdir ,
		UnimplementedFileSerivceServer: service.UnimplementedFileSerivceServer{},
	}

}

func (s *FileService) Download(req *model.DownloadRequest, part service.FileSerivce_DownloadServer) error {

	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
	)

	fpath := filepath.Join(s.fdir,req.Fname)

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

	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}

			return err
		}

		err = part.Send(&model.FilePart{
			Fpath:       fpath,
			Tbytes:      size,
			Bytes:       int64(n),
			IsLastParts: writing == false,
			Md5:         "",
			Contents:   buf[:n],
		})

		if err != nil {
			return err
		}

	}


	return nil
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

		fpath =  filepath.Join(s.fdir,fpart.Fpath)

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

		err = writeToFd(fd, fpart.Contents)
		if err != nil {
			return err
		}
		if fpart.IsLastParts {
			goto END
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

