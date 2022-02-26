package attack

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

type AttackProcess struct {

	TengoObj

	IP  string  `json:"ip"`
	Host string `json:"host"`
	Port int 	`json:"port"`
	Proto string `json:"proto"`
	App  string `json:"app"`
	OS   string `json:"os"`

	Name string `json:"name"`

	Type string `json:"type"`

	Status int `json:"status"`

	Payload string `json:"payload"`

	Result  string `json:"result"`

	Details string `json:"details"`
}


func (ap *AttackProcess) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "setIP":
		return &FieldSet {
			TengoObj: TengoObj{Name:"ip"},
			ap:    ap,
		},nil

	case "setHost":
		return &FieldSet {
			TengoObj: TengoObj{Name:"host"},
			ap:    ap,
		},nil

	case "setProto":
		return &FieldSet {
			TengoObj: TengoObj{Name:"proto"},
			ap:    ap,
		},nil

	case "setApp":
		return &FieldSet {
			TengoObj: TengoObj{Name:"app"},
			ap:    ap,
		},nil

	case "setOS":
		return &FieldSet {
			TengoObj: TengoObj{Name:"os"},
			ap:    ap,
		},nil

		case "setPort":
			return &FieldSet{
			TengoObj: TengoObj{Name:"port"},
			ap:    ap,
		},nil

	case "setName":
		return &FieldSet{
			TengoObj: TengoObj{Name:"name"},
			ap:    ap,
		},nil

	case "setType":
		return &FieldSet{
			TengoObj: TengoObj{Name:"type"},
			ap:    ap,
		},nil

	case "setStatus":
		return &FieldSet{
			TengoObj: TengoObj{Name:"status"},
			ap:    ap,
		},nil

	case "setPayload":
		return &FieldSet{
			TengoObj: TengoObj{Name:"payload"},
			ap:    ap,
		},nil

	case "setResult":
		return &FieldSet{
			TengoObj: TengoObj{Name:"result"},
			ap:    ap,
		},nil

	case "setDetails":
		return &FieldSet{
			TengoObj: TengoObj{Name:"details"},
			ap:    ap,
		},nil

	}


	return nil,fmt.Errorf("undefine set attack process field method:%s",key)
}

type FieldSet struct {
	TengoObj
	ap *AttackProcess
}

func (f *FieldSet) Call(args ... objects.Object) (objects.Object,error) {

	var sdata string
	var idata int
	var ok bool

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	if f.Name == "port" || f.Name == "status" {

		idata, ok = objects.ToInt(args[0])
		if !ok {

			return nil, tengo.ErrInvalidArgumentType{
				Name:     f.Name,
				Expected: "int(compatible)",
				Found:    args[0].TypeName(),
			}
		}

	}else {

		sdata, ok = objects.ToString(args[0])
		if !ok {

			return nil, tengo.ErrInvalidArgumentType{
				Name:     f.Name,
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}
	}

	switch f.Name {

	case "ip":
		f.ap.IP = sdata

	case "host":
		f.ap.Host = sdata

	case "proto":
		f.ap.Proto = sdata

	case "app":
		f.ap.App = sdata

	case "os":
		f.ap.OS = sdata

	case "port":
		f.ap.Port = idata

	case "name":
		f.ap.Name = sdata

	case "type":
		f.ap.Type = sdata

	case "status":
		f.ap.Status = idata

	case "payload":
		f.ap.Payload = sdata

	case "result":
		f.ap.Result = sdata

	case "details":
		f.ap.Details = sdata

	}

	return f.ap,nil
}
