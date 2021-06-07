package types

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type StringType struct {
	Value string
}

func (t StringType) String() string {
	return "StringType{\"" + t.Representation() + "\"}"
}

func (t StringType) Representation() string {
	return t.Value
}

func (t StringType) TypeHash() int {
	return stringType
}

func (t StringType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t StringType) Length() int64 {
	return int64(utf8.RuneCountInString(t.Value))
}

func (t StringType) GetElement(index int64) (ValueType, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	return StringType{Value: string([]rune(t.Value)[idx])}, nil
}

func (t StringType) SetElement(index int64, value ValueType) (ValueType, error) {
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
		return nil, errors.New(fmt.Sprintf("неможливо вставити в рядок об'єкт типу '%s'", v.TypeName()))
	}

	return t, nil
}
