package http

import (
	"fmt"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

type HttpTengo struct {
	name string
}

func (h *HttpTengo) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	panic("implement me")
}

func (h *HttpTengo) IsFalsy() bool {
	panic("implement me")
}

func (h *HttpTengo) Equals(another objects.Object) bool {
	panic("implement me")
}

func (h *HttpTengo) Copy() objects.Object {
	panic("implement me")
}

func (h *HttpTengo) TypeName() string {

	return "HttpTengo"
}

func (h *HttpTengo) String() string {

	return "HttpTengo"
}

/*for create a http client function*/
type HttpClientTengo struct {

	HttpTengo
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
		HttpTengo: HttpTengo{
			name:"httpclient",
		},
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

	if key == "send" {

		return &HttpClientSendTengo{
			HttpTengo: HttpTengo{name:"send"},
			clientTengo:    hc,
		},nil
	}

	return nil,fmt.Errorf("undefine http client method:%s",key)
}

/*for http send  function*/
type HttpClientSendTengo struct {

	HttpTengo
	clientTengo *HttpClientTengo
}

func (c *HttpClientSendTengo) Call(args ... objects.Object) (objects.Object,error){

	if len(args) !=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	reqTengo,ok := args[0].(*HttpRequestTengo)

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "user",
			Expected: "httpRequestTengo(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	res,err:= c.clientTengo.client.Send(reqTengo.req)

	if err!=nil {
		return nil,err
	}

	return &HttpResponseTengo{
		HttpTengo: HttpTengo{name:"response"},
		res:       res,
	},nil
}

/*for new http request */
type HttpRequestTengo struct {

	HttpTengo
	req *HttpRequest
}

func newHttpRequest(args ... objects.Object) (objects.Object,error) {

	if len(args)!=2 {

		return nil,fmt.Errorf("New Http Request Invalid args,must provide <method><uri>")
	}

	method, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "method",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	uri,ok := objects.ToString(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "uri",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	return &HttpRequestTengo{
		HttpTengo: HttpTengo{name:"request"},
		req:       NewHttpRequest(method,uri),
	},nil
}

func (req *HttpRequestTengo) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "auth":
		return &HttpRequestAuth{
			HttpTengo: HttpTengo{name: "auth"},
			reqTengo:  req,
		},nil

	case "addHeader":

		return &HttpRequestAddHeader{
			HttpTengo: HttpTengo{name:"addHeader"},
			reqTengo:  req,
		},nil

	case "putString":
		return &HttpRequestPutString{
			HttpTengo: HttpTengo{name:"putString"},
			reqTengo:  req,
		},nil

	case "putHex":
		return &HttpRequestPutHex{
			HttpTengo: HttpTengo{name:"putHex"},
			reqTengo:  req,
		},nil

	case "upload":
		return &HttpRequestUPload{
			HttpTengo: HttpTengo{name:"upload"},
			reqTengo:  req,
		},nil

	}

	return nil,fmt.Errorf("undefine http request method:%s",key)
}

type HttpRequestAuth struct {

	HttpTengo
	reqTengo *HttpRequestTengo
}

func (auth *HttpRequestAuth) Call(args ... objects.Object) (objects.Object,error){

	if len(args) !=2 {

		return nil,tengo.ErrWrongNumArguments
	}

	user,ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "user",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	passwd,ok := objects.ToString(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "passwd",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	auth.reqTengo.req.BasicAuth(user,passwd)

	return auth.reqTengo,nil
}
type HttpRequestAddHeader struct {

	HttpTengo
	reqTengo *HttpRequestTengo
}

func (ah *HttpRequestAddHeader) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	value,ok := objects.ToString(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "value",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	ah.reqTengo.req.AddHeader(name,value)

	return ah.reqTengo,nil
}


type HttpRequestPutString struct {

	HttpTengo
	reqTengo *HttpRequestTengo
}

func (ps *HttpRequestPutString) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	content, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "content",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	fromFile,ok := objects.ToBool(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "fromFile",
			Expected: "bool(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	ps.reqTengo.req.BodyString(content,fromFile)

	return ps.reqTengo,nil
}

type HttpRequestPutHex struct {

	HttpTengo
	reqTengo *HttpRequestTengo
}

func (ps *HttpRequestPutHex) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	content, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "content",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}



	ps.reqTengo.req.BodyHex(content)

	return ps.reqTengo,nil
}

