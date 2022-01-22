package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ListInstance struct {
	Object
	Values []common.Type
}

func NewListInstance() ListInstance {
	return ListInstance{
		Values: []common.Type{},
		Object: Object{
			typeName:    common.ListTypeName,
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t ListInstance) String(ctx common.Context) string {
	return t.Representation(ctx)
}

func (t ListInstance) Representation(ctx common.Context) string {
	var strValues []string
	for _, value := range t.Values {
		strValues = append(strValues, value.Representation(ctx))
	}

	return "[" + strings.Join(strValues, ", ") + "]"
}

func (t ListInstance) AsBool(ctx common.Context) bool {
	return t.Length(ctx) != 0
}

func (t ListInstance) GetTypeName() string {
	return t.GetPrototype().GetTypeName()
}

func (t ListInstance) SetAttribute(name string, _ common.Type) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if t.Object.HasAttribute(name) || t.GetPrototype().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t ListInstance) GetAttribute(name string) (common.Type, error) {
	if name == ops.AttributesName {
		return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
	}

	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetPrototype().GetAttribute(name)
}

func (ListInstance) GetPrototype() *Class {
	return List
}

func (t ListInstance) Length(_ common.Context) int64 {
	return int64(len(t.Values))
}

func (t ListInstance) GetElement(ctx common.Context, index int64) (common.Type, error) {
	idx, err := getIndex(index, t.Length(ctx))
	if err != nil {
		return nil, err
	}

	return t.Values[idx], nil
}

func (t ListInstance) SetElement(ctx common.Context, index int64, value common.Type) (common.Type, error) {
	idx, err := getIndex(index, t.Length(ctx))
	if err != nil {
		return nil, err
	}

	t.Values[idx] = value
	return t, nil
}

func (t ListInstance) Slice(ctx common.Context, from, to int64) (common.Type, error) {
	fromIdx, err := getIndex(from, t.Length(ctx))
	if err != nil {
		return nil, err
	}

	toIdx, err := getIndex(to, t.Length(ctx))
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

func compareLists(_ common.Context, self common.Type, other common.Type) (int, error) {
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
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
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
	handler func(ListInstance, common.Type) (common.Type, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		List,
		Any,
		doc,
		func(ctx common.Context, left common.Type, right common.Type) (common.Type, error) {
			if leftInstance, ok := left.(ListInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newListClass() *Class {
	attributes := mergeAttributes(
		map[string]common.Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(List, ToList, ""),

			// TODO: add doc
			ops.LengthOperatorName: newLengthOperator(List, getLength, ""),

			ops.MulOp.Caption(): newListBinaryOperator(
				// TODO: add doc
				ops.MulOp.Caption(), "", func(self ListInstance, other common.Type) (common.Type, error) {
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
				ops.AddOp.Caption(), "", func(self ListInstance, other common.Type) (common.Type, error) {
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
		makeLogicalOperators(List),
		makeComparisonOperators(List, compareLists),
		makeCommonOperators(List),
	)
	return NewBuiltinClass(
		common.ListTypeName,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewListInstance(), nil
		},
	)
}
