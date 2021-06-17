package types

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"strconv"
)

type IntegerType struct {
	Value int64
}

func NewIntegerType(value string) (IntegerType, error) {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return IntegerType{}, util.RuntimeError(err.Error())
	}

	return IntegerType{Value: number}, nil
}

func (t IntegerType) String() string {
	return fmt.Sprintf("%d", t.Value)
}

func (t IntegerType) Representation() string {
	return t.String()
}

func (t IntegerType) TypeHash() int {
	return IntegerTypeHash
}

func (t IntegerType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t IntegerType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t IntegerType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}
