package types

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"strconv"
	"strings"
)

type RealType struct {
	Value float64
}

func NewRealType(value string) (RealType, error) {
	number, err := strconv.ParseFloat(strings.TrimSuffix(value, "f"), 64)
	if err != nil {
		return RealType{}, util.RuntimeError(err.Error())
	}

	return RealType{Value: number}, nil
}

func (t RealType) String() string {
	return "RealType{" + t.Representation() + "}"
}

func (t RealType) Representation() string {
	return fmt.Sprintf("%f", t.Value)
}

func (t RealType) TypeHash() int {
	return realType
}

func (t RealType) TypeName() string {
	return GetTypeName(t.TypeHash())
}
