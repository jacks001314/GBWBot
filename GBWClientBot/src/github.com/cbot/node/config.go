package node

type Config struct {
	TaskId string `json:"taskId"`

	PNodeId string `json:"pnodeId"`

	AttackType string `json:"attackType"`

	SbotHost string `json:"sbotHost"`

	SbotRPCPort int `json:"sbotRPCPort"`

	SbotFileServerPort int `json:"sbotFileServerPort"`

	SbotLdapServerPort int `json:"sbotLdapServerPort"`

	MaxThreads int `json:"maxThreads"`

	SourceCapacity int `json:"sourceCapacity"`

	AttackProcessCapacity int `json:"attackProcessCapacity"`
}
