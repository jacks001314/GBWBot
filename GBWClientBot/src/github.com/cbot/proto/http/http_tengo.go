package http

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
	"net/url"
)

type HttpTengo struct {

	TengoObj
}

func urlEncode(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	urlRaw, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "url",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return objects.FromInterface(url.QueryEscape(urlRaw))
}

func urlDecode(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	urlRaw, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "url",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	durl,err := url.QueryUnescape(urlRaw)

	if err!=nil {
		return objects.FromInterface(urlRaw)
	}

	return objects.FromInterface(durl)
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

		"urlEncode": &objects.UserFunction{
			Name:  "urlEncode",
			Value: urlEncode,
		},

		"urlDecode": &objects.UserFunction{
			Name:  "urlDecode",
			Value: urlDecode,
		},

	},
}

func (HttpTengo) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "http":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}
