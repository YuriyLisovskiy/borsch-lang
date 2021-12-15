package types

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type StringType struct {
	Object

	Value string
	package_ *PackageType
}

func NewStringType(value string) StringType {
	return StringType{
		Value:    value,
		Object: *newBuiltinObject(
			StringTypeHash, map[string]Type{
				"__документ__": &NilType{}, // TODO: set doc
				"__пакет__":    BuiltinPackage,
			},
		),
		package_: BuiltinPackage,
	}
}

func (t StringType) String() string {
	return t.Value
}

func (t StringType) Representation() string {
	return "\"" + t.String() + "\""
}

func (t StringType) AsBool() bool {
	return t.Length() != 0
}

func (t StringType) Length() int64 {
	return int64(utf8.RuneCountInString(t.Value))
}

func (t StringType) GetElement(index int64) (Type, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	return StringType{Value: string([]rune(t.Value)[idx])}, nil
}

func (t StringType) SetElement(index int64, value Type) (Type, error) {
	switch v := value.(type) {
	case StringType:
		idx, err := getIndex(index, t.Length())
		if err != nil {
			return nil, err
		}

		if utf8.RuneCountInString(v.Value) != 1 {
			return nil, errors.New("неможливо вставити жодного, або більше ніж один символ в рядок")
		}

		runes := []rune(v.Value)
		target := []rune(t.Value)
		target[idx] = runes[0]
		t.Value = string(target)
	default:
		return nil, errors.New(fmt.Sprintf("неможливо вставити в рядок об'єкт типу '%s'", v.GetTypeName()))
	}

	return t, nil
}

func (t StringType) Slice(from, to int64) (Type, error) {
	fromIdx, err := getIndex(from, t.Length())
	if err != nil {
		return nil, err
	}

	toIdx, err := getIndex(to, t.Length())
	if err != nil {
		return nil, err
	}

	if fromIdx > toIdx {
		return nil, errors.New("індекс рядка за межами послідовності")
	}

	return StringType{Value: t.Value[fromIdx:toIdx]}, nil
}

func (t StringType) SetAttribute(name string, _ Type) (Type, error) {
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t StringType) Pow(Type) (Type, error) {
	return nil, nil
}

func (t StringType) Plus() (Type, error) {
	return nil, nil
}

func (t StringType) Minus() (Type, error) {
	return nil, nil
}

func (t StringType) BitwiseNot() (Type, error) {
	return nil, nil
}

func (t StringType) Mul(other Type) (Type, error) {
	switch o := other.(type) {
	case IntegerType:
		count := int(o.Value)
		if count < 0 {
			return StringType{Value: ""}, nil
		}

		return StringType{
			Value: strings.Repeat(t.Value, count),
		}, nil
	default:
		return nil, nil
	}
}

func (t StringType) Div(Type) (Type, error) {
	return nil, nil
}

func (t StringType) Mod(Type) (Type, error) {
	return nil, nil
}

func (t StringType) Add(other Type) (Type, error) {
	switch o := other.(type) {
	case StringType:
		return StringType{
			Value: t.Value + o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t StringType) Sub(Type) (Type, error) {
	return nil, nil
}

func (t StringType) BitwiseLeftShift(Type) (Type, error) {
	return nil, nil
}

func (t StringType) BitwiseRightShift(Type) (Type, error) {
	return nil, nil
}

func (t StringType) BitwiseAnd(Type) (Type, error) {
	return nil, nil
}

func (t StringType) BitwiseXor(Type) (Type, error) {
	return nil, nil
}

func (t StringType) BitwiseOr(Type) (Type, error) {
	return nil, nil
}

func (t StringType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case StringType:
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

func (t StringType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t StringType) And(other Type) (Type, error) {
	return logicalAnd(t, other)
}

func (t StringType) Or(other Type) (Type, error) {
	return logicalOr(t, other)
}
