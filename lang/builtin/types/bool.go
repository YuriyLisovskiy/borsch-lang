package types

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/lang/util"
	"math"
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

func (t BoolType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case BoolType:
		if t.Value == right.Value {
			return 0, nil
		}
	default:
		return 0, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t BoolType) Add(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: boolToInt(t.Value) + boolToInt(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: boolToInt(t.Value) + o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: boolToFloat64(t.Value) + o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) Sub(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: boolToInt(t.Value) - boolToInt(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: boolToInt(t.Value) - o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: boolToFloat64(t.Value) - o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) Mul(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: boolToInt(t.Value) * boolToInt(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: boolToInt(t.Value) * o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: boolToFloat64(t.Value) * o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) Div(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		if o.Value {
			return RealType{
				Value: boolToFloat64(t.Value),
			}, nil
		}
	case IntegerType:
		if o.Value != 0 {
			return RealType{
				Value: boolToFloat64(t.Value) / float64(o.Value),
			}, nil
		}
	case RealType:
		if o.Value != 0.0 {
			return RealType{
				Value: boolToFloat64(t.Value) / o.Value,
			}, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func (t BoolType) Pow(other ValueType) (ValueType, error) {
		switch o := other.(type) {
		case RealType:
			return RealType{
				Value: math.Pow(boolToFloat64(t.Value), o.Value),
			}, nil
		case IntegerType:
			return IntegerType{
				Value: int64(math.Pow(boolToFloat64(t.Value), float64(o.Value))),
			}, nil
		case BoolType:
			return IntegerType{
				Value: int64(math.Pow(boolToFloat64(t.Value), boolToFloat64(o.Value))),
			}, nil
		default:
			return nil, nil
		}
}

func (t BoolType) Mod(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		if o.Value {
			return IntegerType{
				Value: boolToInt(t.Value) % boolToInt(o.Value),
			}, nil
		}
	case IntegerType:
		if o.Value != 0 {
			return IntegerType{
				Value: boolToInt(t.Value) % o.Value,
			}, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення за модулем на нуль")
}
