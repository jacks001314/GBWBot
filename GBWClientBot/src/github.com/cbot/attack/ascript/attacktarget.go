package ascript

import (
	"github.com/cbot/attack"
	"github.com/cbot/targets/source"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

type AttackTarget struct {
	attack.TengoObj

	attack attack.Attack

	target source.Target
}

func newAttackTarget(at attack.Attack, target source.Target) *AttackTarget {

	return &AttackTarget{

		TengoObj: attack.TengoObj{Name: "AttackTarget"},
		attack:   at,
		target:   target,
	}
}

func (at *AttackTarget) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case "ip":
		return objects.FromInterface(at.target.IP())

	case "host":
		return objects.FromInterface(at.target.Host())

	case "port":

		port := at.target.Port()

		if port <= 0 {

			port = at.attack.DefaultPort()
		}

		return objects.FromInterface(port)

	case "app":
		return objects.FromInterface(at.target.App())

	case "proto":

		proto := at.target.Proto()

		if proto == "" {

			proto = at.attack.DefaultProto()
		}

		return objects.FromInterface(proto)

	}

	return nil, nil
}
