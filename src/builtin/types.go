package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"strconv"
	"strings"
)

const (
	noneType = iota
	realNumberType
	integerNumberType
	stringType
)

type ValueType interface {
	String() string
	Representation() string
	TypeHash() int
	TypeName() string
}

// NoneType represents none type.
type NoneType struct {
}

func (NoneType) String() string {
	return "Порожнеча"
}

func (t NoneType) Representation() string {
	return "\"" + t.String() + "\""
}

func (t NoneType) TypeHash() int {
	return noneType
}

func (t NoneType) TypeName() string {
	return "ніякий"
}

// RealNumberType represents numbers as float64
type RealNumberType struct {
	Value float64
}

func NewRealNumberType(value string) (RealNumberType, error) {
	number, err := strconv.ParseFloat(strings.TrimSuffix(value, "f"), 64)
	if err != nil {
		return RealNumberType{}, util.RuntimeError(err.Error())
	}

	return RealNumberType{Value: number}, nil
}

func (t RealNumberType) String() string {
	return fmt.Sprintf("%f", t.Value)
}

func (t RealNumberType) Representation() string {
	return "\"" + t.String() + "\""
}

func (t RealNumberType) TypeHash() int {
	return realNumberType
}

func (t RealNumberType) TypeName() string {
	return "дійсне число"
}

// IntegerNumberType represents numbers as float64
type IntegerNumberType struct {
	Value int64
}

func NewIntegerNumberType(value string) (IntegerNumberType, error) {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return IntegerNumberType{}, util.RuntimeError(err.Error())
	}

	return IntegerNumberType{Value: number}, nil
}

func (t IntegerNumberType) String() string {
	return fmt.Sprintf("%d", t.Value)
}

func (t IntegerNumberType) Representation() string {
	return "\"" + t.String() + "\""
}

func (t IntegerNumberType) TypeHash() int {
	return integerNumberType
}

func (t IntegerNumberType) TypeName() string {
	return "ціле число"
}

// StringType is string representation.
type StringType struct {
	Value string
}

func (t StringType) String() string {
	return t.Value
}

func (t StringType) Representation() string {
	return "\"" + t.String() + "\""
}

func (t StringType) TypeHash() int {
	return stringType
}

func (t StringType) TypeName() string {
	return "рядок"
}
