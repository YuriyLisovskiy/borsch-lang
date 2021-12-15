package types

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

// BoolType TODO: move methods impl to attributes
type BoolType struct {
	Object

	Value    bool
	package_ *PackageType
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

	return BoolType{
		Value: boolean,
		Object: *newBuiltinObject(
			BoolTypeHash, map[string]Type{
				"__документ__": &NilType{}, // TODO: set doc
				"__пакет__":    BuiltinPackage,
			},
		),
		package_: BuiltinPackage,
	}, nil
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

func (t BoolType) AsBool() bool {
	return t.Value
}

func (t BoolType) SetAttribute(name string, _ Type) (Type, error) {
	if t.Object.HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t BoolType) Pow(other Type) (Type, error) {
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

func (t BoolType) Plus() (Type, error) {
	return IntegerType{Value: boolToInt64(t.Value)}, nil
}

func (t BoolType) Minus() (Type, error) {
	return IntegerType{Value: -boolToInt64(t.Value)}, nil
}

func (t BoolType) BitwiseNot() (Type, error) {
	return IntegerType{Value: ^boolToInt64(t.Value)}, nil
}

func (t BoolType) Mul(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: boolToInt64(t.Value) * boolToInt64(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: boolToInt64(t.Value) * o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: boolToFloat64(t.Value) * o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) Div(other Type) (Type, error) {
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

func (t BoolType) Mod(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		if o.Value {
			return IntegerType{
				Value: boolToInt64(t.Value) % boolToInt64(o.Value),
			}, nil
		}
	case IntegerType:
		if o.Value != 0 {
			return IntegerType{
				Value: boolToInt64(t.Value) % o.Value,
			}, nil
		}
	default:
		return nil, nil
	}

	return nil, errors.New("ділення за модулем на нуль")
}

func (t BoolType) Add(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: boolToInt64(t.Value) + boolToInt64(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: boolToInt64(t.Value) + o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: boolToFloat64(t.Value) + o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) Sub(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{
			Value: boolToInt64(t.Value) - boolToInt64(o.Value),
		}, nil
	case IntegerType:
		return IntegerType{
			Value: boolToInt64(t.Value) - o.Value,
		}, nil
	case RealType:
		return RealType{
			Value: boolToFloat64(t.Value) - o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) BitwiseLeftShift(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: boolToInt64(t.Value) << boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: boolToInt64(t.Value) << o.Value}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) BitwiseRightShift(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: boolToInt64(t.Value) >> boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: boolToInt64(t.Value) >> o.Value}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) BitwiseAnd(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: boolToInt64(t.Value) & boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: boolToInt64(t.Value) & o.Value}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) BitwiseXor(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: boolToInt64(t.Value) ^ boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: boolToInt64(t.Value) ^ o.Value}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) BitwiseOr(other Type) (Type, error) {
	switch o := other.(type) {
	case BoolType:
		return IntegerType{Value: boolToInt64(t.Value) | boolToInt64(o.Value)}, nil
	case IntegerType:
		return IntegerType{Value: boolToInt64(t.Value) | o.Value}, nil
	default:
		return nil, nil
	}
}

func (t BoolType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case BoolType:
		if t.Value == right.Value {
			return 0, nil
		}
	case IntegerType, RealType:
		if t.Value == right.AsBool() {
			return 0, nil
		}
	default:
		return 0, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				"%s", t.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t BoolType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t BoolType) And(other Type) (Type, error) {
	return logicalAnd(t, other)
}

func (t BoolType) Or(other Type) (Type, error) {
	return logicalOr(t, other)
}
