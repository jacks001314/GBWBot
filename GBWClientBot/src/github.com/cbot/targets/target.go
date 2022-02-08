package targets

type Target interface {

	IP() 		string

	Host() 		string

	Port() 		int

	Proto() 	string

	App()		string

}

