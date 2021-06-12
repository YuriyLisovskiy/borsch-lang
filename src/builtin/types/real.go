package types

import (
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
	return strconv.FormatFloat(t.Value, 'f', -1, 64)
}

func (t RealType) Representation() string {
	return t.String()
}

func (t RealType) TypeHash() int {
	return RealTypeHash
}

func (t RealType) TypeName() string {
	return GetTypeName(t.TypeHash())
}
