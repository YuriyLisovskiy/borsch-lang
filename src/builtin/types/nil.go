package types

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

type NilType struct {
}

func (t NilType) String() string {
	return "нуль"
}

func (t NilType) Representation() string {
	return t.String()
}

func (t NilType) TypeHash() int {
	return NilTypeHash
}

func (t NilType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t NilType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t NilType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t NilType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
		return 0, nil
	default:
		return 0, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	}
}
