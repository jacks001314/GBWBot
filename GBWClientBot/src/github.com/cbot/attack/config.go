package attack

type Config struct {
	TaskId string `json:"taskId"`

	NodeId string `json:"nodeId"`

	MaxThreads int `json:"maxThreads"`

	SourceCapacity int `json:"sourceCapacity"`

	AttackProcessCapacity int `json:"attackProcessCapacity"`

	SBotHost string `json:"sbotHost"`

	SBotPort int `json:"sbotPort"`
}
