package types

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
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
			typeName:    GetTypeName(StringTypeHash),
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t StringInstance) String() string {
	return t.Value
}

func (t StringInstance) Representation() string {
	return "\"" + t.String() + "\""
}

func (t StringInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t StringInstance) AsBool() bool {
	return t.Length() != 0
}

func (t StringInstance) SetAttribute(name string, _ Type) (Type, error) {
	if t.Object.HasAttribute(name) || t.GetClass().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t StringInstance) GetAttribute(name string) (Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (StringInstance) GetClass() *Class {
	return String
}

func (t StringInstance) Length() int64 {
	return int64(utf8.RuneCountInString(t.Value))
}

func (t StringInstance) GetElement(index int64) (Type, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	return NewStringInstance(string([]rune(t.Value)[idx])), nil
}

func (t StringInstance) SetElement(index int64, value Type) (Type, error) {
	switch v := value.(type) {
	case StringInstance:
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
		return nil, errors.New(fmt.Sprintf("неможливо вставити в рядок об'єкт типу '%s'", v.GetTypeName()))
	}

	return t, nil
}

func (t StringInstance) Slice(from, to int64) (Type, error) {
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

	return NewStringInstance(t.Value[fromIdx:toIdx]), nil
}

func compareStrings(self, other Type) (int, error) {
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
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
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
	handler func(StringInstance, Type) (Type, error),
) *FunctionInstance {
	return newBinaryOperator(
		name, StringTypeHash, AnyTypeHash, doc, func(left Type, right Type) (Type, error) {
			if leftInstance, ok := left.(StringInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newStringClass() *Class {
	attributes := mergeAttributes(
		map[string]Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(StringTypeHash, ToString, ""),
			ops.MulOp.Caption(): newStringBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self StringInstance, other Type) (Type, error) {
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
				ops.AddOp.Caption(), "", func(self StringInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case StringInstance:
						return NewStringInstance(self.Value + o.Value), nil
					default:
						return nil, nil
					}
				},
			),
		},
		makeLogicalOperators(StringTypeHash),
		makeComparisonOperators(StringTypeHash, compareStrings),
	)
	return NewBuiltinClass(
		StringTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (Type, error) {
			return NewStringInstance(""), nil
		},
	)
}
