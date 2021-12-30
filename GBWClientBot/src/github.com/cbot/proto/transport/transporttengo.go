package transport

import (
	"fmt"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
	"time"
)

type TransportTengo struct {

	name string
}


func (t *TransportTengo) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	panic("implement me")
}

func (t *TransportTengo) IsFalsy() bool {
	panic("implement me")
}

func (t *TransportTengo) Equals(another objects.Object) bool {
	panic("implement me")
}

func (t *TransportTengo) Copy() objects.Object {
	panic("implement me")
}

func (t *TransportTengo) TypeName() string {

	return "TransportTengo"
}

func (t *TransportTengo) String() string {

	return "TransportTengo"
}

/*for transport connection*/
type ConnectionTengo struct {

	TransportTengo
	conn *Connection
}

/*
*
  args[0] ---network
  args[1] ---host
  args[2] ---port
  args[3] ---isSSL
  args[4] ---connectionTimeout
  args[5] ---readTimeout
  args[6] ---writeTimeout
*/

func newConnection(args ... objects.Object) (objects.Object,error){

	var conn *Connection
	var err error

	if len(args)!=7 {

		return nil,tengo.ErrWrongNumArguments
	}

	network,ok := objects.ToString(args[0])
	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "network",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	host,ok := objects.ToString(args[1])
	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "host",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	port,ok := objects.ToInt(args[2])
	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	isSSL,ok := objects.ToBool(args[3])
	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "isSSL",
			Expected: "bool(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	connTimeout,ok := objects.ToInt64(args[4])
	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "connTimeout",
			Expected: "int64(compatible)",
			Found:    args[4].TypeName(),
		}
	}

	readTimeout,ok := objects.ToInt64(args[5])
	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "readTimeout",
			Expected: "int64(compatible)",
			Found:    args[5].TypeName(),
		}
	}

	writeTimeout,ok := objects.ToInt64(args[6])
	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "writeTimeout",
			Expected: "int64(compatible)",
			Found:    args[6].TypeName(),
		}
	}

	if isSSL {

		conn,err = Dial(network,fmt.Sprintf("%s:%d",host,port),
				DialConnectTimeout(time.Duration(connTimeout)*time.Millisecond),
				DialReadTimeout(time.Duration(readTimeout)*time.Millisecond),
				DialWriteTimeout(time.Duration(writeTimeout)*time.Millisecond),
				DialTLSSkipVerify(true),
				DialUseTLS(true))

	}else {

		conn,err = Dial(network,fmt.Sprintf("%s:%d",host,port),
			DialConnectTimeout(time.Duration(connTimeout)*time.Millisecond),
			DialReadTimeout(time.Duration(readTimeout)*time.Millisecond),
			DialWriteTimeout(time.Duration(writeTimeout)*time.Millisecond))
	}

	if err!=nil {
		return nil,err
	}

	return &ConnectionTengo{
		TransportTengo: TransportTengo{name:"TransportTengo"},
		conn:           conn,
	},nil

}

func (ct *ConnectionTengo) IndexGet(index objects.Object)(value objects.Object,err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "close":
		return &ConnectionClose{
			TransportTengo: TransportTengo{name:"close"},
			connTengo:      ct,
		},nil

	case "flush":
		return &ConnectionFlush{
			TransportTengo: TransportTengo{name:"flush"},
			connTengo:      ct,
		},nil

	case "writeBytes":
		return &ConnectionWriteBytes{
			TransportTengo: TransportTengo{name:"writeBytes"},
			connTengo:      ct,
		},nil

	case "writeHex":
		return &ConnectionWriteHex{
			TransportTengo: TransportTengo{name:"writeHex"},
			connTengo:      ct,
		},nil

	case "writeString":
		return &ConnectionWriteString{
			TransportTengo: TransportTengo{name:"writeString"},
			connTengo:      ct,
		},nil

	case "readLine":
		return &ConnectionReadLine{
			TransportTengo: TransportTengo{name:"readLine"},
			connTengo:      ct,
		},nil

	case "readLineAsString":
		return &ConnectionReadLineAsString{
			TransportTengo: TransportTengo{name:"readLineAsString"},
			connTengo:      ct,
		},nil

	case "readBytes":
		return &ConnectionReadBytes{
			TransportTengo: TransportTengo{name:"readBytes"},
			connTengo:      ct,
		},nil

	case "readBytesAsString":
		return &ConnectionReadBytesAsString{
			TransportTengo: TransportTengo{name:"readBytesAsString"},
			connTengo:      ct,
		},nil


	}


	return nil,fmt.Errorf("Cannot support method:%s",key)
}

/*for connection close*/
type ConnectionClose struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (cc *ConnectionClose) Call(args ... objects.Object) (objects.Object,error) {
	cc.connTengo.conn.Close()
	return nil,nil
}

/*for write flush*/
type ConnectionFlush struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (cf *ConnectionFlush) Call(args ... objects.Object) (objects.Object,error) {

	cf.connTengo.conn.Flush()
	return nil,nil
}

/*for connection.WriteBytes*/
type ConnectionWriteBytes struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (cwb *ConnectionWriteBytes) Call(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	data,ok:= objects.ToByteSlice(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "[]byte(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return nil,cwb.connTengo.conn.WriteBytes(data)
}


/*for connection.WriteHex*/
type ConnectionWriteHex struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (cwh *ConnectionWriteHex) Call(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	data,ok:= objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return nil,cwh.connTengo.conn.WriteHex(data)
}

/*for connection.WriteString*/
type ConnectionWriteString struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (cws *ConnectionWriteString) Call(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	data,ok:= objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return nil,cws.connTengo.conn.WriteString(data)
}

/*for connection.ReadLine*/
type ConnectionReadLine struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (cr *ConnectionReadLine) Call(args ... objects.Object) (objects.Object,error) {

	data,err := cr.connTengo.conn.ReadLine()
	if err!=nil {

		return nil,err
	}

	return objects.FromInterface(data)
}

/*for connection.ReadLineASString*/
type ConnectionReadLineAsString struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (crs *ConnectionReadLineAsString) Call(args ... objects.Object) (objects.Object,error) {

	data,err := crs.connTengo.conn.ReadLine()
	if err!=nil {

		return nil,err
	}

	return objects.FromInterface(string(data))
}

/*for connection.ReadBytes*/
type ConnectionReadBytes struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (crbs *ConnectionReadBytes) Call(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	n,ok:= objects.ToInt(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "n",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	data,err := crbs.connTengo.conn.ReadBytes(n)
	if err!=nil {

		return nil,err
	}

	return objects.FromInterface(data)
}

/*for connection.ReadBytesASString*/
type ConnectionReadBytesAsString struct {

	TransportTengo
	connTengo *ConnectionTengo
}

func (crbs *ConnectionReadBytesAsString) Call(args ... objects.Object) (objects.Object,error) {

	if len(args)!=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	n,ok:= objects.ToInt(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "n",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	data,err := crbs.connTengo.conn.ReadBytes(n)
	if err!=nil {

		return nil,err
	}

	return objects.FromInterface(string(data))
}

var moduleMap objects.Object = &objects.ImmutableMap{
	Value: map[string]objects.Object{
		"newConnection": &objects.UserFunction{
			Name:  "newConnection",
			Value: newConnection,
		},
	},
}


func (TransportTengo) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "transport":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}

