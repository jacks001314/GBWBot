package ssh

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

type TengoSSHClient struct {

	TengoObj

	sshClient *SSHClient

}

func newSSHClientWithPasswd(args ... objects.Object) (objects.Object,error) {

	if len(args) != 5 {

		return nil, tengo.ErrWrongNumArguments
	}

	host,ok:= objects.ToString(args[0])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "host",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	port,ok:= objects.ToInt(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	user,ok:= objects.ToString(args[2])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "user",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	passwd,ok:= objects.ToString(args[3])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "passwd",
			Expected: "string(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	timeout,ok:= objects.ToInt64(args[4])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "timeout",
			Expected: "int64(compatible)",
			Found:    args[4].TypeName(),
		}
	}

	sshClient,err := LoginWithPasswd(host,port,user,passwd,timeout)
	if err!=nil {

		return nil,nil
	}

	return &TengoSSHClient{
		TengoObj:  TengoObj{name:"SSHClient"},
		sshClient: sshClient,
	},nil

}

func newSSHClientWithKey(args ... objects.Object) (objects.Object,error) {

	if len(args) != 5 {

		return nil, tengo.ErrWrongNumArguments
	}

	host,ok:= objects.ToString(args[0])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "host",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	port,ok:= objects.ToInt(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	user,ok:= objects.ToString(args[2])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "user",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	key,ok:= objects.ToString(args[3])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "key",
			Expected: "string(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	timeout,ok:= objects.ToInt64(args[4])

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "timeout",
			Expected: "int64(compatible)",
			Found:    args[4].TypeName(),
		}
	}

	sshClient,err := LoginWithPrivKey(host,port,user,key, timeout)
	if err!=nil {

		return nil,nil
	}

	return &TengoSSHClient{
		TengoObj:  TengoObj{name:"SSHClient"},
		sshClient: sshClient,
	},nil

}

func (t *TengoSSHClient) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "runCmd":

		return &SSHClientMethod{
			TengoObj: TengoObj{name: "runCmd"},
			t:   t,
		}, nil

	case "upload":

		return &SSHClientMethod{
			TengoObj: TengoObj{name: "upload"},
			t:   t,
		}, nil

	case "download":
		return &SSHClientMethod{
			TengoObj: TengoObj{name: "download"},
			t:   t,
		}, nil

	case "close":
		return &SSHClientMethod{
			TengoObj: TengoObj{name: "close"},
			t:   t,
		}, nil

	}

	return nil,fmt.Errorf("Unknown ssh client method:%s",key)
}

type SSHClientMethod struct {

	TengoObj

	t *TengoSSHClient

}

func (m *SSHClientMethod) runCmd(cmdObj objects.Object) (objects.Object,error) {

	cmd,ok:= objects.ToString(cmdObj)

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "cmd",
			Expected: "string(compatible)",
			Found:    cmdObj.TypeName(),
		}

	}

	res,err := m.t.sshClient.RunCmd(cmd)
	if err!=nil {

		return nil,err
	}

	return objects.FromInterface(string(res))
}


func (m *SSHClientMethod) upload(fpathObj objects.Object,remoteDirObj objects.Object) (objects.Object,error) {

	fpath,ok:= objects.ToString(fpathObj)

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "fpath",
			Expected: "string(compatible)",
			Found:    fpathObj.TypeName(),
		}

	}

	remoteDir,ok:= objects.ToString(remoteDirObj)

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "remoteDir",
			Expected: "string(compatible)",
			Found:    remoteDirObj.TypeName(),
		}

	}

	sftpClient,err := NewSftpClient(m.t.sshClient)
	if err!=nil {

		return nil,err
	}

	defer sftpClient.Close()

	return nil,sftpClient.UPloadFile(fpath,remoteDir)

}

func (m *SSHClientMethod) download(rfileObj objects.Object,localDirObj objects.Object) (objects.Object,error) {

	rfile,ok:= objects.ToString(rfileObj)

	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "rfile",
			Expected: "string(compatible)",
			Found:    rfileObj.TypeName(),
		}

	}

	localDir,ok:= objects.ToString(localDirObj)


	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "localDir",
			Expected: "string(compatible)",
			Found:    localDirObj.TypeName(),
		}

	}

	sftpClient,err := NewSftpClient(m.t.sshClient)
	if err!=nil {

		return nil,err
	}

	defer sftpClient.Close()

	return nil,sftpClient.DownloadFile(rfile,localDir)

}


func (m *SSHClientMethod) Call(args ... objects.Object) (objects.Object,error) {

	switch m.name {

	case "runCmd":

		if len(args)!=1 {

			return nil,tengo.ErrWrongNumArguments
		}

		return m.runCmd(args[0])

	case "upload":

		if len(args)!=2 {

			return nil,tengo.ErrWrongNumArguments
		}

		return m.upload(args[0],args[1])

	case "download":
		if len(args)!=2 {

			return nil,tengo.ErrWrongNumArguments
		}

		return m.download(args[0],args[1])

	case "close":

		m.t.sshClient.Close()

		return nil,nil
	}


	return nil,fmt.Errorf("Unknown ssh client method:%s",m.name)

}

var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{

		"newSSHClientWithPass": &objects.UserFunction{
			Name:  "new_ssh_client_with_passwd",
			Value: newSSHClientWithPasswd,
		},

		"newSSHClientWithKey": &objects.UserFunction{
			Name:  "new_ssh_client_with_key",
			Value: newSSHClientWithKey,
		},
	},
}

func (TengoSSHClient) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "ssh":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}