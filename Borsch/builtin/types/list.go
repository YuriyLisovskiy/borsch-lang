package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ListInstance struct {
	Object
	Values []Type
}

func NewListInstance() ListInstance {
	return ListInstance{
		Values: []Type{},
		Object: Object{
			typeName:    GetTypeName(ListTypeHash),
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t ListInstance) String() string {
	return t.Representation()
}

func (t ListInstance) Representation() string {
	var strValues []string
	for _, v := range t.Values {
		strValues = append(strValues, v.Representation())
	}

	return "[" + strings.Join(strValues, ", ") + "]"
}

func (t ListInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t ListInstance) AsBool() bool {
	return t.Length() != 0
}

func (t ListInstance) SetAttribute(name string, _ Type) (Type, error) {
	if t.Object.HasAttribute(name) || t.GetClass().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t ListInstance) GetAttribute(name string) (Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (ListInstance) GetClass() *Class {
	return List
}

func (t ListInstance) Length() int64 {
	return int64(len(t.Values))
}

func (t ListInstance) GetElement(index int64) (Type, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	return t.Values[idx], nil
}

func (t ListInstance) SetElement(index int64, value Type) (Type, error) {
	idx, err := getIndex(index, t.Length())
	if err != nil {
		return nil, err
	}

	t.Values[idx] = value
	return t, nil
}

func (t ListInstance) Slice(from, to int64) (Type, error) {
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

	listInstance := NewListInstance()
	listInstance.Values = t.Values[fromIdx : toIdx+1]
	return listInstance, nil
}

func compareLists(self Type, other Type) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case ListInstance:
		return -2, util.RuntimeError(
			fmt.Sprintf(
				"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	default:
		return -2, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newListBinaryOperator(
	name string,
	doc string,
	handler func(ListInstance, Type) (Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name, ListTypeHash, AnyTypeHash, doc, func(left Type, right Type) (Type, error) {
			if leftInstance, ok := left.(ListInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newListClass() *Class {
	attributes := mergeAttributes(
		map[string]Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(ListTypeHash, ToList, ""),

			// TODO: add doc
			ops.LengthOperatorName: newLengthOperator(ListTypeHash, getLength, ""),

			ops.MulOp.Caption(): newListBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self ListInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case IntegerInstance:
						count := int(o.Value)
						list := NewListInstance()
						if count > 0 {
							for c := 0; c < count; c++ {
								list.Values = append(list.Values, self.Values...)
							}
						}

						return list, nil
					default:
						return nil, nil
					}
				},
			),
			ops.AddOp.Caption(): newListBinaryOperator(
				// TODO: add doc
				ops.AddOp.Caption(), "", func(self ListInstance, other Type) (Type, error) {
					switch o := other.(type) {
					case ListInstance:
						self.Values = append(self.Values, o.Values...)
						return self, nil
					default:
						return nil, nil
					}
				},
			),
		},
		makeLogicalOperators(ListTypeHash),
		makeComparisonOperators(ListTypeHash, compareLists),
	)
	return NewBuiltinClass(
		ListTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (Type, error) {
			return NewListInstance(), nil
		},
	)
}
