package types

import (
	"github.com/YuriyLisovskiy/borsch/src/util"
	"strconv"
)

type BoolType struct {
	Value bool
}

func NewBoolType(value string) (BoolType, error) {
	switch value {
	case "істина":
		value = "t"
	case "хиба":
		value = "f"
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		return BoolType{}, util.RuntimeError(err.Error())
	}

	return BoolType{Value: boolean}, nil
}

func (t BoolType) String() string {
	if t.Value {
		return "істина"
	}

	return "хиба"
}

func (t BoolType) Representation() string {
	return t.String()
}

func (t BoolType) TypeHash() int {
	return BoolTypeHash
}

func (t BoolType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t BoolType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t BoolType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}
