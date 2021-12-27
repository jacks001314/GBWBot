package main

import (
	"context"
	"fmt"
	"github.com/cbot/proto/http"
	"github.com/d5/tengo/compiler/token"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/script"
	"github.com/d5/tengo/stdlib"
	"io/ioutil"
)

func testHttp(){

	host  := "www.163.com"
	port  := 443

	client := http.NewHttpClient(host,port,true,10000)
	request := http.NewHttpRequest("get","/").AddHeader("User-Agent","GOClient")

	if resp,err:= client.Send(request);err!=nil {

		fmt.Println(err)
	}else {

		fmt.Println(resp.GetHeaderValue("Content-Type"))
		fmt.Println(resp.GetStatusCode())
		fmt.Println(resp.GetBodyAsString())
	}
}

type Test struct {

	//tengo.ObjectImpl
	c int

}

func (t *Test) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {
	panic("implement me")
}

func (t *Test) IsFalsy() bool {
	panic("implement me")
}

func (t *Test) Equals(another objects.Object) bool {
	panic("implement me")
}

func (t *Test) Copy() objects.Object {
	panic("implement me")
}

func (t *Test) TypeName() string {

	return "test"
}

func (t *Test) String() string {

	return "test"
}

func (t *Test) Call(args ... objects.Object) (objects.Object,error){

	fmt.Println(args[1])
	fmt.Println(t.c)
	tt := args[2].(*Test)
	fmt.Println(tt.c)

	return &objects.Int{
		Value:      12,
	},nil

}

func (t *Test) CanCall() bool {

	return true
}

func newTest(args ...objects.Object) (objects.Object, error) {

	var t Test
	t.c = 1234
	fmt.Println(",,,,,,,,,,,,,,,,,,,,,,,,,,,")
	return &t,nil
}

type TestGetCall struct {
	Test
	fmap map[string]interface{}
	value string
}

func newGetCall(args ...objects.Object) (objects.Object, error) {

	return &TestGetCall{
		Test:  Test{},
		fmap:map[string]interface{}{
		"get":newGetFunc,
		"set":newSetFunc,
	},
	value: "",
	},nil
}

func (tc *TestGetCall)IndexGet(index objects.Object) (value objects.Object, err error){

	k,_:= objects.ToString(index)

	if k == "get" {

		return &GetFunc{
			Test: Test{},
			name: "get",
			tgc: tc,
		},nil
	}

	if k== "set" {

		return &SetFunc{
			Test: Test{},
			name: "set",
			tgc:  tc,
		},nil
	}

	return nil,nil
}

type GetFunc struct {

	Test
	name string
	tgc *TestGetCall
}

func newGetFunc() *GetFunc {

	return &GetFunc{
		Test: Test{},
		name: "get",
	}
}

func (tc *GetFunc) Call(args ... objects.Object) (objects.Object,error){

	return objects.FromInterface(tc.tgc.value)
}

type SetFunc struct {

	Test
	name string
	tgc *TestGetCall
}

func newSetFunc() *SetFunc {

	return &SetFunc{
		Test: Test{},
		name: "set",
	}
}

func (tc *SetFunc) Call(args ... objects.Object) (objects.Object,error){

	tc.tgc.value ,_= objects.ToString(args[0])

	return nil,nil
}


var moduleMap objects.Object = &objects.ImmutableMap{
	Value: map[string]objects.Object{
		"newTest": &objects.UserFunction{
			Name:  "new_test",
			Value: newTest,
		},
		"newCall": &objects.UserFunction{
			Name:  "new_tcall",
			Value: newGetCall,
		},
	},
}

func (Test) Import(moduleName string) (interface{}, error) {

	fmt.Println(moduleName)
	switch moduleName {
	case "test":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}


func testScript(){

	path:="D:\\shajf_dev\\self\\GBWBot\\GBWClientBot\\src\\github.com\\cbot\\cmd\\test\\t.tengo"
	data,_:=ioutil.ReadFile(path)
	script := script.New(data)

	mm := objects.NewModuleMap()
	mm.Add("http", http.HttpTengo{})
	mm.AddMap(stdlib.GetModuleMap("fmt"))
	script.SetImports(mm)

	// run the script
	_, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}

	//objects.Map{}

}

func main() {

	testScript()
	//testHttp()
}

