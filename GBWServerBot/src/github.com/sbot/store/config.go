package store

type Config struct {

	DB 		string  `json:"db"`
	Table   string  `json:"table"`
	Host 	string  `json:"host"`
	Port    int 	`json:"port"`
	User    string  `json:"user"`
	Pass    string  `json:"pass"`
	Codes   string   `json:"codes"`
	
	Timeout uint64   `json:"timeout"`
}


