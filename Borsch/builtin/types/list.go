package types

import (
	"errors"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type ListInstance struct {
	BuiltinInstance
	Values []common.Value
}

func NewListInstance() ListInstance {
	return ListInstance{
		Values: []common.Value{},
		BuiltinInstance: BuiltinInstance{
			CommonInstance{
				ObjectBase: ObjectBase{
					typeName:    common.ListTypeName,
					Attributes:  nil,
					callHandler: nil,
				},
				prototype: List,
			},
		},
	}
}

func (t ListInstance) String(state common.State) (string, error) {
	return t.Representation(state)
}

func (t ListInstance) Representation(state common.State) (string, error) {
	var strValues []string
	for _, value := range t.Values {
		strValue, err := value.Representation(state)
		if err != nil {
			return "", err
		}

		strValues = append(strValues, strValue)
	}

	return "[" + strings.Join(strValues, ", ") + "]", nil
}

func (t ListInstance) AsBool(state common.State) (bool, error) {
	return t.Length(state) != 0, nil
}

func (t ListInstance) Length(common.State) int64 {
	return int64(len(t.Values))
}

func (t ListInstance) GetElement(state common.State, index int64) (common.Value, error) {
	idx, err := getIndex(index, t.Length(state))
	if err != nil {
		return nil, err
	}

	return t.Values[idx], nil
}

func (t ListInstance) SetElement(state common.State, index int64, value common.Value) (common.Value, error) {
	idx, err := getIndex(index, t.Length(state))
	if err != nil {
		return nil, err
	}

	t.Values[idx] = value
	return t, nil
}

func (t ListInstance) Slice(state common.State, from, to int64) (common.Value, error) {
	length := t.Length(state)
	fromIdx := normalizeBound(from, length)
	toIdx := normalizeBound(to, length)
	if fromIdx > toIdx {
		return nil, errors.New("індекс списку за межами послідовності")
	}

	listInstance := NewListInstance()
	listInstance.Values = t.Values[fromIdx:toIdx]
	return listInstance, nil
}

func compareLists(_ common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case ListInstance:
		return -2, util.OperandsNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	default:
		return -2, util.OperatorNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newListBinaryOperator(
	name string,
	doc string,
	handler func(ListInstance, common.Value) (common.Value, error),
) *FunctionInstance {
	return newBinaryMethod(
		name,
		List,
		Any,
		doc,
		func(_ common.State, left common.Value, right common.Value) (common.Value, error) {
			if leftInstance, ok := left.(ListInstance); ok {
				return handler(leftInstance, right)
			}

			return nil, util.IncorrectUseOfFunctionError(name)
		},
	)
}

func newListClass() *Class {
	initAttributes := func() map[string]common.Value {
		return MergeAttributes(
			map[string]common.Value{
				// TODO: add doc
				common.ConstructorName: newBuiltinConstructor(List, ToList, ""),

				// TODO: add doc
				common.LengthOperatorName: newLengthOperator(List, getLength, ""),

				common.MulOp.Name(): newListBinaryOperator(
					// TODO: add doc
					common.MulOp.Name(), "", func(self ListInstance, other common.Value) (common.Value, error) {
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
				common.AddOp.Name(): newListBinaryOperator(
					// TODO: add doc
					common.AddOp.Name(), "", func(self ListInstance, other common.Value) (common.Value, error) {
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
			MakeLogicalOperators(List),
			MakeComparisonOperators(List, compareLists),
			MakeCommonOperators(List),
		)
	}

	return NewBuiltinClass(
		common.ListTypeName,
		nil,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Value, error) {
			return NewListInstance(), nil
		},
	)
}
