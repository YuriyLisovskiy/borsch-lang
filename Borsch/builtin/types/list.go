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
		BuiltinInstance: BuiltinInstance{
			ClassInstance{
				class:      List,
				attributes: map[string]common.Value{},
				address:    "",
			},
		},
		Values: []common.Value{},
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

func toList(_ common.State, args ...common.Value) (common.Value, error) {
	list := NewListInstance()
	if len(args) == 0 {
		return list, nil
	}

	for _, arg := range args {
		list.Values = append(list.Values, arg)
	}

	return list, nil
}

func compareLists(_ common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case ListInstance:
		return -2, util.OperandsNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	default:
		return -2, util.OperatorNotSupportedError(op, self, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func listOperator(
	operator common.Operator,
	handler func(common.State, ListInstance, common.Value) (common.Value, error),
) common.Value {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       List,
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
			left, ok := (*args)[0].(ListInstance)
			if !ok {
				return nil, util.InvalidUseOfOperator(operator, left, (*args)[1])
			}

			return handler(state, left, (*args)[1])
		},
		[]FunctionReturnType{
			{
				Type:       List,
				IsNullable: false,
			},
		},
		true,
		nil,
		"", // TODO: add doc
	)
}

func listOperator_Mul(_ common.State, left ListInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case IntegerInstance:
		count := int(other.Value)
		list := NewListInstance()
		if count > 0 {
			for c := 0; c < count; c++ {
				list.Values = append(list.Values, left.Values...)
			}
		}

		return list, nil
	default:
		return nil, nil
	}
}

func listOperator_Add(_ common.State, left ListInstance, right common.Value) (common.Value, error) {
	switch other := right.(type) {
	case ListInstance:
		left.Values = append(left.Values, other.Values...)
		return left, nil
	default:
		return nil, nil
	}
}

func newListClass() *Class {
	return &Class{
		Name:    common.ListTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		Parent:  BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Value) {
			*attrs = MergeAttributes(
				map[string]common.Value{
					// TODO: add doc
					common.ConstructorName: makeVariadicConstructor(List, toList, ""),

					// TODO: add doc
					common.LengthOperatorName: makeLengthOperator(List, ""),

					common.MulOp.Name(): listOperator(common.MulOp, listOperator_Mul),
					common.AddOp.Name(): listOperator(common.AddOp, listOperator_Add),
				},
				MakeLogicalOperators(List),
				MakeComparisonOperators(List, compareLists),
				MakeCommonOperators(List),
			)
		},
		GetEmptyInstance: func() (common.Value, error) {
			return NewListInstance(), nil
		},
	}
}
