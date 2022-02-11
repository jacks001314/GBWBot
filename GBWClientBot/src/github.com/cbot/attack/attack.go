package attack

import (
	"github.com/cbot/targets"
)

type Attack interface {

	Name() string

	DefaultPort() int

	DefaultProto() string

	Accept(target targets.Target) bool

	Run(target targets.Target) error

	PubProcess(process *AttackProcess)

}
