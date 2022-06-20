package server

import (
	"net/http"
)

type HttpFileServer struct {

	fhandler http.Handler
	fileDir string
	addr string
}

func NewHttpFileServer(fileDir,addr string) *HttpFileServer {

	return &HttpFileServer{
		fhandler: http.FileServer(http.Dir(fileDir)),
		fileDir:  fileDir,
		addr :addr,
	}
}

func (s *HttpFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Infof("Receive a  file download request,from:%s,url:%s", r.RemoteAddr, r.URL.Path)

	s.fhandler.ServeHTTP(w, r)

}

func (s *HttpFileServer) Start() {

	server := &http.Server{Addr:s.addr,
		Handler: s}

	log.Infof("FileServer starting on directory: %s\nListening on http://%s\n", s.fileDir, s.addr)

	if err := server.ListenAndServe(); err != nil {

		log.Errorf("FileServer ListenAndServe Failed:%v", err)

	}
}

