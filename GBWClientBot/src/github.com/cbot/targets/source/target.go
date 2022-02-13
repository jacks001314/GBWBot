package source

type Target interface {

	IP() 		string

	Host() 		string

	Port() 		int

	Proto() 	string

	App()		string

	Source() Source
}

