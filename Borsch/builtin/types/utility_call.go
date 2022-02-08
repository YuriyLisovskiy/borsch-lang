package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func Call(
	state common.State,
	function *FunctionInstance,
	args *[]common.Value,
	kwargs *map[string]common.Value,
) (common.Value, error) {
	if args == nil {
		args = &[]common.Value{}
	}

	if err := CheckFunctionArguments(function, args, nil); err != nil {
		return nil, err
	}

	if kwargs == nil {
		kwargs = &map[string]common.Value{}
	}

	updateKwargs(*args, kwargs, function.Parameters)
	ctx := function.GetContext()
	if ctx == nil {
		ctx = state.GetContext()
	}

	ctx = ctx.GetChild()
	ctx.PushScope(*kwargs)
	funcState := state.WithContext(ctx)
	result, err := function.Call(funcState, args, kwargs)
	if err != nil {
		return nil, err
	}

	if err := CheckResult(funcState, result, function); err != nil {
		return nil, err
	}

	return result, nil
}

func CallAttribute(
	state common.State,
	object common.Value,
	attribute common.Value,
	attributeName string,
	args *[]common.Value,
	kwargs *map[string]common.Value,
	isMethod bool,
) (common.Value, error) {
	switch function := attribute.(type) {
	case *FunctionInstance:
		if isMethod {
			if len(function.Parameters) == 0 {
				functionStr, err := function.Representation(state)
				if err != nil {
					return nil, err
				}

				return nil, errors.New(fmt.Sprintf("%s is not a method", functionStr))
			}

			if args == nil {
				args = &[]common.Value{}
			}

			*args = append([]common.Value{object}, *args...)
		}

		return Call(state, function, args, kwargs)
	default:
		return nil, utilities.ObjectIsNotCallable(attributeName, attribute.GetTypeName())
	}
}

func CallByName(
	state common.State,
	object common.Value,
	funcName string,
	args *[]common.Value,
	kwargs *map[string]common.Value,
	isMethod bool,
) (common.Value, error) {
	attribute, err := object.GetAttribute(funcName)
	if err != nil {
		return nil, err
	}

	return CallAttribute(state, object, attribute, funcName, args, kwargs, isMethod)
}

func updateKwargs(args []common.Value, kwargs *map[string]common.Value, funcArgs []FunctionParameter) {
	argsLen := len(args)
	var i int
	for i = 0; i < argsLen && !funcArgs[i].IsVariadic; i++ {
		(*kwargs)[funcArgs[i].Name] = args[i]
	}

	if i < argsLen {
		list := NewListInstance()
		list.Values = args[i:]
		(*kwargs)[funcArgs[i].Name] = list
	}
}
