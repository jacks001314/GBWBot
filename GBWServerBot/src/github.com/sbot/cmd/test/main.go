package main

import (
	"os"
	"text/template"
)

type Test struct {


	Port int
	Proto string
	App string
}

func main()  {

	tpath := `D:\shajf_dev\self\GBWBot\GBWServerBot\src\github.com\sbot\scripts\source\local.tengo.tpl`

	test := &Test{
		Port:  22,
		Proto: "ssh",
		App:   "ssh",
	}

	t,_:=template.ParseFiles(tpath)

	t.Execute(os.Stdout,test)
}
