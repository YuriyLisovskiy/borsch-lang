package interpreter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalBinaryOperator(
	state types.State,
	valueToSet types.Object,
	current common.OperatorEvaluatable,
	next common.OperatorEvaluatable,
	operator func(left, right types.Object) (types.Object, error),
) (types.Object, error) {
	left, err := current.Evaluate(state, valueToSet)
	if err != nil {
		return nil, err
	}

	if !reflect.ValueOf(next).IsNil() {
		right, err := next.Evaluate(state, valueToSet)
		if err != nil {
			return nil, err
		}

		return operator(left, right)
	}

	return left, nil
}

func evalUnaryOperator(
	state types.State,
	target common.OperatorEvaluatable,
	operator func(self types.Object) (types.Object, error),
) (types.Object, error) {
	if target != nil {
		value, err := target.Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		return operator(value)
	}

	panic("unreachable")
}

// evalSlicingOperation: "ranges_" len should be greater than 0
func evalSlicingOperation(
	state types.State,
	variable types.Object,
	ranges_ []*Range,
	valueToSet types.Object,
) (types.Object, error) {
	errMsg := ""
	if ranges_[0].IsSlicing {
		errMsg = "ліва межа має бути цілого типу"
	} else {
		errMsg = "індекс має бути цілого типу"
	}

	leftIdx, err := mustInt(
		state, ranges_[0].LeftBound, func(t types.Object) string {
			return fmt.Sprintf("%s, отримано %s", errMsg, t.Class().Name)
		},
	)
	if err != nil {
		return nil, err
	}

	var element types.Object

	switch variable.(type) {
	case types.I__set_item__, types.I__get_item__:
		if len(ranges_) == 1 {
			if valueToSet != nil {
				return types.SetItem(variable, types.Int(leftIdx), valueToSet)
			}

			return types.GetItem(variable, types.Int(leftIdx))
		}

		element, err = types.GetItem(variable, types.Int(leftIdx))
		if err != nil {
			return nil, err
		}

		element, err = evalSlicingOperation(state, element, ranges_[1:], valueToSet)
		if err != nil {
			return nil, err
		}

		return types.SetItem(variable, types.Int(leftIdx), element)
	}

	operatorDescription := ""
	if ranges_[0].IsSlicing {
		operatorDescription = "зрізу"
	} else {
		operatorDescription = "довільного доступу"
	}

	return nil, errors.New(
		fmt.Sprintf(
			"неможливо застосувати оператор %s до об'єкта з типом '%s'",
			operatorDescription, variable.Class().Name,
		),
	)

	// switch iterable := variable.(type) {
	// case types.I__get_item__:
	// case types.I__set_item__:

	// case common.SequentialType:
	// errMsg := ""
	// if ranges_[0].IsSlicing {
	// 	errMsg = "ліва межа має бути цілого типу"
	// } else {
	// 	errMsg = "індекс має бути цілого типу"
	// }
	//
	// leftIdx, err := mustInt(
	// 	state, ranges_[0].LeftBound, func(t types.Object) string {
	// 		return fmt.Sprintf("%s, отримано %s", errMsg, t.Class().Name)
	// 	},
	// )
	// if err != nil {
	// 	return nil, err
	// }
	//
	// var element types.Object
	// if ranges_[0].RightBound != nil {
	// 	rightIdx, err := mustInt(
	// 		state, ranges_[0].RightBound, func(t types.Object) string {
	// 			return fmt.Sprintf("права межа має бути цілого типу, отримано %s", t.Class().Name)
	// 		},
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	element, err = iterable.Slice(state, leftIdx, rightIdx)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	if len(ranges_) == 1 {
	// 		// valueToSet is ignored, return error maybe.
	// 		return element, nil
	// 	}
	// } else if ranges_[0].IsSlicing {
	// 	element, err = iterable.Slice(state, leftIdx, iterable.Length(state))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	if len(ranges_) == 1 {
	// 		// valueToSet is ignored, return error maybe.
	// 		return element, nil
	// 	}
	// } else {
	// if len(ranges_) == 1 {
	// 	if valueToSet != nil {
	// 		return iterable.SetElement(state, leftIdx, valueToSet)
	// 	}
	//
	// 	return iterable.GetElement(state, leftIdx)
	// }
	//
	// element, err = iterable.GetElement(state, leftIdx)
	// if err != nil {
	// 	return nil, err
	// }
	// }

	// element, err = evalSlicingOperation(state, element, ranges_[1:], valueToSet)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return iterable.SetElement(state, leftIdx, element)
	// default:
	// 	operatorDescription := ""
	// 	if ranges_[0].IsSlicing {
	// 		operatorDescription = "зрізу"
	// 	} else {
	// 		operatorDescription = "довільного доступу"
	// 	}
	//
	// 	return nil, errors.New(
	// 		fmt.Sprintf(
	// 			"неможливо застосувати оператор %s до об'єкта з типом '%s'",
	// 			operatorDescription, variable.Class().Name,
	// 		),
	// 	)
	// }
}

func mustInt(state types.State, expression *Expression, errFunc func(types.Object) string) (int64, error) {
	value, err := expression.Evaluate(state, nil)
	if err != nil {
		return 0, err
	}

	switch integer := value.(type) {
	case types.Int:
		return integer.GoInt64()
	default:
		return 0, errors.New(errFunc(value))
	}
}

func unpack(state types.State, lhs []*Expression, rhs []*Expression) (types.Object, error) {
	lhsLen := len(lhs)
	rhsLen := len(rhs)
	if lhsLen > rhsLen {
		return unpackList(state, lhs, rhs[0])
	}

	sequence, result, err := getSequenceOrResult(state, lhs, rhs)
	if err != nil {
		return nil, err
	}

	if result != nil {
		return result, err
	}

	if lhsLen > len(sequence) {
		// TODO: return unable to unpack
		panic(fmt.Sprintf("unable to unpack %d elements to %d vars", len(sequence), lhsLen))
	}

	var i int
	list := types.NewListSized(lhsLen - 1)
	for i = 0; i < lhsLen-1; i++ {
		element, err := lhs[i].Evaluate(state, sequence[i])
		if err != nil {
			return nil, err
		}

		list.Items[i] = element
	}

	if i < len(sequence)-1 {
		rest := types.NewList()
		rest.Items = sequence[i:]
		list.Items = append(list.Items, rest)
	} else {
		element, err := lhs[i].Evaluate(state, sequence[i])
		if err != nil {
			return nil, err
		}

		list.Items = append(list.Items, element)
	}

	return list, nil
}

func getSequenceOrResult(state types.State, lhs []*Expression, rhs []*Expression) (
	[]types.Object,
	types.Object,
	error,
) {
	rhsLen := len(rhs)
	var sequence []types.Object
	if rhsLen == 1 {
		element, err := rhs[0].Evaluate(state, nil)
		if err != nil {
			return nil, nil, err
		}

		switch list := element.(type) {
		case *types.List:
			if len(lhs) == 1 {
				result, err := lhs[0].Evaluate(state, list)
				if err != nil {
					return nil, nil, err
				}

				return nil, result, nil
			}

			sequence = list.Items
		default:
			sequence = append(sequence, element)
		}
	} else {
		for _, expr := range rhs {
			element, err := expr.Evaluate(state, nil)
			if err != nil {
				return nil, nil, err
			}

			sequence = append(sequence, element)
		}
	}

	return sequence, nil, nil
}

func unpackList(state types.State, lhs []*Expression, rhs *Expression) (types.Object, error) {
	element, err := rhs.Evaluate(state, nil)
	if err != nil {
		return nil, err
	}

	switch list := element.(type) {
	case *types.List:
		lhsLen := len(lhs)
		rhsLen := list.Length()
		if lhsLen > rhsLen {
			// TODO: return error
			panic(fmt.Sprintf("unable to unpack %d elements of %s to %d vars", rhsLen, element.Class().Name, lhsLen))
		}

		var i int
		resultList := types.NewListSized(lhsLen - 1)
		for i = 0; i < lhsLen-1; i++ {
			item, err := lhs[i].Evaluate(state, list.Items[i])
			if err != nil {
				return nil, err
			}

			resultList.Items[i] = item
		}

		if i < list.Length()-1 {
			rest := types.NewList()
			rest.Items = list.Items[i:]
			resultList.Items = append(resultList.Items, rest)
		} else {
			element, err := lhs[i].Evaluate(state, list.Items[i])
			if err != nil {
				return nil, err
			}

			resultList.Items = append(resultList.Items, element)
		}

		return resultList, nil
	}

	// TODO: return error
	return nil, errors.New(fmt.Sprintf("unable to unpack %s", element.Class().Name))
}

/*
func evalReturnTypes(state common.State, returnTypes []*ReturnType) ([]types.FunctionReturnType, error) {
	var result []types.FunctionReturnType
	if len(returnTypes) == 0 {
		result = append(
			result, types.FunctionReturnType{
				Type:       types.NilClass,
				IsNullable: false,
			},
		)
	} else {
		for _, returnType := range returnTypes {
			r, err := returnType.Evaluate(state.GetContext())
			if err != nil {
				return nil, err
			}

			result = append(result, *r)
		}
	}

	return result, nil
}
*/

func getCurrentValue(ctx types.Context, prevValue types.Object, ident string) (types.Object, error) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return types.GetAttrString(prevValue, ident)
	}

	return ctx.GetVar(ident)
}

func setCurrentValue(ctx types.Context, prevValue types.Object, ident string, valueToSet types.Object) (
	types.Object,
	error,
) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return prevValue, types.SetAttrString(prevValue, ident, valueToSet)
	}

	return valueToSet, ctx.SetVar(ident, valueToSet)
}

func checkForNilAttribute(ident string) error {
	switch ident {
	case "нуль", "нульовий":
		return errors.New(fmt.Sprintf("'%s' не є атрибутом", ident))
	}

	return nil
}

func updateArgs(state types.State, arguments []*Expression, args *types.Tuple) error {
	for _, expressionArgument := range arguments {
		arg, err := expressionArgument.Evaluate(state, nil)
		if err != nil {
			return err
		}

		*args = append(*args, arg)
	}

	return nil
}
