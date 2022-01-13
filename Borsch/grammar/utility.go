package grammar

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func callMethod(object types.Type, funcName string, args *[]types.Type, kwargs *map[string]types.Type) (
	types.Type,
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

		*args = append([]types.Type{object}, *args...)
		if kwargs == nil {
			kwargs = &map[string]types.Type{}
		}

		argsLen := len(*args)
		for i := 0; i < argsLen; i++ {
			(*kwargs)[function.Arguments[i].Name] = (*args)[i]
		}

		if err := types.CheckFunctionArguments(function, args, kwargs); err != nil {
			return nil, err
		}

		res, err := function.Call(args, kwargs)
		if err != nil {
			return nil, util.RuntimeError(fmt.Sprintf(err.Error(), funcName))
		}

		return res, nil
	default:
		return nil, util.ObjectIsNotCallable(funcName, attribute.GetTypeName())
	}
}

func mustBool(result types.Type) (types.BoolInstance, error) {
	switch value := result.(type) {
	case types.BoolInstance:
		return value, nil
	default:
		var args []types.Type
		boolResult, err := callMethod(value, ops.BoolOperatorName, &args, nil)
		if err != nil {
			return types.BoolInstance{}, err
		}

		return boolResult.(types.BoolInstance), nil
	}
}

func evalBinaryOperator(
	ctx *Context,
	valueToSet types.Type,
	operatorName string,
	current OperatorEvaluatable,
	next OperatorEvaluatable,
) (types.Type, error) {
	left, err := current.Evaluate(ctx, valueToSet)
	if err != nil {
		return nil, err
	}

	if !reflect.ValueOf(next).IsNil() {
		right, err := next.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(left, operatorName, &[]types.Type{right}, nil)
	}

	return left, nil
}

func evalUnaryOperator(ctx *Context, operatorName string, operator OperatorEvaluatable) (types.Type, error) {
	if operator != nil {
		value, err := operator.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		return callMethod(value, operatorName, &[]types.Type{}, nil)
	}

	panic("unreachable")
}

func evalSingleGetByIndexOperation(variable types.Type, index types.Type) (types.Type, error) {
	switch iterable := variable.(type) {
	case types.SequentialType:
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
	ctx *Context,
	variable types.Type,
	indices []*Expression,
	value types.Type,
) (types.Type, error) {
	switch iterable := variable.(type) {
	case types.SequentialType:
		index, err := indices[0].Evaluate(ctx)
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
