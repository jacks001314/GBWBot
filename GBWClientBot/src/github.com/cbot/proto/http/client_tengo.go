package http

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

/*for create a http client function*/
type HttpClientTengo struct {

	TengoObj

	client *HttpClient
}

func newHttpClient(args ... objects.Object) (objects.Object,error) {

	if len(args)!=4 {

		return nil,fmt.Errorf("New Http Client Invalid args,must provide <host><port><isSSL><timeout>")
	}

	host, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "host",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	port,ok := objects.ToInt(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	isSSL,ok := objects.ToBool(args[2])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "bool(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	timeOut,ok := objects.ToInt64(args[3])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int64(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	return &HttpClientTengo{
		TengoObj: TengoObj{name:"HttpClient"},
		client:   NewHttpClient(host,port,isSSL,timeOut),
	},nil
}

func (hc *HttpClientTengo) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "send":
		return &HttpClientMethodTengo{
			TengoObj:    TengoObj{name:"send"},
			clientTengo: hc,
		},nil

	default:
		return nil,fmt.Errorf("undefine http client method:%s",key)
	}

}

/*for http send  function*/
type HttpClientMethodTengo struct {

	TengoObj
	clientTengo *HttpClientTengo
}

func (hcm *HttpClientMethodTengo) sendRequest(args ... objects.Object)(objects.Object,error) {

	if len(args) !=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	request,ok := args[0].(*HttpRequestTengo)

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "user",
			Expected: "httpRequestTengo(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	response,err:= hcm.clientTengo.client.Send(request.req)

	if err!=nil {
		return nil,err
	}

	return &HttpResponseTengo {
		TengoObj:TengoObj{name:"response"},
		res:response,
	},nil
}

func (hcm *HttpClientMethodTengo) Call(args ... objects.Object) (objects.Object,error){

	switch hcm.name {
	case "send":
		return hcm.sendRequest(args...)

	default:
		return nil,fmt.Errorf("unknown http client method:%s",hcm.name)

	}

}



