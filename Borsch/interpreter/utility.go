package interpreter

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalBinaryOperator(
	state common.State,
	valueToSet common.Value,
	operatorName string,
	current common.OperatorEvaluatable,
	next common.OperatorEvaluatable,
) (common.Value, error) {
	left, err := current.Evaluate(state, valueToSet)
	if err != nil {
		return nil, err
	}

	if !reflect.ValueOf(next).IsNil() {
		right, err := next.Evaluate(state, valueToSet)
		if err != nil {
			return nil, err
		}

		operator, err := left.GetOperator(operatorName)
		if err != nil {
			return nil, err
		}

		return types.CallAttribute(state, left, operator, operatorName, &[]common.Value{right}, nil, true)
	}

	return left, nil
}

func evalUnaryOperator(
	state common.State,
	operatorName string,
	operator common.OperatorEvaluatable,
) (common.Value, error) {
	if operator != nil {
		value, err := operator.Evaluate(state, nil)
		if err != nil {
			return nil, err
		}

		operatorFunc, err := value.GetOperator(operatorName)
		if err != nil {
			return nil, err
		}

		return types.CallAttribute(state, value, operatorFunc, operatorName, nil, nil, true)
	}

	panic("unreachable")
}

// evalSlicingOperation: "ranges_" len should be greater than 0
func evalSlicingOperation(
	state common.State,
	variable common.Value,
	ranges_ []*Range,
	valueToSet common.Value,
) (common.Value, error) {
	switch iterable := variable.(type) {
	case common.SequentialType:
		errMsg := ""
		if ranges_[0].IsSlicing {
			errMsg = "ліва межа має бути цілого типу"
		} else {
			errMsg = "індекс має бути цілого типу"
		}

		leftIdx, err := mustInt(
			state, ranges_[0].LeftBound, func(t common.Value) string {
				return fmt.Sprintf("%s, отримано %s", errMsg, t.GetTypeName())
			},
		)
		if err != nil {
			return nil, err
		}

		var element common.Value
		if ranges_[0].RightBound != nil {
			rightIdx, err := mustInt(
				state, ranges_[0].RightBound, func(t common.Value) string {
					return fmt.Sprintf("права межа має бути цілого типу, отримано %s", t.GetTypeName())
				},
			)
			if err != nil {
				return nil, err
			}

			element, err = iterable.Slice(state, leftIdx, rightIdx)
			if err != nil {
				return nil, err
			}

			if len(ranges_) == 1 {
				// valueToSet is ignored, return error maybe.
				return element, nil
			}
		} else if ranges_[0].IsSlicing {
			element, err = iterable.Slice(state, leftIdx, iterable.Length(state))
			if err != nil {
				return nil, err
			}

			if len(ranges_) == 1 {
				// valueToSet is ignored, return error maybe.
				return element, nil
			}
		} else {
			if len(ranges_) == 1 {
				if valueToSet != nil {
					return iterable.SetElement(state, leftIdx, valueToSet)
				}

				return iterable.GetElement(state, leftIdx)
			}

			element, err = iterable.GetElement(state, leftIdx)
			if err != nil {
				return nil, err
			}
		}

		element, err = evalSlicingOperation(state, element, ranges_[1:], valueToSet)
		if err != nil {
			return nil, err
		}

		return iterable.SetElement(state, leftIdx, element)
	default:
		operatorDescription := ""
		if ranges_[0].IsSlicing {
			operatorDescription = "зрізу"
		} else {
			operatorDescription = "довільного доступу"
		}

		return nil, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до об'єкта з типом '%s'",
				operatorDescription, variable.GetTypeName(),
			),
		)
	}
}

func mustInt(state common.State, expression *Expression, errFunc func(common.Value) string) (int64, error) {
	value, err := expression.Evaluate(state, nil)
	if err != nil {
		return 0, err
	}

	switch integer := value.(type) {
	case types.IntegerInstance:
		return integer.Value, nil
	default:
		return 0, errors.New(errFunc(value))
	}
}

func unpack(state common.State, lhs []*Expression, rhs []*Expression) (common.Value, error) {
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
	list := types.NewListInstance()
	for i = 0; i < lhsLen-1; i++ {
		element, err := lhs[i].Evaluate(state, sequence[i])
		if err != nil {
			return nil, err
		}

		list.Values = append(list.Values, element)
	}

	if i < len(sequence)-1 {
		rest := types.NewListInstance()
		rest.Values = sequence[i:]
		list.Values = append(list.Values, rest)
	} else {
		element, err := lhs[i].Evaluate(state, sequence[i])
		if err != nil {
			return nil, err
		}

		list.Values = append(list.Values, element)
	}

	return list, nil
}

func getSequenceOrResult(state common.State, lhs []*Expression, rhs []*Expression) (
	[]common.Value,
	common.Value,
	error,
) {
	rhsLen := len(rhs)
	var sequence []common.Value
	if rhsLen == 1 {
		element, err := rhs[0].Evaluate(state, nil)
		if err != nil {
			return nil, nil, err
		}

		switch list := element.(type) {
		case types.ListInstance:
			if len(lhs) == 1 {
				result, err := lhs[0].Evaluate(state, list)
				if err != nil {
					return nil, nil, err
				}

				return nil, result, nil
			}

			sequence = list.Values
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

func unpackList(state common.State, lhs []*Expression, rhs *Expression) (common.Value, error) {
	element, err := rhs.Evaluate(state, nil)
	if err != nil {
		return nil, err
	}

	switch list := element.(type) {
	case types.ListInstance:
		lhsLen := int64(len(lhs))
		rhsLen := list.Length(state)
		if lhsLen > rhsLen {
			// TODO: return error
			panic(fmt.Sprintf("unable to unpack %d elements of %s to %d vars", rhsLen, element.GetTypeName(), lhsLen))
		}

		var i int64
		resultList := types.NewListInstance()
		for i = 0; i < lhsLen-1; i++ {
			item, err := lhs[i].Evaluate(state, list.Values[i])
			if err != nil {
				return nil, err
			}

			resultList.Values = append(resultList.Values, item)
		}

		if i < list.Length(state)-1 {
			rest := types.NewListInstance()
			rest.Values = list.Values[i:]
			resultList.Values = append(resultList.Values, rest)
		} else {
			element, err := lhs[i].Evaluate(state, list.Values[i])
			if err != nil {
				return nil, err
			}

			resultList.Values = append(resultList.Values, element)
		}

		return resultList, nil
	}

	// TODO: return error
	panic(fmt.Sprintf("unable to unpack %s", element.GetTypeName()))
}

func evalReturnTypes(state common.State, returnTypes []*ReturnType) ([]types.FunctionReturnType, error) {
	var result []types.FunctionReturnType
	if len(returnTypes) == 0 {
		result = append(
			result, types.FunctionReturnType{
				Type:       types.Nil,
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

func getCurrentValue(ctx common.Context, prevValue common.Value, ident string) (common.Value, error) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return prevValue.GetAttribute(ident)
	}

	return ctx.GetVar(ident)
}

func setCurrentValue(ctx common.Context, prevValue common.Value, ident string, valueToSet common.Value) (
	common.Value,
	error,
) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return prevValue, prevValue.SetAttribute(ident, valueToSet)
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

func updateArgs(state common.State, arguments []*Expression, args *[]common.Value) error {
	for _, expressionArgument := range arguments {
		arg, err := expressionArgument.Evaluate(state, nil)
		if err != nil {
			return err
		}

		*args = append(*args, arg)
	}

	return nil
}
