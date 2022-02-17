package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func (node *Call) Evaluate(
	state common.State,
	variable common.Object,
	selfInstance common.Object,
	isLambda *bool,
) (common.Object, error) {
	switch object := variable.(type) {
	case *types.Class:
		var args []common.Object
		instance, err := object.GetEmptyInstance()
		if err != nil {
			return nil, err
		}

		_, err = node.evalFunctionByName(state, instance, common.ConstructorName, &args, nil, true)
		if err != nil {
			return nil, err
		}

		return args[0], nil
	case *types.FunctionInstance:
		*isLambda = object.IsLambda()
		var args []common.Object
		if selfInstance != nil {
			switch selfInstance.(type) {
			case *types.Class, *types.PackageInstance:
				// ignore
			case types.ObjectInstance:
				if !*isLambda {
					args = append(args, selfInstance)
				}
			}
		}

		return node.evalFunction(state, object, &args, nil)
	case types.ObjectInstance:
		args := []common.Object{variable}
		return node.evalFunctionByName(state, object.GetClass(), common.CallOperatorName, &args, nil, true)
	default:
		return nil, utilities.ObjectIsNotCallable(node.Ident.String(), object.GetTypeName())
	}
}

func (node *Call) evalFunctionByName(
	state common.State,
	object common.Object,
	functionName string,
	args *[]common.Object,
	kwargs *map[string]common.Object,
	isMethod bool,
) (common.Object, error) {
	if err := updateArgs(state, node.Arguments, args); err != nil {
		return nil, err
	}

	return types.CallByName(state, object, functionName, args, kwargs, isMethod)
}

func (node *Call) evalFunction(
	state common.State,
	function *types.FunctionInstance,
	args *[]common.Object,
	kwargs *map[string]common.Object,
) (common.Object, error) {
	if err := updateArgs(state, node.Arguments, args); err != nil {
		return nil, err
	}

	return types.Call(state, function, args, kwargs)
}
