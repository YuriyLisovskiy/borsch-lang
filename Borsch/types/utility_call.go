package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func Call(
	state common.State,
	function *FunctionInstance,
	args *[]common.Type,
	kwargs *map[string]common.Type,
) (common.Type, error) {
	if args == nil {
		args = &[]common.Type{}
	}

	if err := CheckFunctionArguments(state, function, args, nil); err != nil {
		return nil, err
	}

	if kwargs == nil {
		kwargs = &map[string]common.Type{}
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
		return nil, err
	}

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
				args = &[]common.Type{}
			}

			*args = append([]common.Type{object}, *args...)
		}

		return Call(state, function, args, kwargs)
	default:
		return nil, util.ObjectIsNotCallable(funcName, attribute.GetTypeName())
	}
}

func updateKwargs(args []common.Type, kwargs *map[string]common.Type, funcArgs []FunctionParameter) {
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
