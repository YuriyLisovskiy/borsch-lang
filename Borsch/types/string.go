package types

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type StringInstance struct {
	BuiltinInstance
	Value string
}

func NewStringInstance(value string) StringInstance {
	return StringInstance{
		Value: value,
		BuiltinInstance: BuiltinInstance{
			CommonInstance{
				Object: Object{
					typeName:    common.StringTypeName,
					Attributes:  nil,
					callHandler: nil,
				},
				prototype: String,
			},
		},
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

func (t StringInstance) AsBool(state common.State) bool {
	return t.Length(state) != 0
}

func (t StringInstance) Length(_ common.State) int64 {
	return int64(utf8.RuneCountInString(t.Value))
}

func (t StringInstance) GetElement(state common.State, index int64) (common.Type, error) {
	idx, err := getIndex(index, t.Length(state))
	if err != nil {
		return nil, err
	}

	return NewStringInstance(string([]rune(t.Value)[idx])), nil
}

func (t StringInstance) SetElement(state common.State, index int64, value common.Type) (common.Type, error) {
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

func (t StringInstance) Slice(state common.State, from, to int64) (common.Type, error) {
	length := t.Length(state)
	fromIdx := normalizeBound(from, length)
	toIdx := normalizeBound(to, length)
	if fromIdx > toIdx {
		return nil, errors.New("індекс рядка за межами послідовності")
	}

	return NewStringInstance(t.Value[fromIdx:toIdx]), nil
}

func compareStrings(_ common.State, self, other common.Type) (int, error) {
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
		return 0, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
				"%s", left.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newStringBinaryOperator(
	name string,
	doc string,
	handler func(StringInstance, common.Type) (common.Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		String,
		Any,
		doc,
		func(_ common.State, left common.Type, right common.Type) (common.Type, error) {
			if leftInstance, ok := left.(StringInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newStringClass() *Class {
	initAttributes := func() map[string]common.Type {
		return mergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				ops.ConstructorName: newBuiltinConstructor(String, ToString, ""),
				ops.MulOp.Name(): newStringBinaryOperator(
					// TODO: add doc
					ops.MulOp.Name(), "", func(self StringInstance, other common.Type) (common.Type, error) {
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
				ops.AddOp.Name(): newStringBinaryOperator(
					// TODO: add doc
					ops.AddOp.Name(), "", func(self StringInstance, other common.Type) (common.Type, error) {
						switch o := other.(type) {
						case StringInstance:
							return NewStringInstance(self.Value + o.Value), nil
						default:
							return nil, nil
						}
					},
				),
			},
			makeLogicalOperators(String),
			makeComparisonOperators(String, compareStrings),
			makeCommonOperators(String),
		)
	}

	return NewBuiltinClass(
		common.StringTypeName,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewStringInstance(""), nil
		},
	)
}
