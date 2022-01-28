package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sbot/utils/fileutils"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var errURLFormat = errors.New("Invalid url format")


type FileServer struct {

	lock sync.Mutex

	cfg *Config


	fhandler http.Handler

	subs []*FileDownloadRequestSub
}

type FileDownloadRequest struct {

	Url string
	Fname string

	AttackType string
	AttackIP string
	TargetIP string
	TargetPort int

	TargetOutIP string
	DownloadTool string

	UserAgent string

}

type Config struct {

	//the root dir that file stores
	RootDir string `json:"rootDir"`

	//the bind host,default 0.0.0.0
	Host string 	`json:"host"`

	// the bind port
	Port int 		`json:"port"`


}

type FileDownloadRequestSub struct {

	requests chan *FileDownloadRequest
}

func (f *FileServer) NewFileDownloadRequestSub() *FileDownloadRequestSub {

	f.lock.Lock()
	defer f.lock.Unlock()

	sub := &FileDownloadRequestSub{requests:make(chan *FileDownloadRequest)}
	f.subs = append(f.subs,sub)

	return sub
}

func (f *FileServer) pub( fdr *FileDownloadRequest) {

	f.lock.Lock()
	defer f.lock.Unlock()

	for _,fsub := range f.subs {

		fsub.requests<-fdr
	}


}

func (s *FileDownloadRequestSub) Sub() chan *FileDownloadRequest {

	return s.requests
}

func getPath(fname string) string {

	dir := fname
	if ext:= filepath.Ext(fname);ext!="" {
		dir = fname[:len(fname)-len(ext)]
	}

	return fmt.Sprintf("/%s/%s",dir,fname)
}

func (f *FileServer) isFileExisted(fname string) bool {

	fpath := filepath.Join(f.cfg.RootDir,getPath(fname))

	return fileutils.FileIsExisted(fpath)
}

func (f *FileServer) makeHttpRequest(fname string,r *http.Request) *http.Request {

	path := getPath(fname)

	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = path

	return r2
}

func getArgsMap(content string) map[string]string {

	results := make(map[string]string)

	args := strings.Split(content,"&")

	for _,arg := range args {

		kv := strings.Split(arg,"=")

		if len(kv) <=1 {
			results[arg] = ""
		}else {

			results[kv[0]] = kv[1]
		}
	}

	return results
}

func (f *FileServer) makeFileDownloadRequst(r *http.Request) (*FileDownloadRequest,error) {

	path := r.URL.Path

	if path == "" ||path[0]!= '/' ||strings.LastIndex(path,"/")!=0 {

		return nil,errURLFormat
	}

	base64Content := path[1:]

	content,err:= base64.StdEncoding.DecodeString(base64Content)

	if err!=nil {

		return nil,err
	}

	argsMap := getArgsMap(string(content))

	port,err:= strconv.ParseInt(argsMap["tPort"],10,32)

	if err!=nil {

		return nil,errURLFormat
	}

	if !f.isFileExisted(argsMap["fname"]){

		return nil,errURLFormat
	}

	return &FileDownloadRequest{
		Url:          r.URL.Path,
		Fname:        argsMap["fname"],
		AttackType:   argsMap["atype"],
		AttackIP:     argsMap["pip"],
		TargetIP:     argsMap["tip"],
		TargetPort:   int(port),
		TargetOutIP:  r.RemoteAddr,
		DownloadTool: argsMap["dt"],
		UserAgent:    r.Header.Get("User-Agent"),
	},nil
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Infof("Receive a file download request,from:%s,url:%s",r.RemoteAddr,r.URL.Path)

	dreq,err:= f.makeFileDownloadRequst(r)

	if err!=nil {

		http.NotFound(w,r)
		return
	}

	newReq := f.makeHttpRequest(dreq.Fname,r)

	f.fhandler.ServeHTTP(w,newReq)

	f.pub(dreq)
}

func (f *FileServer) Start() {

	server := &http.Server{Addr:fmt.Sprintf("%s:%d",f.cfg.Host,f.cfg.Port),
		Handler: f}

	log.Infof("FileServer starting on directory: %s\nListening on http://%s:%s\n",f.cfg.RootDir,f.cfg.Host, f.cfg.Port)

	if err := server.ListenAndServe(); err != nil  {

		log.Errorf("FileServer ListenAndServe Failed:%v",err)

	}

}


func NewFileServer(cfg *Config) *FileServer {

	return &FileServer{
		lock:     sync.Mutex{},
		cfg:      cfg,
		fhandler: http.FileServer(http.Dir(cfg.RootDir)),
		subs:     make([]*FileDownloadRequestSub,0),
	}

}
