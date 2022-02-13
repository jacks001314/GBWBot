package node

type Config struct {

	sbotHost string `json:"sbotHost"`

	sbotRPCPort int 	`json:"sbotRPCPort"`

	sbotFileServerPort int `json:"sbotFileServerPort"`


	NodeId string `json:"nodeId"`

}
