package server

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func getArgsMap(content string) map[string]string {

	results := make(map[string]string)

	args := strings.Split(content,"&")

	for _,arg := range args {

		kv := strings.Split(arg,"=")

		if len(kv) <=1 {
			results[arg] = ""
		}else {

			results[kv[0]] = kv[1]
		}
	}

	return results
}

func DecodeCryptArgs(cryptData string) map[string]string {

	 hdata,err := hex.DecodeString(cryptData)

	 if err!=nil {

	 	return make(map[string]string,0)
	 }

	 content,err:= base64.StdEncoding.DecodeString(string(hdata))

	if err!=nil {

		return make(map[string]string,0)
	}

	return getArgsMap(string(content))

}

