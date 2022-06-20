package jenkins

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

/*for create a http client function*/
type JeckinsClientTengo struct {

	TengoObj

	jenkinsClient *JenkinsClient
}


func newJeckinsClient(args ... objects.Object) (objects.Object,error) {

	if len(args)!=6 {

		return nil,fmt.Errorf("New Jeckins Client Invalid args,must provide <host><port><proto><user><passwd><buildPoolTime>")
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

	proto, ok := objects.ToString(args[2])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "proto",
			Expected: "string(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	user, ok := objects.ToString(args[3])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "user",
			Expected: "string(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	passwd, ok := objects.ToString(args[4])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "passwd",
			Expected: "string(compatible)",
			Found:    args[4].TypeName(),
		}
	}

	buildPoolTime,ok := objects.ToInt64(args[5])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "buildPoolTime",
			Expected: "int64(compatible)",
			Found:    args[5].TypeName(),
		}
	}

	jc,err := NewJenkinsClient(host,port,proto,user,passwd,uint64(buildPoolTime))

	if err!=nil {
		return nil,err
	}

	return &JeckinsClientTengo{
		TengoObj:      TengoObj{name:"JenkinsClient"},
		jenkinsClient: jc,
	},nil

}

func (jc *JeckinsClientTengo) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "createJob":
		return &JenkinsClientMethodTengo{
			TengoObj:    TengoObj{name:"createJob"},
			clientTengo: jc,
		},nil

	case "deleteJob":
		return &JenkinsClientMethodTengo{
			TengoObj:    TengoObj{name:"deleteJob"},
			clientTengo: jc,
		},nil

	case "buildJob":
		return &JenkinsClientMethodTengo{
			TengoObj:    TengoObj{name:"buildJob"},
			clientTengo: jc,
		},nil

	case "buildJobWaitResult":
		return &JenkinsClientMethodTengo{
			TengoObj:    TengoObj{name:"buildJobWaitResult"},
			clientTengo: jc,
		},nil

	default:
		return nil,fmt.Errorf("undefine jenkins client method:%s",key)
	}

}

type JenkinsClientMethodTengo struct {

	TengoObj
	clientTengo *JeckinsClientTengo
}

func (jcm *JenkinsClientMethodTengo) createJob(args ... objects.Object)(objects.Object,error) {

	if len(args) !=2 {

		return nil,tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	config, ok := objects.ToString(args[1])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "config",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	job,err := jcm.clientTengo.jenkinsClient.CreateJob(name,config)

	if err!=nil {
		return nil,err
	}

	return objects.FromInterface(job)
}

func (jcm *JenkinsClientMethodTengo) deleteJob(args ... objects.Object)(objects.Object,error) {

	if len(args) !=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	jcm.clientTengo.jenkinsClient.DeleteJob(name)

	return nil,nil
}

func (jcm *JenkinsClientMethodTengo) buildJob(args ... objects.Object)(objects.Object,error) {

	if len(args) !=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	result,err := jcm.clientTengo.jenkinsClient.BuildJob(name)

	if err!=nil {

		return objects.FromInterface(int64(-1))
	}

	return objects.FromInterface(result)
}

func (jcm *JenkinsClientMethodTengo) buildJobWaitResult(args ... objects.Object)(objects.Object,error) {

	if len(args) !=1 {

		return nil,tengo.ErrWrongNumArguments
	}

	name, ok := objects.ToString(args[0])
	if !ok {

		return nil,tengo.ErrInvalidArgumentType{
			Name:     "name",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	result,err := jcm.clientTengo.jenkinsClient.BuildJobWaitResult(name)

	if err!=nil {

		return objects.FromInterface(err.Error())
	}

	return objects.FromInterface(result)
}

func (jcm *JenkinsClientMethodTengo) Call(args ... objects.Object) (objects.Object,error){

	switch jcm.name {

	case "createJob":
		return jcm.createJob(args ...)

	case "deleteJob":
		return jcm.deleteJob(args ...)

	case "buildJob":
		return jcm.buildJob(args ...)

	case "buildJobWaitResult":
		return jcm.buildJobWaitResult(args ...)

	default:
		return nil,fmt.Errorf("unknown jenkins client method:%s",jcm.name)

	}

}

var moduleMap objects.Object = &objects.ImmutableMap{
	Value: map[string]objects.Object{
		"newJeckinsClient": &objects.UserFunction{
			Name:  "new_jenkins_client",
			Value: newJeckinsClient,
		},
	},
}

func (JeckinsClientTengo) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "jenkins":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}


