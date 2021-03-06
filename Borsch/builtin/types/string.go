package types

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type StringInstance struct {
	BuiltinInstance
	Value string
}

func NewStringInstance(value string) StringInstance {
	return StringInstance{
		BuiltinInstance: BuiltinInstance{
			ClassInstance: ClassInstance{
				class:      String,
				attributes: map[string]common.Value{},
				address:    "",
			},
		},
		Value: value,
	}
}

func (t StringInstance) String(common.State) (string, error) {
	return t.Value, nil
}

func (t StringInstance) Representation(state common.State) (string, error) {
	value, err := t.String(state)
	if err != nil {
		return "", err
	}

	return "\"" + value + "\"", nil
}

func (t StringInstance) AsBool(state common.State) (bool, error) {
	return t.Length(state) != 0, nil
}

func (t StringInstance) Length(_ common.State) int64 {
	return int64(utf8.RuneCountInString(t.Value))
}

func (t StringInstance) GetElement(state common.State, index int64) (common.Value, error) {
	idx, err := getIndex(index, t.Length(state))
	if err != nil {
		return nil, err
	}

	return NewStringInstance(string([]rune(t.Value)[idx])), nil
}

func (t StringInstance) SetElement(state common.State, index int64, value common.Value) (common.Value, error) {
	switch v := value.(type) {
	case StringInstance:
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

func (t StringInstance) Slice(state common.State, from, to int64) (common.Value, error) {
	length := t.Length(state)
	fromIdx := normalizeBound(from, length)
	toIdx := normalizeBound(to, length)
	if fromIdx > toIdx {
		return nil, errors.New("індекс рядка за межами послідовності")
	}

	return NewStringInstance(t.Value[fromIdx:toIdx]), nil
}

func compareStrings(_ common.State, op common.Operator, self, other common.Value) (int, error) {
	left, ok := self.(StringInstance)
	if !ok {
		return 0, util.IncorrectUseOfFunctionError("compareStrings")
	}

	switch right := other.(type) {
	case NilInstance:
	case StringInstance:
		if left.Value == right.Value {
			return 0, nil
		}

		if left.Value < right.Value {
			return -1, nil
		}

		return 1, nil
	default:
		return 0, util.OperatorNotSupportedError(op, left.GetTypeName(), right.GetTypeName())
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newStringBinaryOperator(
	name string,
	doc string,
	handler func(StringInstance, common.Value) (common.Value, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		String,
		Any,
		doc,
		func(_ common.State, left common.Value, right common.Value) (common.Value, error) {
			if leftInstance, ok := left.(StringInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newStringClass() *Class {
	initAttributes := func(attrs *map[string]common.Value) {
		*attrs = MergeAttributes(
			map[string]common.Value{
				// TODO: add doc
				common.ConstructorName: newBuiltinConstructor(String, ToString, ""),
				common.MulOp.Name(): newStringBinaryOperator(
					// TODO: add doc
					common.MulOp.Name(), "", func(self StringInstance, other common.Value) (common.Value, error) {
						switch o := other.(type) {
						case IntegerInstance:
							count := int(o.Value)
							if count < 0 {
								return NewStringInstance(""), nil
							}

							return NewStringInstance(strings.Repeat(self.Value, count)), nil
						default:
							return nil, nil
						}
					},
				),
				common.AddOp.Name(): newStringBinaryOperator(
					// TODO: add doc
					common.AddOp.Name(), "", func(self StringInstance, other common.Value) (common.Value, error) {
						switch o := other.(type) {
						case StringInstance:
							return NewStringInstance(self.Value + o.Value), nil
						default:
							return nil, nil
						}
					},
				),
			},
			MakeLogicalOperators(String),
			MakeComparisonOperators(String, compareStrings),
			MakeCommonOperators(String),
		)
	}

	return &Class{
		Name:            common.StringTypeName,
		IsFinal:         true,
		Bases:           []*Class{},
		Parent:          BuiltinPackage,
		AttrInitializer: initAttributes,
		GetEmptyInstance: func() (common.Value, error) {
			return NewStringInstance(""), nil
		},
	}
}
