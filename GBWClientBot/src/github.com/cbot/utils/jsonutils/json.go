package jsonutils

import (
	"encoding/json"
	"io/ioutil"
	"github.com/cbot/utils/fileutils"
)

func WriteJsonPretty(v interface{},fpath string) error {


	 if data,err := json.MarshalIndent(v,"","\t");err !=nil {
	 	return err
	 }else {

	 	return fileutils.WriteFile(fpath,data)
	 }
}


func UNMarshalFromFile(v interface{},fpath string) (err error) {

	if data,err := ioutil.ReadFile(fpath) ; err !=nil{
		return err
	}else {

		return json.Unmarshal(data,v)
	}
}