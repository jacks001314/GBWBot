package http

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

/*for http response*/
type HttpResponseTengo struct {

	TengoObj
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

		return &HttpResponseMethod{
			TengoObj: TengoObj{name:"getStatusCode"},
			resTengo:  res,
		},nil

	case "getBodyAsByte":
		return &HttpResponseMethod{
			TengoObj: TengoObj{name:"getBodyAsByte"},
			resTengo:  res,
		},nil

	case "getBodyAsString":
		return &HttpResponseMethod{
			TengoObj: TengoObj{name:"getBodyAsString"},
			resTengo:  res,
		},nil

	case "getProtocol":
		return &HttpResponseMethod{
			TengoObj: TengoObj{name:"getProtocol"},
			resTengo:  res,
		},nil

	case "getHeader":
		return &HttpResponseMethod{
			TengoObj: TengoObj{name:"getHeader"},
			resTengo:  res,
		},nil

	case "getHeaders":
		return &HttpResponseMethod{
			TengoObj: TengoObj{name:"getHeaders"},
			resTengo:  res,
		},nil

	default:
		return nil,fmt.Errorf("Undefine http response function:%s",key)
	}
}

type HttpResponseMethod struct {

	TengoObj
	resTengo *HttpResponseTengo
}

func (hrm *HttpResponseMethod) getStatusCode()(objects.Object,error){

	return objects.FromInterface(hrm.resTengo.res.GetStatusCode())
}


func (hrm *HttpResponseMethod) getBodyAsByte() (objects.Object,error) {

	content,_:= hrm.resTengo.res.GetBodyAsByte()

	return objects.FromInterface(content)
}

func (hrm *HttpResponseMethod) getBodyAsString() (objects.Object,error) {

	content,_:= hrm.resTengo.res.GetBodyAsString()
	return objects.FromInterface(content)
}

func (hrm *HttpResponseMethod) getProtocol() (objects.Object,error) {

	content:= hrm.resTengo.res.Protocol()
	return objects.FromInterface(content)
}

func (hrm *HttpResponseMethod) getHeader(args ... objects.Object) (objects.Object,error) {

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

	return objects.FromInterface(hrm.resTengo.res.GetHeaderValue(key))
}

func (hrm *HttpResponseMethod) getHeaders(args ... objects.Object) (objects.Object,error) {

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

	return objects.FromInterface(hrm.resTengo.res.GetHeaderValues(key))
}

func (hrm *HttpResponseMethod) Call(args ... objects.Object) (objects.Object,error) {

	switch hrm.name {
	case "getStatusCode":
		return hrm.getStatusCode()

	case "getBodyAsByte":
		return hrm.getBodyAsByte()

	case "getBodyAsString":
		return hrm.getBodyAsString()

	case "getProtocol":
		return hrm.getProtocol()

	case "getHeader":
		return hrm.getHeader(args ...)
	case "getHeaders":
		return hrm.getHeaders(args ...)

	default:
		return nil,fmt.Errorf("unknown http response method:%s",hrm.name)

	}
}




