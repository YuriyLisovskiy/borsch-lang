package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ListType struct {
	Values []ValueType
	object   *ObjectType
	package_ *PackageType
}

func NewListType() ListType {
	return ListType{
		Values: []ValueType{},
		object: newObjectType(
			ListTypeHash, map[string]ValueType{
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

func (t ListType) TypeHash() int {
	return ListTypeHash
}

func (t ListType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t ListType) AsBool() bool {
	return t.Length() != 0
}

func (t ListType) Length() int64 {
	return int64(len(t.Values))
}

func (t ListType) GetElement(index int64) (ValueType, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	return t.Values[idx], nil
}

func (t ListType) SetElement(index int64, value ValueType) (ValueType, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	t.Values[idx] = value
	return t, nil
}

func (t ListType) Slice(from, to int64) (ValueType, error) {
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

func (t ListType) GetAttr(name string) (ValueType, error) {
	return t.object.GetAttribute(name)
}

func (t ListType) SetAttr(name string, value ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t ListType) Pow(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) Plus() (ValueType, error) {
	return nil, nil
}

func (t ListType) Minus() (ValueType, error) {
	return nil, nil
}

func (t ListType) BitwiseNot() (ValueType, error) {
	return nil, nil
}

func (t ListType) Mul(other ValueType) (ValueType, error) {
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

func (t ListType) Div(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) Add(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case ListType:
		t.Values = append(t.Values, o.Values...)
		return t, nil
	default:
		return nil, nil
	}
}

func (t ListType) Sub(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) BitwiseLeftShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) BitwiseRightShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) BitwiseAnd(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) BitwiseXor(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) BitwiseOr(ValueType) (ValueType, error) {
	return nil, nil
}

func (t ListType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case ListType:
		return -2, util.RuntimeError(fmt.Sprintf(
			"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	default:
		return -2, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t ListType) Not() (ValueType, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t ListType) And(other ValueType) (ValueType, error) {
	return logicalAnd(t, other)
}

func (t ListType) Or(other ValueType) (ValueType, error) {
	return logicalOr(t, other)
}
