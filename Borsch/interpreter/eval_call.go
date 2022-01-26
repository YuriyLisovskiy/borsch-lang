package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (a *Call) Evaluate(ctx common.Context, variable common.Type, selfInstance common.Type) (common.Type, error) {
	switch function := variable.(type) {
	case *types.Class:
		callable, err := function.GetAttribute(ops.ConstructorName)
		if err != nil {
			return nil, err
		}

		switch __constructor__ := callable.(type) {
		case *types.FunctionInstance:
			instance, err := function.GetEmptyInstance()
			if err != nil {
				return nil, err
			}

			args := []common.Type{instance}
			kwargs := map[string]common.Type{__constructor__.Arguments[0].Name: instance}

			// TODO: check if constructor returns nothing.
			_, err = a.evalFunction(ctx, __constructor__, &args, &kwargs, 1)
			if err != nil {
				return nil, err
			}

			return args[0], nil
		default:
			return nil, util.ObjectIsNotCallable(a.Ident, callable.GetTypeName())
		}
	case *types.FunctionInstance:
		var args []common.Type
		kwargs := map[string]common.Type{}
		argsShift := 0
		if selfInstance != nil {
			switch selfInstance.(type) {
			case *types.Class, *types.PackageInstance:
				// ignore
			case types.ObjectInstance:
				argsShift++
				args = append(args, selfInstance)
				kwargs[function.Arguments[0].Name] = selfInstance
			}
		}

		return a.evalFunction(ctx, function, &args, &kwargs, argsShift)
	case types.ObjectInstance:
		operator, err := function.GetPrototype().GetAttribute(ops.CallOperatorName)
		if err != nil {
			return nil, err
		}

		switch __call__ := operator.(type) {
		case *types.FunctionInstance:
			args := []common.Type{variable}
			kwargs := map[string]common.Type{__call__.Arguments[0].Name: variable}
			return a.evalFunction(ctx, __call__, &args, &kwargs, 1)
		default:
			return nil, util.ObjectIsNotCallable(a.Ident, operator.GetTypeName())
		}
	default:
		return nil, util.ObjectIsNotCallable(a.Ident, function.GetTypeName())
	}
}

func (a *Call) evalFunction(
	ctx common.Context,
	function *types.FunctionInstance,
	args *[]common.Type,
	kwargs *map[string]common.Type,
	argsShift int,
) (common.Type, error) {
	variadicArgs := types.NewListInstance()
	variadicArgsIndex := -1
	for i, expressionArgument := range a.Arguments {
		arg, err := expressionArgument.Evaluate(ctx, nil)
		if err != nil {
			return nil, err
		}

		*args = append(*args, arg)
		if variadicArgsIndex == -1 {
			if i+argsShift >= len(function.Arguments) {
				// TODO: return ukr error!
				return nil, util.RuntimeError("too many arguments")
			}

			if function.Arguments[i+argsShift].IsVariadic {
				variadicArgsIndex = i + argsShift
				variadicArgs.Values = append(variadicArgs.Values, arg)
			} else {
				(*kwargs)[function.Arguments[i+argsShift].Name] = arg
			}
		} else {
			variadicArgs.Values = append(variadicArgs.Values, arg)
		}
	}

	if variadicArgsIndex != -1 {
		(*kwargs)[function.Arguments[variadicArgsIndex].Name] = variadicArgs
	}

	if err := types.CheckFunctionArguments(ctx, function, args, kwargs); err != nil {
		return nil, err
	}

	ctx.PushScope(*kwargs)
	res, err := function.Call(ctx, args, kwargs)
	if err != nil {
		return nil, err
	}

	if err := types.CheckResult(ctx, res, function); err != nil {
		return nil, err
	}

	ctx.PopScope()
	return res, nil
}
