package source

import (
	"fmt"
	"github.com/cbot/proto/http"
	"github.com/cbot/targets"
	"github.com/cbot/targets/genip"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"github.com/d5/tengo/v2"
	"io/ioutil"
)

type ScriptSource struct {

	TengoObj

	name  string
	/*tengo script instanse Compiled*/
	scomp  *script.Compiled

	spool *SourcePool

	types []string
}


var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{
		"newEntry": &objects.UserFunction{
			Name:  "new_source_entry",
			Value: newEntry,
		},

	},
}

func (ScriptSource) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "source":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
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
	mm.Add("source", ScriptSource{})
	mm.Add("http", http.HttpTengo{})
	mm.Add("ipgen",genip.IPGen{})

	script.SetImports(mm)


	return script.Compile()

}

/*Create a script source by script content*/

func NewScriptSourceFromContent(spool *SourcePool,name string,rtypes []string,sdata []byte) (*ScriptSource,error) {

	com,err := scriptCompile(sdata)

	if err!= nil {

		return nil,err
	}

	return &ScriptSource{
		TengoObj: TengoObj{name:name},
		name: name,
		scomp:    com,
		spool: spool,
		types: rtypes,
	},nil

}

/*create a script source by file*/
func NewScriptSourceFromFile(spool *SourcePool,name string,rtypes []string,fname string) (*ScriptSource,error){

	sdata,err:= ioutil.ReadFile(fname)

	if err!=nil {
		return nil,err
	}


	return NewScriptSourceFromContent(spool,name,rtypes,sdata)

}

func (s *ScriptSource) Put(entry targets.Target) error{

	s.spool.put(s,entry)

	return nil
}

func (s *ScriptSource) Start() error {


	s.scomp.Set("scriptSource", s)

	if err := s.scomp.Run(); err != nil {

		return err

	}

	return nil
}

func (s* ScriptSource) Stop() {


}

func (s* ScriptSource)GetTypes() []string {


	return s.types
}


func (s* ScriptSource) AtEnd() {

	s.spool.StopSource(s)
}

func (s *ScriptSource) Name() string {

	return s.name
}

func (s *ScriptSource) IndexGet(index objects.Object)(value objects.Object,err error){

	key,ok := objects.ToString(index)

	if !ok {
		return nil,tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "put":

		return &SourcePut{
			TengoObj: TengoObj{name: "put"},
			source:   s,
		}, nil

	case "atEnd":

		return &SourceAtEnd{
			TengoObj: TengoObj{name: "atEnd"},
			source:   s,
		}, nil

	}
		return nil,fmt.Errorf("Unknown source method:%s",key)
}


type SourcePut struct {

	TengoObj

	source *ScriptSource
}

func (sp *SourcePut) Call(args ... objects.Object) (objects.Object,error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	entry := args[0].(*ScriptSourceEntry)

	entry.source = sp.source

	sp.source.Put(entry)

	return sp.source,nil
}

type SourceAtEnd struct {

	TengoObj

	source *ScriptSource
}

func (se *SourceAtEnd) Call(args ... objects.Object) (objects.Object,error) {

	se.source.AtEnd()

	return se.source,nil
}
