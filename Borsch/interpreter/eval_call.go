package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *Call) Evaluate(
	state common.State,
	variable types.Object,
	self types.Object,
	isLambda *bool,
) (types.Object, error) {
	args := types.Tuple{}
	switch object := variable.(type) {
	case *types.Method:
		*isLambda = object.Name == builtin.LambdaSignature
		if !*isLambda {
			switch o := self.(type) {
			case *types.Package:
				// ignore
			case *types.Class:
				if o.ClassType != nil {
					// got class instance
					args = append(args, self)
					// } else {
					// 	// got class
					// 	if err := updateArgs(state, node.Arguments, &args); err != nil {
					// 		return nil, err
					// 	}
					//
					// 	ctx := state.GetContext()
					// 	instance, err := o.New(ctx, o, args)
					// 	if err != nil {
					// 		return nil, err
					// 	}
					//
					// 	if o.Construct != nil {
					// 		err = o.Construct(ctx, instance, args)
					// 		if err != nil {
					// 			return nil, err
					// 		}
					// 	}
					//
					// 	return instance, nil
				}
			case nil:
				// ignore
			default:
				args = append(args, self)
			}
		}
	case *types.Class:
		if err := updateArgs(state, node.Arguments, &args); err != nil {
			return nil, err
		}

		ctx := state.GetContext()
		instance, err := object.New(ctx, object, args)
		if err != nil {
			return nil, err
		}

		if object.Construct != nil {
			err = object.Construct(ctx, instance, args)
			if err != nil {
				return nil, err
			}
		}

		return instance, nil
	}

	if err := updateArgs(state, node.Arguments, &args); err != nil {
		return nil, err
	}

	return types.Call(state.GetContext(), variable, args)
}
