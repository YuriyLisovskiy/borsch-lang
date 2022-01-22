package grammar

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func callMethod(
	ctx common.Context,
	object common.Type,
	funcName string,
	args *[]common.Type,
	kwargs *map[string]common.Type,
) (common.Type, error) {
	attribute, err := object.GetAttribute(funcName)
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	switch function := attribute.(type) {
	case *types.FunctionInstance:
		if len(function.Arguments) == 0 {
			return nil, errors.New(fmt.Sprintf("%s is not a method", function.Representation(ctx)))
		}

		*args = append([]common.Type{object}, *args...)
		if kwargs == nil {
			kwargs = &map[string]common.Type{}
		}

		argsLen := len(*args)
		for i := 0; i < argsLen; i++ {
			(*kwargs)[function.Arguments[i].Name] = (*args)[i]
		}

		if err := types.CheckFunctionArguments(ctx, function, args, kwargs); err != nil {
			return nil, err
		}

		res, err := function.Call(ctx, args, kwargs)
		if err != nil {
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), funcName))
		}

		return res, nil
	default:
		return nil, util.ObjectIsNotCallable(funcName, attribute.GetTypeName())
	}
}

func evalBinaryOperator(
	ctx common.Context,
	valueToSet common.Type,
	operatorName string,
	current common.OperatorEvaluatable,
	next common.OperatorEvaluatable,
) (common.Type, error) {
	left, err := current.Evaluate(ctx, valueToSet)
	if err != nil {
		return nil, err
	}

	if !reflect.ValueOf(next).IsNil() {
		right, err := next.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(ctx, left, operatorName, &[]common.Type{right}, nil)
	}

	return left, nil
}

func evalUnaryOperator(
	ctx common.Context,
	operatorName string,
	operator common.OperatorEvaluatable,
) (common.Type, error) {
	if operator != nil {
		value, err := operator.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(ctx, value, operatorName, &[]common.Type{}, nil)
	}

	panic("unreachable")
}

// evalSlicingOperation: "ranges_" len should be greater than 0
func evalSlicingOperation(
	ctx common.Context,
	variable common.Type,
	ranges_ []*Range,
	valueToSet common.Type,
) (common.Type, error) {
	switch iterable := variable.(type) {
	case common.SequentialType:
		errMsg := ""
		if ranges_[0].IsSlicing {
			errMsg = "ліва межа має бути цілого типу"
		} else {
			errMsg = "індекс має бути цілого типу"
		}

		leftIdx, err := mustIntIndex(ctx, ranges_[0].LeftBound, errMsg)
		if err != nil {
			return nil, err
		}

		var element common.Type
		if ranges_[0].RightBound != nil {
			rightIdx, err := mustIntIndex(ctx, ranges_[0].RightBound, "права межа має бути цілого типу")
			if err != nil {
				return nil, err
			}

			element, err = iterable.Slice(ctx, leftIdx, rightIdx)
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
					return iterable.SetElement(ctx, leftIdx, valueToSet)
				}

				return iterable.GetElement(ctx, leftIdx)
			}

			element, err = iterable.GetElement(ctx, leftIdx)
			if err != nil {
				return nil, err
			}
		}

		return evalSlicingOperation(ctx, element, ranges_[1:], valueToSet)
	default:
		operatorDescription := ""
		if ranges_[0].IsSlicing {
			operatorDescription = "зрізу"
		} else {
			operatorDescription = "довільного доступу"
		}

		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до об'єкта з типом '%s'",
				operatorDescription, variable.GetTypeName(),
			),
		)
	}
}

func mustIntIndex(ctx common.Context, expression *Expression, errMessage string) (int64, error) {
	value, err := expression.Evaluate(ctx, nil)
	if err != nil {
		return 0, err
	}

	switch integer := value.(type) {
	case types.IntegerInstance:
		return integer.Value, nil
	default:
		return 0, util.RuntimeError(errMessage)
	}
}

func unpack(ctx common.Context, lhs []*Expression, rhs []*Expression) (common.Type, error) {
	lhsLen := len(lhs)
	rhsLen := len(rhs)
	if lhsLen > rhsLen {
		return unpackList(ctx, lhs, rhs[0])
	}

	sequence, result, err := getSequenceOrResult(ctx, lhs, rhs)
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
		element, err := lhs[i].Evaluate(ctx, sequence[i])
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
		element, err := lhs[i].Evaluate(ctx, sequence[i])
		if err != nil {
			return nil, err
		}

		list.Values = append(list.Values, element)
	}

	return list, nil
}

func getSequenceOrResult(ctx common.Context, lhs []*Expression, rhs []*Expression) (
	[]common.Type,
	common.Type,
	error,
) {
	rhsLen := len(rhs)
	var sequence []common.Type
	if rhsLen == 1 {
		element, err := rhs[0].Evaluate(ctx, nil)
		if err != nil {
			return nil, nil, err
		}

		switch list := element.(type) {
		case types.ListInstance:
			if len(lhs) == 1 {
				result, err := lhs[0].Evaluate(ctx, list)
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
			element, err := expr.Evaluate(ctx, nil)
			if err != nil {
				return nil, nil, err
			}

			sequence = append(sequence, element)
		}
	}

	return sequence, nil, nil
}

func unpackList(ctx common.Context, lhs []*Expression, rhs *Expression) (common.Type, error) {
	element, err := rhs.Evaluate(ctx, nil)
	if err != nil {
		return nil, err
	}

	switch list := element.(type) {
	case types.ListInstance:
		lhsLen := int64(len(lhs))
		rhsLen := list.Length(ctx)
		if lhsLen > rhsLen {
			// TODO: return error
			panic(fmt.Sprintf("unable to unpack %d elements of %s to %d vars", rhsLen, element.GetTypeName(), lhsLen))
		}

		var i int64
		resultList := types.NewListInstance()
		for i = 0; i < lhsLen-1; i++ {
			item, err := lhs[i].Evaluate(ctx, list.Values[i])
			if err != nil {
				return nil, err
			}

			resultList.Values = append(resultList.Values, item)
		}

		if i < list.Length(ctx)-1 {
			rest := types.NewListInstance()
			rest.Values = list.Values[i:]
			resultList.Values = append(resultList.Values, rest)
		} else {
			element, err := lhs[i].Evaluate(ctx, list.Values[i])
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

func evalReturnTypes(ctx common.Context, returnTypes []*ReturnType) ([]types.FunctionReturnType, error) {
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
			r, err := returnType.Evaluate(ctx)
			if err != nil {
				return nil, err
			}

			result = append(result, *r)
		}
	}

	return result, nil
}

func getCurrentValue(ctx common.Context, prevValue common.Type, ident string) (common.Type, error) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return prevValue.GetAttribute(ident)
	}

	return ctx.GetVar(ident)
}

func setCurrentValue(ctx common.Context, prevValue common.Type, ident string, valueToSet common.Type) (
	common.Type,
	error,
) {
	if prevValue != nil {
		if err := checkForNilAttribute(ident); err != nil {
			return nil, err
		}

		return prevValue.SetAttribute(ident, valueToSet)
	}

	return valueToSet, ctx.SetVar(ident, valueToSet)
}

func checkForNilAttribute(ident string) error {
	switch ident {
	case "нуль", "нульовий":
		return util.RuntimeError(fmt.Sprintf("'%s' не є атрибутом", ident))
	}

	return nil
}
