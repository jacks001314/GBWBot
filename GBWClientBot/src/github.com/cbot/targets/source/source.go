package source

import (
	"errors"
	"github.com/cbot/targets"
)

var endError error = errors.New("Source Read Over!")

type Source interface {

	Put(target targets.Target) error

	Start() error

	Stop()

	AtEnd()

	GetTypes() []string

	Name() string
}


