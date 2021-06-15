package types

import (
	"github.com/YuriyLisovskiy/borsch/src/util"
)

type NoneType struct {
}

func (t NoneType) String() string {
	return "Порожнеча"
}

func (t NoneType) Representation() string {
	return t.String()
}

func (t NoneType) TypeHash() int {
	return NoneTypeHash
}

func (t NoneType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t NoneType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}
