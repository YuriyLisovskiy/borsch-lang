package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func (node *Call) Evaluate(
	state types.State,
	variable types.Object,
	selfInstance types.Object,
	isLambda *bool,
) (types.Object, error) {
	switch object := variable.(type) {
	case *types.Class:
		var args types.Tuple
		if err := updateArgs(state, node.Arguments, &args); err != nil {
			return nil, err
		}

		instance, err := object.New(object, args)
		if err != nil {
			return nil, err
		}

		err = object.Construct(instance, args)
		if err != nil {
			return nil, err
		}

		return instance, nil
	case *types.Method:
		*isLambda = object.Name == builtin.LambdaSignature
		var args types.Tuple
		if err := updateArgs(state, node.Arguments, &args); err != nil {
			return nil, err
		}

		if *isLambda || selfInstance == nil {
			return types.Call(state, object, args)
		}

		return object.Call(state, selfInstance, args)
	default:
		return nil, utilities.ObjectIsNotCallable(node.Ident.String(), object.Class().Name)
	}
}
