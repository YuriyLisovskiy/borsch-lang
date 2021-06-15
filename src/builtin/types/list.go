package types

import (
	"errors"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"strings"
)

type ListType struct {
	Values []ValueType
}

func NewListType() ListType {
	return ListType{
		Values: []ValueType{},
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
	return nil, util.AttributeError(t.TypeName(), name)
}
