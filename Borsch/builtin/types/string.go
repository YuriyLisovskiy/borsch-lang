package types

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type String struct {
	BuiltinInstance
	Value string
}

func NewStringInstance(value string) String {
	return String{
		BuiltinInstance: BuiltinInstance{
			ClassInstance: ClassInstance{
				class:      StringClass,
				attributes: map[string]common.Value{},
				address:    "",
			},
		},
		Value: value,
	}
}

func (t String) String(common.State) (string, error) {
	return t.Value, nil
}

func (t String) Representation(state common.State) (string, error) {
	value, err := t.String(state)
	if err != nil {
		return "", err
	}

	return "\"" + value + "\"", nil
}

func (t String) AsBool(state common.State) (bool, error) {
	return t.Length(state) != 0, nil
}

func (t String) Length(_ common.State) int64 {
	return int64(utf8.RuneCountInString(t.Value))
}

func (t String) GetElement(state common.State, index int64) (common.Value, error) {
	idx, err := getIndex(index, t.Length(state))
	if err != nil {
		return nil, err
	}

	return NewStringInstance(string([]rune(t.Value)[idx])), nil
}

func (t String) SetElement(state common.State, index int64, value common.Value) (common.Value, error) {
	switch v := value.(type) {
	case String:
		idx, err := getIndex(index, t.Length(state))
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

func (t String) Slice(state common.State, from, to int64) (common.Value, error) {
	length := t.Length(state)
	fromIdx := normalizeBound(from, length)
	toIdx := normalizeBound(to, length)
	if fromIdx > toIdx {
		return nil, errors.New("індекс рядка за межами послідовності")
	}

	return NewStringInstance(string([]rune(t.Value)[fromIdx:toIdx])), nil
}

func compareStrings(_ common.State, op common.Operator, self, other common.Value) (int, error) {
	left, ok := self.(String)
	if !ok {
		return 0, utilities.IncorrectUseOfFunctionError("compareStrings")
	}

	switch right := other.(type) {
	case NilInstance:
	case String:
		if left.Value == right.Value {
			return 0, nil
		}

		if left.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	default:
		return 0, utilities.OperatorNotSupportedError(op, left, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func stringBinaryOperator(
	operator common.Operator,
	handler func(common.State, String, common.Value) (common.Value, error),
) common.Value {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       StringClass,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "інший",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			left, ok := (*args)[0].(String)
			if !ok {
				return nil, utilities.InvalidUseOfOperator(operator, left, (*args)[1])
			}

			right := (*args)[1]
			result, err := handler(state, left, right)
			if err != nil {
				return nil, err
			}

			if result == nil {
				return nil, utilities.OperatorNotSupportedError(operator, left, right)
			}

			return result, nil
		},
		[]FunctionReturnType{
			{
				Type:       Any,
				IsNullable: false,
			},
		},
		true,
		nil,
		"", // TODO: add doc
	)
}

func stringOperator_Mul(_ common.State, left String, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case Int:
		count := int(other)
		if count < 0 {
			return NewStringInstance(""), nil
		}

		return NewStringInstance(strings.Repeat(left.Value, count)), nil
	default:
		return nil, nil
	}
}

func stringOperator_Add(_ common.State, left String, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case String:
		return NewStringInstance(left.Value + other.Value), nil
	default:
		return nil, nil
	}
}

func newStringClass() *Class {
	_ = func(attrs *map[string]common.Value) {
		*attrs = MergeAttributes(
			map[string]common.Value{
				// TODO: add doc
				common.ConstructorName: makeVariadicConstructor(StringClass, ToString, ""),
				common.MulOp.Name():    stringBinaryOperator(common.MulOp, stringOperator_Mul),
				common.AddOp.Name():    stringBinaryOperator(common.AddOp, stringOperator_Add),
			},
			MakeLogicalOperators(StringClass),
			MakeComparisonOperators(StringClass, compareStrings),
			MakeCommonOperators(StringClass),
		)
	}

	return &Class{
		Name:    common.StringTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		// Parent:          BuiltinPackage,
		// AttrInitializer: initAttributes,
		// GetEmptyInstance: func() (common.Value, error) {
		// 	return NewStringInstance(""), nil
		// },
	}
}
