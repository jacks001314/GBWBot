package attack

import "github.com/cbot/targets/source"

type Attack interface {
	Name() string

	DefaultPort() int

	DefaultProto() string

	Accept(target source.Target) bool

	Run(target source.Target) error

	PubProcess(process *AttackProcess)
}
