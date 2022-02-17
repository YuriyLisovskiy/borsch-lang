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
	args Tuple,
	kwargs *map[string]common.Object,
) (common.Object, error) {
	if args == nil {
		args = Tuple{}
	}

	if err := CheckFunctionArguments(function, args, nil); err != nil {
		return nil, err
	}

	if kwargs == nil {
		kwargs = &map[string]common.Object{}
	}

	updateKwargs(&args, kwargs, function.Parameters)
	ctx := function.GetContext()
	if ctx == nil {
		ctx = state.GetContext()
	}

	ctx.PushScope(*kwargs)
	funcState := state.WithContext(ctx)
	result, err := function.Call(funcState, args, kwargs)
	if err != nil {
		return nil, utilities.NewCallError(err, function.Name)
	}

	if err := CheckResult(funcState, result, function); err != nil {
		return nil, err
	}

	ctx.PopScope()
	return result, nil
}

func CallAttribute(
	state common.State,
	object common.Object,
	attribute common.Object,
	attributeName string,
	args *[]common.Object,
	kwargs *map[string]common.Object,
	isMethod bool,
) (common.Object, error) {
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
				args = &[]common.Object{}
			}

			*args = append([]common.Object{object}, *args...)
		}

		return Call(state, function, args, kwargs)
	default:
		return nil, utilities.ObjectIsNotCallable(attributeName, attribute.GetTypeName())
	}
}

func CallByName(
	state common.State,
	object common.Object,
	funcName string,
	args *[]common.Object,
	kwargs *map[string]common.Object,
	isMethod bool,
) (common.Object, error) {
	attribute, err := object.GetAttribute(funcName)
	if err != nil {
		return nil, err
	}

	return CallAttribute(state, object, attribute, funcName, args, kwargs, isMethod)
}

func updateKwargs(args []common.Object, kwargs *map[string]common.Object, funcArgs []FunctionParameter) {
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
