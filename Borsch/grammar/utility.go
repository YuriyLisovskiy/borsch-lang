package grammar

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustBool(result common.Type) (types.BoolInstance, error) {
	switch value := result.(type) {
	case types.BoolInstance:
		return value, nil
	default:
		var args []common.Type
		boolResult, err := callMethod(value, ops.BoolOperatorName, &args, nil)
		if err != nil {
			return types.BoolInstance{}, err
		}

		return boolResult.(types.BoolInstance), nil
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
