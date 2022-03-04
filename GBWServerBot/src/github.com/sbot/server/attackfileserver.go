package server

import (
	"fmt"
	"github.com/sbot/handler"
	"github.com/sbot/utils/fileutils"
	"github.com/sbot/utils/netutils"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	AttackTaskPathKey = "/attack/tasks/"
)

type AttackFileServer struct {
	fhandler http.Handler

	attackFileHandle *handler.AttackFileServerHandle

	attackFileDir string

	host string

	port int
}

type AttackFileDownloadRequest struct {
	TaskId string `json:"taskId"`
	NodeId string `json:"nodeId"`

	Url   string `json:"url"`
	Fname string `json:"fname"`

	AttackType string `json:"attackType"`
	AttackIP   string `json:"attackIP"`
	TargetIP   string `json:"targetIP"`
	TargetPort int    `json:"targetPort"`

	TargetOutIP  string `json:"targetOutIP"`
	DownloadTool string `json:"downloadTool"`

	UserAgent string `json:"userAgent"`
}

func (afs *AttackFileServer) isFileExisted(req *AttackFileDownloadRequest) bool {

	fpath := filepath.Join(afs.attackFileDir, req.TaskId, req.Fname)

	return fileutils.FileIsExisted(fpath)
}

func (afs *AttackFileServer) makeHttpRequest(req *AttackFileDownloadRequest, r *http.Request) *http.Request {

	path := fmt.Sprintf("/%s/%s", req.TaskId, req.Fname)

	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = path

	return r2
}

func (afs *AttackFileServer) makeAttackFileDownloadRequest(path string, r *http.Request) (*AttackFileDownloadRequest, error) {

	paths := strings.Split(path, "/")

	if len(paths) < 5 {

		errS := fmt.Sprintf("Invalid url path:%s", path)
		log.WithField("urlPath", path).Error(errS)

		return nil, fmt.Errorf(errS)

	}

	taskId := paths[3]
	nodeId := paths[4]
	fname := paths[len(paths)-1]

	return &AttackFileDownloadRequest{
		TaskId:       taskId,
		NodeId:       nodeId,
		Url:          path,
		Fname:        fname,
		AttackType:   "attack",
		AttackIP:     "0.0.0.0",
		TargetIP:     r.RemoteAddr,
		TargetPort:   0,
		TargetOutIP:  r.RemoteAddr,
		DownloadTool: "wget",
		UserAgent:    r.Header.Get("User-Agent"),
	}, nil
}

func (afs *AttackFileServer) makeAttackFileDownloadRequstFromCryptPath(path string, r *http.Request) (*AttackFileDownloadRequest, error) {

	if path == "" || path[0] != '/' || strings.LastIndex(path, "/") != 0 {

		errS := fmt.Sprintf("Invalid url path:%s", path)
		log.WithField("urlPath", path).Error(errS)

		return nil, fmt.Errorf(errS)
	}

	ucr, err := netutils.DeCryptToURLPath(path[1:])

	if err != nil {
		errS := fmt.Sprintf("Invalid url path:%s", path)
		log.WithField("urlPath", path).Error(errS)

		return nil, fmt.Errorf(errS)
	}

	return &AttackFileDownloadRequest{
		TaskId:       ucr.TaskId,
		NodeId:       ucr.NodeId,
		Url:          r.URL.Path,
		Fname:        ucr.Fname,
		AttackType:   ucr.AttackType,
		AttackIP:     ucr.AttackIP,
		TargetIP:     ucr.TargetIP,
		TargetPort:   ucr.TargetPort,
		TargetOutIP:  r.RemoteAddr,
		DownloadTool: ucr.DownloadTool,
		UserAgent:    r.Header.Get("User-Agent"),
	}, nil
}

func (afs *AttackFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var dreq *AttackFileDownloadRequest
	var err error

	log.Infof("Receive a attack file download request,from:%s,url:%s", r.RemoteAddr, r.URL.Path)

	path := r.URL.Path

	if strings.HasPrefix(path, AttackTaskPathKey) {

		dreq, err = afs.makeAttackFileDownloadRequest(path, r)

	} else {

		dreq, err = afs.makeAttackFileDownloadRequstFromCryptPath(path, r)

	}

	if err != nil || !afs.isFileExisted(dreq) {

		log.WithField("path", filepath.Join(afs.attackFileDir, dreq.TaskId, dreq.Fname)).Error("Attack File is not Foud")

		http.NotFound(w, r)

		return
	}

	afs.fhandler.ServeHTTP(w, afs.makeHttpRequest(dreq, r))

	afs.attackFileHandle.Handle(dreq)

}

func (afs *AttackFileServer) Start() {

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", afs.host, afs.port),
		Handler: afs}

	log.Infof("FileServer starting on directory: %s\nListening on http://%s:%s\n", afs.attackFileDir, afs.host, afs.port)

	if err := server.ListenAndServe(); err != nil {

		log.Errorf("FileServer ListenAndServe Failed:%v", err)

	}

}

func NewAttackFileServer(attackFileHandle *handler.AttackFileServerHandle, attackFileDir string, host string, port int) *AttackFileServer {

	return &AttackFileServer{
		fhandler:         http.FileServer(http.Dir(attackFileDir)),
		attackFileHandle: attackFileHandle,
		attackFileDir:    attackFileDir,
		host:             host,
		port:             port,
	}

}
