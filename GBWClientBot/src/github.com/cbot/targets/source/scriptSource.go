package source

import (

	"fmt"
	"github.com/cbot/proto/http"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"github.com/d5/tengo/v2"
	"io/ioutil"
	"sync"
)

type ScriptSource struct {

	TengoObj

	locker sync.Mutex


	/*tengo script instanse Compiled*/
	scomp  *script.Compiled

	/*readers for this source to read*/
	readers map[string]*SourceReader

	/*the reader types that this source can provided*/
	rtypesMap map[string]bool

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
	script.SetImports(mm)


	return script.Compile()

}

/*Create a script source by script content*/

func NewScriptSourceFromContent(rtypes []string,sdata []byte) (*ScriptSource,error) {

	com,err := scriptCompile(sdata)

	if err!= nil {

		return nil,err
	}

	rtypesMap := make(map[string]bool,0)

	for _,rtype:= range rtypes {

		rtypesMap[rtype] = true
	}

	return &ScriptSource{
		TengoObj: TengoObj{name: "scriptSource"},
		locker:   sync.Mutex{},
		scomp:    com,
		readers:  make(map[string]*SourceReader,0),
		rtypesMap: rtypesMap,
	},nil

}

/*create a script source by file*/
func NewScriptSourceFromFile(rtypes []string,fname string) (*ScriptSource,error){

	sdata,err:= ioutil.ReadFile(fname)

	if err!=nil {
		return nil,err
	}


	return NewScriptSourceFromContent(rtypes,sdata)

}


func (s *ScriptSource) canRead(rtypes []string) bool {

	for _,rtype := range rtypes {

		if _,ok := s.rtypesMap[rtype]; ok {

			return true
		}
	}

	return false
}

/*create a source reader
*@rtypes  ----that the types wanted to been read by reader
 @capacity ----channel capacity
 */
func (s *ScriptSource) OpenReader(name string,rtypes []string, capacity int) (*SourceReader,error)  {

	s.locker.Lock()
	defer s.locker.Unlock()

	if v,ok := s.readers[name]; ok {

		//existed
		return v,nil
	}

	if !s.canRead(rtypes) {

		return nil,fmt.Errorf("This Script source cannot provide rtyps:%v to been read,only provide rtypes:%v",
			rtypes,s.rtypesMap)
	}

	reader := NewSourceReader(name,rtypes,capacity)

	s.readers[name] = reader

	return reader,nil
}

func (s *ScriptSource) CloseReader(r *SourceReader) {

	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.readers,r.name)
}


func (s*ScriptSource) Put(entry SourceEntry) error{

	s.locker.Lock()
	defer s.locker.Unlock()

	for _,reader := range s.readers {

		reader.Push(entry)

	}


	return nil
}

func (s *ScriptSource) Start() error {

	go func () error {

		s.scomp.Set("scriptSource", s)
		s.scomp.Set("shajf","fuck")
		if err := s.scomp.Run(); err != nil {

			return err
		}

		return nil
	}()

	return nil
}

func (s* ScriptSource) Stop() {


}

func (s* ScriptSource) AtEnd() {

	for _,reader := range s.readers {

		reader.isEnd = true
	}

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

	entry := args[0].(SourceEntry)

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
