package types

import (
	"errors"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type List []common.Object

func (t List) String(state common.State) (string, error) {
	return t.Representation(state)
}

func (t List) Representation(state common.State) (string, error) {
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

func (t List) AsBool(state common.State) (bool, error) {
	return t.Length(state) != 0, nil
}

func (t List) Length(common.State) int64 {
	return int64(len(t.Values))
}

func (t List) GetElement(state common.State, index int64) (common.Object, error) {
	idx, err := getIndex(index, t.Length(state))
	if err != nil {
		return nil, err
	}

	return t.Values[idx], nil
}

func (t List) SetElement(state common.State, index int64, value common.Object) (common.Object, error) {
	idx, err := getIndex(index, t.Length(state))
	if err != nil {
		return nil, err
	}

	t.Values[idx] = value
	return t, nil
}

func (t List) Slice(state common.State, from, to int64) (common.Object, error) {
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

func toList(_ common.State, args ...common.Object) (common.Object, error) {
	list := NewListInstance()
	if len(args) == 0 {
		return list, nil
	}

	for _, arg := range args {
		list.Values = append(list.Values, arg)
	}

	return list, nil
}

func compareLists(_ common.State, op common.Operator, self common.Object, other common.Object) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case List:
		return -2, utilities.OperandsNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	default:
		return -2, utilities.OperatorNotSupportedError(op, self, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func listOperator(
	operator common.Operator,
	handler func(common.State, List, common.Object) (common.Object, error),
) common.Object {
	return NewFunctionInstance(
		operator.Name(),
		[]FunctionParameter{
			{
				Type:       ListClass,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       AnyClass,
				Name:       "інший",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
			left, ok := (*args)[0].(List)
			if !ok {
				return nil, utilities.InvalidUseOfOperator(operator, left, (*args)[1])
			}

			return handler(state, left, (*args)[1])
		},
		[]FunctionReturnType{
			{
				Type:       ListClass,
				IsNullable: false,
			},
		},
		true,
		nil,
		"", // TODO: add doc
	)
}

func listOperator_Mul(_ common.State, left List, right common.Object) (common.Object, error) {
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

func listOperator_Add(_ common.State, left List, right common.Object) (common.Object, error) {
	switch other := right.(type) {
	case List:
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
		AttrInitializer: func(attrs *map[string]common.Object) {
			*attrs = MergeAttributes(
				map[string]common.Object{
					// TODO: add doc
					common.ConstructorName: makeVariadicConstructor(ListClass, toList, ""),

					// TODO: add doc
					common.LengthOperatorName: makeLengthOperator(ListClass, ""),

					common.MulOp.Name(): listOperator(common.MulOp, listOperator_Mul),
					common.AddOp.Name(): listOperator(common.AddOp, listOperator_Add),
				},
				MakeLogicalOperators(ListClass),
				MakeComparisonOperators(ListClass, compareLists),
				MakeCommonOperators(ListClass),
			)
		},
		GetEmptyInstance: func() (common.Object, error) {
			return NewListInstance(), nil
		},
	}
}
