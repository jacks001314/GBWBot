package ascript

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/proto/http"
	"github.com/cbot/targets"
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
	attackTengo  *script.Compiled

}


/*compile tengo script*/
func scriptCompile(sdata []byte) (*script.Compiled, error) {

	script := script.New(sdata)

	script.Add("scriptSource",nil)

	mm := objects.NewModuleMap()

	/*add all stdlibs*/
	builtinMaps := objects.NewModuleMap()
	for name,im:= range stdlib.BuiltinModules {
		builtinMaps.AddBuiltinModule(name,im)
	}

	mm.AddMap(builtinMaps)
	mm.Add("attack", AttackScript{})
	mm.Add("http", http.HttpTengo{})


	script.SetImports(mm)

	return script.Compile()

}

/*Create an attack script  by script content*/
func NewAttackScriptFromContent(attackTasks *attack.AttackTasks,
	name string,
	attackType string,
	defaultPort int,
	defaultProto string,
	data []byte) (*AttackScript,error) {

	com,err := scriptCompile(data)

	if err!= nil {

		return nil,err
	}

	return &AttackScript{
		TengoObj:     attack.TengoObj{Name:name},
		attackTasks:  attackTasks,
		name:         name,
		attackType:   attackType,
		defaultPort:  defaultPort,
		defaultProto: defaultProto,
		attackTengo:  com,
	},nil


}

/*create an attack script  by file*/
func NewAttackScriptFromFile(attackTasks *attack.AttackTasks,
	name string,
	attackType string,
	defaultPort int,
	defaultProto string,fname string) (*AttackScript,error){

	data,err:= ioutil.ReadFile(fname)

	if err!=nil {
		return nil,err
	}


	return NewAttackScriptFromContent(attackTasks,name,attackType,defaultPort,defaultProto,data)

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

func (as *AttackScript)Accept(target targets.Target) bool {

	types := target.Source().GetTypes()

	for _,t:= range types {

		if strings.EqualFold(t,as.attackType){

			return true
		}
	}

	return false
}


func (as *AttackScript) Run(target targets.Target) error {

	defer as.attackTasks.PubSyn()

	attackTarget := newAttackTarget(as,target)

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

		TengoObj: attack.TengoObj{Name:"AttackProcess"},
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
	},nil

}


func (as *AttackScript) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "pubProcess":

		return &PubProcessMethod{
			TengoObj: attack.TengoObj{Name: "pubProcess"},
			as:   as,
		}, nil

	}

	return nil,fmt.Errorf("Unknown Attack script method:%s",key)
}


type PubProcessMethod struct {

	attack.TengoObj

	as *AttackScript
}

func (pp *PubProcessMethod) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	ap := args[0].(*attack.AttackProcess)

	pp.as.attackTasks.PubAttackProcess(ap)

	return pp.as,nil
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
