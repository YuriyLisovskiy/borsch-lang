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
	BuiltinObject
	Values []common.Type
}

func NewListInstance() ListInstance {
	return ListInstance{
		Values: []common.Type{},
		BuiltinObject: BuiltinObject{
			CommonObject{
				Object: Object{
					typeName:    common.ListTypeName,
					Attributes:  nil,
					callHandler: nil,
				},
				prototype: List,
			},
		},
	}
}

func (t ListInstance) String(state common.State) string {
	return t.Representation(state)
}

func (t ListInstance) Representation(state common.State) string {
	var strValues []string
	for _, value := range t.Values {
		strValues = append(strValues, value.Representation(state))
	}

	return "[" + strings.Join(strValues, ", ") + "]"
}

func (t ListInstance) AsBool(state common.State) bool {
	return t.Length(state) != 0
}

func (t ListInstance) Length(common.State) int64 {
	return int64(len(t.Values))
}

func (t ListInstance) GetElement(state common.State, index int64) (common.Type, error) {
	idx, err := getIndex(index, t.Length(state)-1)
	if err != nil {
		return nil, err
	}

	return t.Values[idx], nil
}

func (t ListInstance) SetElement(state common.State, index int64, value common.Type) (common.Type, error) {
	idx, err := getIndex(index, t.Length(state)-1)
	if err != nil {
		return nil, err
	}

	t.Values[idx] = value
	return t, nil
}

func (t ListInstance) Slice(state common.State, from, to int64) (common.Type, error) {
	fromIdx, err := getIndex(from, t.Length(state))
	if err != nil {
		return nil, err
	}

	toIdx, err := getIndex(to, t.Length(state))
	if err != nil {
		return nil, err
	}

	if fromIdx > toIdx {
		return nil, errors.New("індекс списку за межами послідовності")
	}

	listInstance := NewListInstance()
	listInstance.Values = t.Values[fromIdx:toIdx]
	return listInstance, nil
}

func compareLists(_ common.State, self common.Type, other common.Type) (int, error) {
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
		func(_ common.State, left common.Type, right common.Type) (common.Type, error) {
			if leftInstance, ok := left.(ListInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newListClass() *Class {
	initAttributes := func() map[string]common.Type {
		return mergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				ops.ConstructorName: newBuiltinConstructor(List, ToList, ""),

				// TODO: add doc
				ops.LengthOperatorName: newLengthOperator(List, getLength, ""),

				ops.MulOp.Name(): newListBinaryOperator(
					// TODO: add doc
					ops.MulOp.Name(), "", func(self ListInstance, other common.Type) (common.Type, error) {
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
				ops.AddOp.Name(): newListBinaryOperator(
					// TODO: add doc
					ops.AddOp.Name(), "", func(self ListInstance, other common.Type) (common.Type, error) {
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
	}

	return NewBuiltinClass(
		common.ListTypeName,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewListInstance(), nil
		},
	)
}
