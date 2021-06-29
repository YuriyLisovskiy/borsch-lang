package types

import (
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/lang/util"
	"strings"
	"unicode/utf8"
)

type StringType struct {
	Value string
}

func (t StringType) String() string {
	return t.Value
}

func (t StringType) Representation() string {
	return "\"" + t.String() + "\""
}

func (t StringType) TypeHash() int {
	return StringTypeHash
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

func (t StringType) Slice(from, to int64) (ValueType, error) {
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

func (t StringType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t StringType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t StringType) CompareTo(other ValueType) (int, error) {
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
			"%s", t.TypeName(), right.TypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t StringType) Add(other ValueType) (ValueType, error) {
	switch o := other.(type) {
	case StringType:
		return StringType{
			Value: t.Value + o.Value,
		}, nil
	default:
		return nil, nil
	}
}

func (t StringType) Sub(ValueType) (ValueType, error) {
	return nil, nil
}

func (t StringType) Mul(other ValueType) (ValueType, error) {
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

func (t StringType) Div(ValueType) (ValueType, error) {
	return nil, nil
}

func (t StringType) Pow(ValueType) (ValueType, error) {
	return nil, nil
}

func (t StringType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}
