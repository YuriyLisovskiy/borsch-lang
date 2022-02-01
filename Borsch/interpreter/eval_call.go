package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (a *Call) Evaluate(state common.State, variable common.Type, selfInstance common.Type) (common.Type, error) {
	switch object := variable.(type) {
	case *types.Class:
		var args []common.Type
		instance, err := object.GetEmptyInstance()
		if err != nil {
			return nil, err
		}

		_, err = a.evalFunctionByName(state, instance, ops.ConstructorName, &args, nil, true)
		if err != nil {
			return nil, err
		}

		return args[0], nil
	case *types.FunctionInstance:
		var args []common.Type
		if selfInstance != nil {
			switch selfInstance.(type) {
			case *types.Class, *types.PackageInstance:
				// ignore
			case types.ObjectInstance:
				args = append(args, selfInstance)
			}
		}

		return a.evalFunction(state, object, &args, nil)
	case types.ObjectInstance:
		args := []common.Type{variable}
		return a.evalFunctionByName(state, object.GetPrototype(), ops.CallOperatorName, &args, nil, true)
	default:
		return nil, util.ObjectIsNotCallable(a.Ident, object.GetTypeName())
	}
}

func (a *Call) evalFunctionByName(
	state common.State,
	object common.Type,
	functionName string,
	args *[]common.Type,
	kwargs *map[string]common.Type,
	isMethod bool,
) (common.Type, error) {
	if err := a.updateArgs(state, args); err != nil {
		return nil, err
	}

	return builtin.CallByName(state, object, functionName, args, kwargs, isMethod)
}

func (a *Call) evalFunction(
	state common.State,
	function *types.FunctionInstance,
	args *[]common.Type,
	kwargs *map[string]common.Type,
) (common.Type, error) {
	if err := a.updateArgs(state, args); err != nil {
		return nil, err
	}

	return builtin.Call(state, function, args, kwargs)
}

func (a *Call) updateArgs(state common.State, args *[]common.Type) error {
	for _, expressionArgument := range a.Arguments {
		arg, err := expressionArgument.Evaluate(state, nil)
		if err != nil {
			return err
		}

		*args = append(*args, arg)
	}

	return nil
}
