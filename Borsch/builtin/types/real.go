package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type RealType struct {
	Object

	Value float64
	package_ *PackageType
}

func NewRealType(value string) (RealType, error) {
	number, err := strconv.ParseFloat(strings.TrimSuffix(value, "f"), 64)
	if err != nil {
		return RealType{}, util.RuntimeError(err.Error())
	}

	return RealType{
		Value:    number,
		Object: *newBuiltinObject(
			RealTypeHash, map[string]Type{
				"__документ__": &NilType{}, // TODO: set doc
				"__пакет__":    BuiltinPackage,
			},
		),
		package_: BuiltinPackage,
	}, nil
}

func (t RealType) String() string {
	return strconv.FormatFloat(t.Value, 'f', -1, 64)
}

func (t RealType) Representation() string {
	return t.String()
}

func (t RealType) AsBool() bool {
	return t.Value != 0.0
}

func (t RealType) SetAttribute(name string, _ Type) (Type, error) {
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t RealType) Pow(other Type) (Type, error) {
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

func (t RealType) Plus() (Type, error) {
	return t, nil
}

func (t RealType) Minus() (Type, error) {
	return RealType{Value: -t.Value}, nil
}

func (t RealType) BitwiseNot() (Type, error) {
	return nil, nil
}

func (t RealType) Mul(other Type) (Type, error) {
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

func (t RealType) Div(other Type) (Type, error) {
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

func (t RealType) Mod(Type) (Type, error) {
	return nil, nil
}

func (t RealType) Add(other Type) (Type, error) {
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

func (t RealType) Sub(other Type) (Type, error) {
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

func (t RealType) BitwiseLeftShift(Type) (Type, error) {
	return nil, nil
}

func (t RealType) BitwiseRightShift(Type) (Type, error) {
	return nil, nil
}

func (t RealType) BitwiseAnd(Type) (Type, error) {
	return nil, nil
}

func (t RealType) BitwiseXor(Type) (Type, error) {
	return nil, nil
}

func (t RealType) BitwiseOr(Type) (Type, error) {
	return nil, nil
}

func (t RealType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case BoolType:
		rightVal := boolToFloat64(right.Value)
		if t.Value == rightVal {
			return 0, nil
		}

		if t.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case IntegerType:
		rightVal := float64(right.Value)
		if t.Value == rightVal {
			return 0, nil
		}

		if t.Value < rightVal {
			return -1, nil
		}

		return 1, nil
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
			"%s", t.GetTypeName(), right.GetTypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t RealType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t RealType) And(other Type) (Type, error) {
	return logicalAnd(t, other)
}

func (t RealType) Or(other Type) (Type, error) {
	return logicalOr(t, other)
}
