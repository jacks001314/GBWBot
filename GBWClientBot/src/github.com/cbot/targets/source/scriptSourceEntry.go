package source

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

type ScriptSourceEntry struct {

	TengoObj

	ip 		string
	host 	string
	port 	int
	proto 	string
	app		string
}

func newEntry(args ... objects.Object) (objects.Object,error) {


	return &ScriptSourceEntry {
		TengoObj: TengoObj{name:"SourceEntry"},
		ip:       "",
		host:     "",
		port:     0,
		proto:    "",
		app:      "",
	},nil
}

func (e *ScriptSourceEntry) IndexGet(index objects.Object)(value objects.Object,err error){

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
		return &EntrySetString {
			TengoObj: TengoObj{name:"ip"},
			entry:    e,
		},nil

	case "setHost":
		return &EntrySetString {
			TengoObj: TengoObj{name:"host"},
			entry:    e,
		},nil

	case "setProto":
		return &EntrySetString {
			TengoObj: TengoObj{name:"proto"},
			entry:    e,
		},nil

	case "setApp":
		return &EntrySetString {
			TengoObj: TengoObj{name:"app"},
			entry:    e,
		},nil

	case "setPort":
		return &EntrySetPort{
			TengoObj: TengoObj{name:"port"},
			entry:    e,
		},nil

	}


	return nil,fmt.Errorf("undefine http client method:%s",key)
}

type EntrySetString struct {
	TengoObj
	entry *ScriptSourceEntry
}

func (s *EntrySetString) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	content, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     s.name,
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	switch s.name {

	case "ip":

		s.entry.ip = content
	case "host":
		s.entry.host = content

	case "proto":
		s.entry.proto = content

	case "app":
		s.entry.app = content

	}

	return s.entry,nil
}

type EntrySetPort struct {
	TengoObj
	entry *ScriptSourceEntry
}

func (s *EntrySetPort) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	port, ok := objects.ToInt(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	s.entry.port = port

	return s.entry,nil
}


func (e *ScriptSourceEntry) IP() string {

	return e.ip
}

func (e *ScriptSourceEntry) Host() string {

	return e.host
}

func (e *ScriptSourceEntry) Port() int {

	return e.port
}

func (e *ScriptSourceEntry) Proto() string {

	return e.proto
}

func (e *ScriptSourceEntry) App() string {

	return e.app
}
