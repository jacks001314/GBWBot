package http

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

/*for new http request */
type HttpRequestTengo struct {

	TengoObj
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
		TengoObj: TengoObj{name:"HttpRequest"},
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
		return &HttpRequestMethod{
			TengoObj: TengoObj{name: "auth"},
			reqTengo:  req,
		},nil

	case "addHeader":

		return &HttpRequestMethod{
			TengoObj: TengoObj{name:"addHeader"},
			reqTengo:  req,
		},nil

	case "putString":
		return &HttpRequestMethod{
			TengoObj: TengoObj{name:"putString"},
			reqTengo:  req,
		},nil

	case "putHex":
		return &HttpRequestMethod{
			TengoObj: TengoObj{name:"putHex"},
			reqTengo:  req,
		},nil

	case "upload":
		return &HttpRequestMethod{
			TengoObj: TengoObj{name:"upload"},
			reqTengo:  req,
		},nil

	default:
		return nil,fmt.Errorf("undefine http request method:%s",key)
	}

}

type HttpRequestMethod struct {

	TengoObj
	reqTengo *HttpRequestTengo
}

func (hrm *HttpRequestMethod) makeRequestAuth(args ... objects.Object) (objects.Object,error){

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

	hrm.reqTengo.req.BasicAuth(user,passwd)

	return hrm.reqTengo,nil
}

func (hrm *HttpRequestMethod) addHeader(args ... objects.Object) (objects.Object,error) {

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

	hrm.reqTengo.req.AddHeader(name,value)

	return hrm.reqTengo,nil
}

func (hrm *HttpRequestMethod) putString(args ... objects.Object) (objects.Object,error) {

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

	hrm.reqTengo.req.BodyString(content,fromFile)

	return hrm.reqTengo,nil
}

func (hrm *HttpRequestMethod) putHex(args ... objects.Object) (objects.Object,error) {

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

	hrm.reqTengo.req.BodyHex(content)

	return hrm.reqTengo,nil
}

func (hrm *HttpRequestMethod) upload(args ... objects.Object) (objects.Object,error) {

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


	hrm.reqTengo.req.UPload(name,fpath,boundary)

	return hrm.reqTengo,nil
}

func (hrm *HttpRequestMethod) Call(args ... objects.Object) (objects.Object,error){

	switch hrm.name {

	case "auth":
		return hrm.makeRequestAuth(args ...)

	case "addHeader":
		return hrm.addHeader(args ...)

	case "putString":
		return hrm.putString(args ...)

	case "putHex":
		return hrm.putHex(args ...)

	case "upload":
		return hrm.upload(args ...)

	default:
		return nil,fmt.Errorf("unknown http request method:%s",hrm.name)
	}
}

