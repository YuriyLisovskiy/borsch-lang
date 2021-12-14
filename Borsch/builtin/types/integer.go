package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type IntegerType struct {
	Value int64
	object   *ObjectType
	package_ *PackageType
}

func NewIntegerType(value string) (IntegerType, error) {
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return IntegerType{}, util.RuntimeError(err.Error())
	}

	return IntegerType{
		Value:    number,
		object: newObjectType(
			IntegerTypeHash, map[string]ValueType{
				"__документ__": &NilType{}, // TODO: set doc
				"__пакет__":    BuiltinPackage,
			},
		),
		package_: BuiltinPackage,
	}, nil
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

func (t IntegerType) AsBool() bool {
	return t.Value != 0
}

func (t IntegerType) GetAttr(name string) (ValueType, error) {
	return t.object.GetAttribute(name)
}

func (t IntegerType) SetAttr(name string, value ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t IntegerType) Pow(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case RealType:
		return RealType{
			Value: math.Pow(float64(t.Value), o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: int64(math.Pow(float64(t.Value), float64(o.Value))),
		}, nil
	case BoolType:
		return IntegerType{
			Value: int64(math.Pow(float64(t.Value), boolToFloat64(o.Value))),
		}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) Plus() (ValueType, error) {
	return t, nil
}

func (t IntegerType) Minus() (ValueType, error) {
	return IntegerType{Value: -t.Value}, nil
}

func (t IntegerType) BitwiseNot() (ValueType, error) {
	return IntegerType{Value: ^t.Value}, nil
}

func (t IntegerType) Mul(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: t.Value * boolToInt64(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: t.Value * o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: float64(t.Value) * o.Value,
		}, nil
	case StringType:
		count := int(t.Value)
		if count <= 0 {
			return StringType{Value: ""}, nil
		}

		return StringType{
			Value: strings.Repeat(o.Value, count),
		}, nil
	case ListType:
		count := int(t.Value)
		list := NewListType()
		if count > 0 {
			for c := 0; c < count; c++ {
				list.Values = append(list.Values, o.Values...)
			}
		}

		return list, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) Div(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		if o.Value {
			return RealType{
				Value: float64(t.Value),
			}, nil
		}
	case IntegerType:
		if o.Value != 0 {
			return RealType{
				Value: float64(t.Value) / float64(o.Value),
			}, nil
		}
	case RealType:
		if o.Value != 0.0 {
			return RealType{
				Value: float64(t.Value) / o.Value,
			}, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення на нуль")
}

func (t IntegerType) Mod(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		if o.Value {
			return IntegerType{
				Value: t.Value % boolToInt64(o.Value),
			}, nil
		}
	case IntegerType:
		if o.Value != 0 {
			return IntegerType{
				Value: t.Value % o.Value,
			}, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення за модулем на нуль")
}

func (t IntegerType) Add(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: t.Value + boolToInt64(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: t.Value + o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: float64(t.Value) + o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) Sub(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: t.Value - boolToInt64(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: t.Value - o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: float64(t.Value) - o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) BitwiseLeftShift(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: t.Value << boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: t.Value << o.Value}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) BitwiseRightShift(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: t.Value >> boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: t.Value >> o.Value}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) BitwiseAnd(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: t.Value & boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: t.Value & o.Value}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) BitwiseXor(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: t.Value ^ boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: t.Value ^ o.Value}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) BitwiseOr(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: t.Value | boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: t.Value | o.Value}, nil
	default:
		return nil, nil
	}
}

func (t IntegerType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case BoolType:
		rightVal := boolToInt64(right.Value)
		if t.Value == rightVal {
			return 0, nil
		}

		if t.Value < rightVal {
			return -1, nil
		}

		return 1, nil
	case IntegerType:
		if t.Value == right.Value {
			return 0, nil
		}

		if t.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	case RealType:
		leftVal := float64(t.Value)
		if leftVal == right.Value {
			return 0, nil
		}

		if leftVal < right.Value {
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

func (t IntegerType) Not() (ValueType, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t IntegerType) And(other ValueType) (ValueType, error) {
	return logicalAnd(t, other)
}

func (t IntegerType) Or(other ValueType) (ValueType, error) {
	return logicalOr(t, other)
}
