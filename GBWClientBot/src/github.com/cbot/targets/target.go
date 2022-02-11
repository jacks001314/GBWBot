package targets

import "github.com/cbot/targets/source"

type Target interface {

	IP() 		string

	Host() 		string

	Port() 		int

	Proto() 	string

	App()		string

	Source() source.Source
}

