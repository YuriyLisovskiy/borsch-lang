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
	Object
	Value string
}

func NewStringInstance(value string) StringInstance {
	return StringInstance{
		Value: value,
		Object: Object{
			typeName:    common.StringTypeName,
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t StringInstance) String(common.Context) string {
	return t.Value
}

func (t StringInstance) Representation(ctx common.Context) string {
	return "\"" + t.String(ctx) + "\""
}

func (t StringInstance) AsBool(ctx common.Context) bool {
	return t.Length(ctx) != 0
}

func (t StringInstance) GetTypeName() string {
	return t.GetPrototype().GetTypeName()
}

func (t StringInstance) SetAttribute(name string, _ common.Type) (common.Type, error) {
	if t.Object.HasAttribute(name) || t.GetPrototype().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t StringInstance) GetAttribute(name string) (common.Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetPrototype().GetAttribute(name)
}

func (StringInstance) GetPrototype() *Class {
	return String
}

func (t StringInstance) Length(_ common.Context) int64 {
	return int64(utf8.RuneCountInString(t.Value))
}

func (t StringInstance) GetElement(ctx common.Context, index int64) (common.Type, error) {
	idx, err := getIndex(index, t.Length(ctx))
	if err != nil {
		return nil, err
	}

	return NewStringInstance(string([]rune(t.Value)[idx])), nil
}

func (t StringInstance) SetElement(ctx common.Context, index int64, value common.Type) (common.Type, error) {
	switch v := value.(type) {
	case StringInstance:
		idx, err := getIndex(index, t.Length(ctx))
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

func (t StringInstance) Slice(ctx common.Context, from, to int64) (common.Type, error) {
	fromIdx, err := getIndex(from, t.Length(ctx))
	if err != nil {
		return nil, err
	}

	toIdx, err := getIndex(to, t.Length(ctx))
	if err != nil {
		return nil, err
	}

	if fromIdx > toIdx {
		return nil, errors.New("індекс рядка за межами послідовності")
	}

	return NewStringInstance(t.Value[fromIdx:toIdx]), nil
}

func compareStrings(_ common.Context, self, other common.Type) (int, error) {
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
		func(ctx common.Context, left common.Type, right common.Type) (common.Type, error) {
			if leftInstance, ok := left.(StringInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newStringClass() *Class {
	attributes := mergeAttributes(
		map[string]common.Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(String, ToString, ""),
			ops.MulOp.Caption(): newStringBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self StringInstance, other common.Type) (common.Type, error) {
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
			ops.AddOp.Caption(): newStringBinaryOperator(
				// TODO: add doc
				ops.AddOp.Caption(), "", func(self StringInstance, other common.Type) (common.Type, error) {
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
	return NewBuiltinClass(
		common.StringTypeName,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewStringInstance(""), nil
		},
	)
}
