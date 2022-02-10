package ascript

import (
	"fmt"
	"github.com/cbot/attack"
	"github.com/cbot/proto/http"
	"github.com/cbot/targets"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"io/ioutil"
)

type AttackScript struct {

	attack.TengoObj

	name string

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
func NewAttackScriptFromContent(name string,data []byte) (*AttackScript,error) {

	com,err := scriptCompile(data)

	if err!= nil {

		return nil,err
	}

	return &AttackScript {
		TengoObj: attack.TengoObj{Name: name},
		name: name,
		attackTengo:    com,
	},nil

}

/*create an attack script  by file*/
func NewAttackScriptFromFile(name string,fname string) (*AttackScript,error){

	data,err:= ioutil.ReadFile(fname)

	if err!=nil {
		return nil,err
	}


	return NewAttackScriptFromContent(name,data)

}

func (a *AttackScript) Run( target targets.Target) {


}


func (a *AttackScript) PubProcess(process *attack.AttackProcess) {


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
