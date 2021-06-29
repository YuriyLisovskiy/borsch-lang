package types

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/lang/util"
	"math"
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

func (t RealType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t RealType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t RealType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case RealType:
		if t.Value == right.Value {
			return 0, nil
		}

		if t.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	default:
		return 0, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t RealType) Add(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return RealType{
			Value: t.Value + boolToFloat64(o.Value),
		}, nil
	case IntegerType:
		return RealType{
			Value: t.Value + float64(o.Value),
		}, nil
	case RealType:
		return RealType{
			Value: t.Value + o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t RealType) Sub(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return RealType{
			Value: t.Value - boolToFloat64(o.Value),
		}, nil
	case IntegerType:
		return RealType{
			Value: t.Value - float64(o.Value),
		}, nil
	case RealType:
		return RealType{
			Value: t.Value - o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t RealType) Mul(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return RealType{
			Value: t.Value * boolToFloat64(o.Value),
		}, nil
	case IntegerType:
		return RealType{
			Value: t.Value * float64(o.Value),
		}, nil
	case RealType:
		return RealType{
			Value: t.Value * o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t RealType) Div(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		if o.Value {
			return RealType{
				Value: t.Value,
			}, nil
		}
	case IntegerType:
		if o.Value != 0 {
			return RealType{
				Value: t.Value / float64(o.Value),
			}, nil
		}
	case RealType:
		if o.Value != 0.0 {
			return RealType{
				Value: t.Value / o.Value,
			}, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func (t RealType) Pow(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case RealType:
		return RealType{
			Value: math.Pow(t.Value, o.Value),
		}, nil
	case IntegerType:
		return RealType{
			Value: math.Pow(t.Value, float64(o.Value)),
		}, nil
	case BoolType:
		return RealType{
			Value: math.Pow(t.Value, boolToFloat64(o.Value)),
		}, nil
	default:
		return nil, nil
	}
}

func (t RealType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}
