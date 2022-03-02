package node

type Config struct {
	TaskId string `json:"taskId"`

	SbotHost string `json:"sbotHost"`

	SbotRPCPort int `json:"sbotRPCPort"`

	SbotFileServerPort int `json:"sbotFileServerPort"`

	MaxThreads int `json:"maxThreads"`

	SourceCapacity int `json:"sourceCapacity"`

	AttackProcessCapacity int `json:"attackProcessCapacity"`
}
