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
	return "IntegerType{" + t.Representation() + "}"
}

func (t IntegerType) Representation() string {
	return fmt.Sprintf("%d", t.Value)
}

func (t IntegerType) TypeHash() int {
	return integerType
}

func (t IntegerType) TypeName() string {
	return GetTypeName(t.TypeHash())
}
