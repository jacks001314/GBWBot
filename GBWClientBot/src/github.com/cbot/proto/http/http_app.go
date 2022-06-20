package http

import (
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
	"strings"
)

type HttpApp struct {

	TengoObj
	proto string
	isSSL bool
	port  int

}

func isCommonPort(port int) bool {
	return port ==80 ||port == 443||port == 8080 ||port ==8443
}

func doDetect(host string,port int,isSSL bool,url string,key string,status int,timeout int64) bool {

	client := NewHttpClient(host,port,isSSL,timeout)

	request := NewHttpRequest("get",url)

	res,err := client.Send(request)

	if err!=nil {
		return false
	}

	if status!=0&&res.GetStatusCode() != status{
		return false
	}

	content,err:= res.GetBodyAsString()

	if err!=nil {
		return false
	}

	return strings.Contains(content,key)

}

func doDetectFromCommonPort(host string,url string,key string,status int,timeout int64) *HttpApp {

	if doDetect(host,80,false,url,key,status,timeout) {

		return &HttpApp{
			TengoObj: TengoObj{name:"http_app"},
			proto:    "http",
			isSSL:    false,
			port:     80,
		}
	}

	if doDetect(host,8080,false,url,key,status,timeout) {

		return &HttpApp{
			TengoObj: TengoObj{name:"http_app"},
			proto:    "http",
			isSSL:    false,
			port:     8080,
		}
	}

	if doDetect(host,443,true,url,key,status,timeout) {

		return &HttpApp{
			TengoObj: TengoObj{name:"http_app"},
			proto:    "https",
			isSSL:    true,
			port:     443,
		}
	}

	if doDetect(host,8443,true,url,key,status,timeout) {

		return &HttpApp{
			TengoObj: TengoObj{name:"http_app"},
			proto:    "https",
			isSSL:    true,
			port:     8443,
		}
	}

	return &HttpApp{
		TengoObj: TengoObj{name:"http_app"},
		proto:    "",
		isSSL:    false,
		port:     0,
	}

}

func DetectHttpApp(host string,port int,url string,key string,status int,timeout int64) *HttpApp {

	if !isCommonPort(port) {

		if doDetect(host,port,false,url,key,status,timeout) {

			return &HttpApp{
				TengoObj: TengoObj{name:"http_app"},
				proto:    "http",
				isSSL:    false,
				port:     port,
			}
		}

		if doDetect(host,port,true,url,key,status,timeout) {

			return &HttpApp{
				TengoObj: TengoObj{name:"http_app"},
				proto:    "https",
				isSSL:    true,
				port:     port,
			}
		}
	}

	return doDetectFromCommonPort(host,url,key,status,timeout)

}


func (h *HttpApp) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "proto":
		return objects.FromInterface(h.proto)

	case "port":
		return objects.FromInterface(h.port)

	case "isSSL":
		return objects.FromInterface(h.isSSL)

	}

	return nil, nil
}