type HttpRequestUPload struct {

	HttpTengo
	reqTengo *HttpRequestTengo
}

func (up *HttpRequestUPload) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 3 {

		return nil, tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	fpath,ok := objects.ToString(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "fpath",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	boundary,ok := objects.ToString(args[2])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "boundary",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}


	up.reqTengo.req.UPload(name,fpath,boundary)

	return up.reqTengo,nil
}

/*for http response*/
type HttpResponseTengo struct {

	HttpTengo
	res *HttpResponse
}

func (res *HttpResponseTengo) IndexGet(index objects.Object)(value objects.Object,err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}
	switch key {
	case "getStatusCode":

		return &HttpResponseStatusCode{
			HttpTengo: HttpTengo{name:"getStatusCode"},
			resTengo:  res,
		},nil

	case "getBodyAsByte":
		return &HttpResponseBodyByte{
			HttpTengo: HttpTengo{name:"getBodyAsByte"},
			resTengo:  res,
		},nil

	case "getBodyAsString":
		return &HttpResponseBodyString{
			HttpTengo: HttpTengo{name:"getBodyAsString"},
			resTengo:  res,
		},nil

	case "getProtocol":
		return &HttpResponseProtocol{
			HttpTengo: HttpTengo{name:"getProtocol"},
			resTengo:  res,
		},nil

	case "getHeader":
		return &HttpResponseHeader{
			HttpTengo: HttpTengo{name:"getHeader"},
			resTengo:  res,
		},nil

	case "getHeaders":
		return &HttpResponseHeaders{
			HttpTengo: HttpTengo{name:"getHeaders"},
			resTengo:  res,
		},nil
	}

	return nil,fmt.Errorf("Undefine http response function:%s",key)

}

type HttpResponseStatusCode struct {

	HttpTengo
	resTengo *HttpResponseTengo
}

func (rc *HttpResponseStatusCode) Call(args ... objects.Object) (objects.Object,error) {

	return objects.FromInterface(rc.resTengo.res.GetStatusCode())
}

type HttpResponseBodyByte struct {

	HttpTengo
	resTengo *HttpResponseTengo
}

func (rb *HttpResponseBodyByte) Call(args ... objects.Object) (objects.Object,error) {

	content,_:= rb.resTengo.res.GetBodyAsByte()
	return objects.FromInterface(content)
}

type HttpResponseBodyString struct {

	HttpTengo
	resTengo *HttpResponseTengo
}

func (rs *HttpResponseBodyString) Call(args ... objects.Object) (objects.Object,error) {

	content,_:= rs.resTengo.res.GetBodyAsString()
	return objects.FromInterface(content)
}

type HttpResponseProtocol struct {

	HttpTengo
	resTengo *HttpResponseTengo
}

func (rp *HttpResponseProtocol) Call(args ... objects.Object) (objects.Object,error) {

	content:= rp.resTengo.res.Protocol()
	return objects.FromInterface(content)
}

type HttpResponseHeader struct {

	HttpTengo
	resTengo *HttpResponseTengo
}

func (rh *HttpResponseHeader) Call(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	key,ok := objects.ToString(args[0])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return objects.FromInterface(rh.resTengo.res.GetHeaderValue(key))
}

type HttpResponseHeaders struct {

	HttpTengo
	resTengo *HttpResponseTengo
}

func (rhs *HttpResponseHeaders) Call(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	key,ok := objects.ToString(args[0])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return objects.FromInterface(rhs.resTengo.res.GetHeaderValues(key))
}


var moduleMap objects.Object = &objects.ImmutableMap{
	Value: map[string]objects.Object{
		"newHttpClient": &objects.UserFunction{
			Name:  "new_http_client",
			Value: newHttpClient,
		},
		"newHttpRequest": &objects.UserFunction{
			Name:  "new_http_request",
			Value: newHttpRequest,
		},

	},
}

func (HttpTengo) Import(moduleName string) (interface{}, error) {

	fmt.Println(moduleName)
	switch moduleName {
	case "http":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}
