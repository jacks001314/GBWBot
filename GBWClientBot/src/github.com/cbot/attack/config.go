package attack

type Config struct {
	MaxThreads int `json:"maxThreads"`

	SourceCapacity int `json:"sourceCapacity"`

	AttackProcessCapacity int `json:"attackProcessCapacity"`

	SBotHost string `json:"sbotHost"`

	SBotPort int `json:"sbotPort"`
}
