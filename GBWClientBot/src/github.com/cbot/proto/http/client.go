package http

import (
	"net/http"
	"time"
)

type HttpClient struct {

	client 				   *http.Client
	proto 	string
	host 	string
	port 	int

}

func NewHttpClient(host string,port int,isSSL bool,timeout int64) (httpClient *HttpClient) {

	tr := http.DefaultTransport.(*http.Transport)
	var proto string = "http"

	if isSSL {
		proto = "https"
		if port == 443 {
			port = 0
		}
	}else {

		if port == 80 {
			port = 0
		}
	}

	client := http.Client{
		Transport:     tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(timeout)*time.Millisecond ,
	}

	return &HttpClient{
		client: &client,
		proto: proto,
		host: host,
		port: port,
	}

}


func (c *HttpClient) Send(req *HttpRequest) (*HttpResponse,error){

	var request *http.Request
	var err error
	var httpResponse HttpResponse
	var res *http.Response

	/*make http request*/
	if request,err = req.Build(c.proto,c.host,c.port); err!=nil {

		return nil,err
	}

	if res,err = c.client.Do(request);err!=nil {
		return nil,err
	}

	httpResponse.resp = res
	
	return &httpResponse,nil
}
