package builtin

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func DeepCopy(object common.Type) (common.Type, error) {
	switch value := object.(type) {
	case *types.ClassInstance:
		copied := value.Copy()
		return copied, nil
	default:
		return value, nil
	}
}

func runUnaryOperator(state common.State, name string, object common.Type, expectedType *types.Class) (
	common.Type,
	error,
) {
	var args []common.Type
	result, err := CallByName(state, object, name, &args, nil, true)
	if err != nil {
		return nil, err
	}

	if result.(types.ObjectInstance).GetPrototype() != expectedType {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'%s' має повертати значення з типом '%s', отримано '%s'",
				name, expectedType.GetTypeName(),
				result.GetTypeName(),
			),
		)
	}

	return result, nil
}

func Call(
	state common.State,
	function *types.FunctionInstance,
	args *[]common.Type,
	kwargs *map[string]common.Type,
) (common.Type, error) {
	if err := types.CheckFunctionArguments(state, function, args, nil); err != nil {
		return nil, err
	}

	if kwargs == nil {
		kwargs = &map[string]common.Type{}
	}

	argsLen := len(*args)
	for i := 0; i < argsLen; i++ {
		(*kwargs)[function.Arguments[i].Name] = (*args)[i]
	}

	ctx := function.GetContext()
	if ctx == nil {
		ctx = state.GetContext()
	}

	ctx = ctx.GetChild()
	ctx.PushScope(*kwargs)
	res, err := function.Call(state.WithContext(ctx), args, kwargs)
	if err != nil {
		return nil, util.RuntimeError(fmt.Sprintf(err.Error(), function.Name))
	}

	return res, nil
}

func CallByName(
	state common.State,
	object common.Type,
	funcName string,
	args *[]common.Type,
	kwargs *map[string]common.Type,
	isMethod bool,
) (common.Type, error) {
	attribute, err := object.GetAttribute(funcName)
	if err != nil {
		return nil, util.RuntimeError(err.Error())
	}

	switch function := attribute.(type) {
	case *types.FunctionInstance:
		if isMethod {
			if len(function.Arguments) == 0 {
				return nil, errors.New(fmt.Sprintf("%s is not a method", function.Representation(state)))
			}

			*args = append([]common.Type{object}, *args...)
		}

		return Call(state, function, args, kwargs)
	default:
		return nil, util.ObjectIsNotCallable(funcName, attribute.GetTypeName())
	}
}

func MustBool(value common.Type, errFunc func(common.Type) error) (bool, error) {
	switch integer := value.(type) {
	case types.BoolInstance:
		return integer.Value, nil
	default:
		return false, errFunc(value)
	}
}
