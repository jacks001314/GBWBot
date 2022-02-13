package source

import (
	"errors"
)

var endError error = errors.New("Source Read Over!")

type Source interface {

	Put(target Target) error

	Start() error

	Stop()

	AtEnd()

	GetTypes() []string

	Name() string
}


