package grammar

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func callMethod(object common.Type, funcName string, args *[]common.Type, kwargs *map[string]common.Type) (
	common.Type,
	error,
) {
	attribute, err := object.GetAttribute(funcName)
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	switch function := attribute.(type) {
	case *types.FunctionInstance:
		if len(function.Arguments) == 0 {
			return nil, errors.New(fmt.Sprintf("%s is not a method", function.Representation()))
		}

		*args = append([]common.Type{object}, *args...)
		if kwargs == nil {
			kwargs = &map[string]common.Type{}
		}

		argsLen := len(*args)
		for i := 0; i < argsLen; i++ {
			(*kwargs)[function.Arguments[i].Name] = (*args)[i]
		}

		if err := types.CheckFunctionArguments(function, args, kwargs); err != nil {
			return nil, err
		}

		res, err := function.Call(nil, args, kwargs)
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

		return callMethod(left, operatorName, &[]common.Type{right}, nil)
	}

	return left, nil
}

func evalUnaryOperator(ctx common.Context, operatorName string, operator common.OperatorEvaluatable) (
	common.Type,
	error,
) {
	if operator != nil {
		value, err := operator.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(value, operatorName, &[]common.Type{}, nil)
	}

	panic("unreachable")
}

func evalSingleGetByIndexOperation(variable common.Type, index common.Type) (common.Type, error) {
	switch iterable := variable.(type) {
	case common.SequentialType:
		switch integerIndex := index.(type) {
		case types.IntegerInstance:
			return iterable.GetElement(integerIndex.Value)
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
				variable.GetTypeName(),
			),
		)
	}
}

func evalSingleSetByIndexOperation(
	ctx common.Context,
	variable common.Type,
	indices []*LogicalAnd,
	value common.Type,
) (common.Type, error) {
	switch iterable := variable.(type) {
	case common.SequentialType:
		index, err := indices[0].Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		switch integerIndex := index.(type) {
		case types.IntegerInstance:
			if len(indices) == 1 {
				return iterable.SetElement(integerIndex.Value, value)
			}

			element, err := iterable.GetElement(integerIndex.Value)
			if err != nil {
				return nil, err
			}

			element, err = evalSingleSetByIndexOperation(ctx, element, indices[1:], value)
			if err != nil {
				return nil, err
			}

			return iterable.SetElement(integerIndex.Value, element)
		default:
			return nil, util.RuntimeError("індекси мають бути цілого типу")
		}
	default:
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор довільного доступу до об'єкта з типом '%s'",
				variable.GetTypeName(),
			),
		)
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

func getSequenceOrResult(ctx common.Context, lhs []*Expression, rhs []*Expression) ([]common.Type, common.Type, error) {
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
		rhsLen := list.Length()
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

		if i < list.Length()-1 {
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

func evalParameters(ctx common.Context, parameters []*Parameter) []types.FunctionArgument {
	var arguments []types.FunctionArgument
	for _, parameter := range parameters {
		arguments = append(arguments, parameter.Evaluate(ctx))
	}

	return arguments
}

func evalReturnTypes(ctx common.Context, returnTypes []*ReturnType) []types.FunctionReturnType {
	var result []types.FunctionReturnType
	if len(returnTypes) == 0 {
		result = append(
			result, types.FunctionReturnType{
				TypeHash:   types.NilTypeHash,
				IsNullable: false,
			},
		)
	} else {
		for _, returnType := range returnTypes {
			result = append(result, returnType.Evaluate(ctx))
		}
	}

	return result
}
