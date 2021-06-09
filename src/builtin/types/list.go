package types

import (
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
	return listType
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
