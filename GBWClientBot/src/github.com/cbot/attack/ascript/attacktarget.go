package ascript

import (
	"github.com/cbot/attack"
	"github.com/cbot/targets"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

type AttackTarget struct {

	attack.TengoObj

	target targets.Target

}

func (at *AttackTarget) IndexGet(index objects.Object)(value objects.Object,err error) {

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
		return objects.FromInterface(at.target.Port())

	case "app":
		return objects.FromInterface(at.target.App())

	case "proto":
		return objects.FromInterface(at.target.Proto())


	}

	return nil,nil
}