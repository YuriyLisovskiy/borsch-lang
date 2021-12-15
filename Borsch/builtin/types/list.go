package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ListType struct {
	Object

	Values []Type
	package_ *PackageType
}

func NewListType() ListType {
	return ListType{
		Values: []Type{},
		Object: *newBuiltinObject(
			ListTypeHash, map[string]Type{
				"__документ__": &NilType{}, // TODO: set doc
				"__пакет__":    BuiltinPackage,
			},
		),
		package_: BuiltinPackage,
	}
}

func (t ListType) String() string {
	return t.Representation()
}

func (t ListType) Representation() string {
	var strValues []string
	for _, v := range t.Values {
		strValues = append(strValues, v.Representation())
	}

	return "[" + strings.Join(strValues, ", ") + "]"
}

func (t ListType) AsBool() bool {
	return t.Length() != 0
}

func (t ListType) Length() int64 {
	return int64(len(t.Values))
}

func (t ListType) GetElement(index int64) (Type, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	return t.Values[idx], nil
}

func (t ListType) SetElement(index int64, value Type) (Type, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	t.Values[idx] = value
	return t, nil
}

func (t ListType) Slice(from, to int64) (Type, error) {
	fromIdx, err := getIndex(from, t.Length())
	if err != nil {
		return nil, err
	}

	toIdx, err := getIndex(to, t.Length())
	if err != nil {
		return nil, err
	}

	if fromIdx > toIdx {
		return nil, errors.New("індекс списку за межами послідовності")
	}

	return ListType{Values: t.Values[fromIdx:toIdx]}, nil
}

func (t ListType) SetAttribute(name string, _ Type) (Type, error) {
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t ListType) Pow(Type) (Type, error) {
	return nil, nil
}

func (t ListType) Plus() (Type, error) {
	return nil, nil
}

func (t ListType) Minus() (Type, error) {
	return nil, nil
}

func (t ListType) BitwiseNot() (Type, error) {
	return nil, nil
}

func (t ListType) Mul(other Type) (Type, error) {
	switch o := other.(type) {
	case IntegerType:
		count := int(o.Value)
		list := NewListType()
		if count > 0 {
			for c := 0; c < count; c++ {
				list.Values = append(list.Values, t.Values...)
			}
		}

		return list, nil
	default:
		return nil, nil
	}
}

func (t ListType) Div(Type) (Type, error) {
	return nil, nil
}

func (t ListType) Mod(Type) (Type, error) {
	return nil, nil
}

func (t ListType) Add(other Type) (Type, error) {
	switch o := other.(type) {
	case ListType:
		t.Values = append(t.Values, o.Values...)
		return t, nil
	default:
		return nil, nil
	}
}

func (t ListType) Sub(Type) (Type, error) {
	return nil, nil
}

func (t ListType) BitwiseLeftShift(Type) (Type, error) {
	return nil, nil
}

func (t ListType) BitwiseRightShift(Type) (Type, error) {
	return nil, nil
}

func (t ListType) BitwiseAnd(Type) (Type, error) {
	return nil, nil
}

func (t ListType) BitwiseXor(Type) (Type, error) {
	return nil, nil
}

func (t ListType) BitwiseOr(Type) (Type, error) {
	return nil, nil
}

func (t ListType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case ListType:
		return -2, util.RuntimeError(fmt.Sprintf(
			"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
			"%s", t.GetTypeName(), right.GetTypeName(),
		))
	default:
		return -2, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.GetTypeName(), right.GetTypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t ListType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t ListType) And(other Type) (Type, error) {
	return logicalAnd(t, other)
}

func (t ListType) Or(other Type) (Type, error) {
	return logicalOr(t, other)
}
