package ssh

import (
	"fmt"
	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/compiler/token"
)

type TengoObj struct {

	name string
}

func (t *TengoObj) BinaryOp(op token.Token, rhs objects.Object) (objects.Object, error) {

	panic("implement me")
}

func (t *TengoObj) IsFalsy() bool {
	panic("implement me")
}

func (t *TengoObj) Equals(another objects.Object) bool {
	panic("implement me")
}

func (t *TengoObj) Copy() objects.Object {
	panic("implement me")
}

func (t *TengoObj) TypeName() string {

	return "SSHTengo"
}

func (t *TengoObj) String() string {

	return fmt.Sprintf("ScriptSourceTengoObject:%s",t.name)
}