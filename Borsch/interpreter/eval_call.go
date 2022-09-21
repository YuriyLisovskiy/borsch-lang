package interpreter

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

// func (node *Call) Evaluate(
// 	state State,
// 	variable types.Object,
// 	self types.Object,
// 	isLambda *bool,
// ) (types.Object, error) {
// 	args := types.Tuple{}
// 	switch object := variable.(type) {
// 	case *types.Method:
// 		*isLambda = object.Name == builtin.LambdaSignature
// 		if !*isLambda {
// 			switch o := self.(type) {
// 			case *types.Package:
// 				// ignore
// 			case *types.Class:
// 				if o.ClassType != nil {
// 					// got class instance
// 					args = append(args, self)
// 				}
// 			case nil:
// 				// ignore
// 			default:
// 				args = append(args, self)
// 			}
// 		}
// 	case *types.Class:
// 		if err := updateArgs(state, node.Arguments, &args); err != nil {
// 			return nil, err
// 		}
//
// 		ctx := state.Context()
// 		instance, err := object.New(ctx, object, args)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		if object.Construct != nil {
// 			err = object.Construct(ctx, instance, args)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
//
// 		return instance, nil
// 	}
//
// 	if err := updateArgs(state, node.Arguments, &args); err != nil {
// 		return nil, err
// 	}
//
// 	return types.Call(state.Context(), variable, args)
// }

func (node *Call) Evaluate(state State, variable types.Object) (types.Object, error) {
	args := types.Tuple{}
	if cls, ok := variable.(*types.Class); ok {
		if err := updateArgs(state, node.Arguments, &args); err != nil {
			return nil, err
		}

		ctx := state.Context()
		instance, err := cls.New(ctx, cls, args)
		if err != nil {
			return nil, err
		}

		if cls.Construct != nil {
			err = cls.Construct(ctx, instance, args)
			if err != nil {
				return nil, err
			}
		}

		return instance, nil
	}

	if err := updateArgs(state, node.Arguments, &args); err != nil {
		return nil, err
	}

	return types.Call(state.Context(), variable, args)
}
