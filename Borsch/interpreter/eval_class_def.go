package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

func (node *ClassDef) Evaluate(state types.State) (types.Object, error) {
	// TODO: add doc
	cls := types.ObjectClass.NewClass(node.Name.String(), "", nil, nil)

	// TODO: add doc
	// cls := &types.Class{
	// 	Name:        node.Name.String(),
	// 	IsFinal:     node.IsFinal,
	// 	ObjectClass: nil,
	// 	Parent:      state.GetCurrentPackage(),
	// }

	ctx := state.GetContext()
	for _, name := range node.Bases {
		base, err := ctx.GetClass(name.String())
		if err != nil {
			return nil, err
		}

		baseClass := base.(*types.Class)
		if baseClass.IsFinal {
			return nil, state.RuntimeError(
				fmt.Sprintf(
					"клас '%s' є закритим для розширення, не наслідуйте цей клас",
					name,
				),
				node,
			)
		}

		cls.Bases = append(cls.Bases, baseClass)
	}

	err := ctx.SetVar(node.Name.String(), cls)
	if err != nil {
		return nil, err
	}

	classContext := ctx.Derive()
	classContext.PushScope(map[string]types.Object{})
	for _, classMember := range node.Members {
		_, err := classMember.Evaluate(state.WithContext(classContext), cls)
		if err != nil {
			return nil, err
		}
	}

	cls.Dict = classContext.PopScope()

	// cls.SetAttributes(classContext.PopScope())
	// cls.Setup()
	// if !cls.IsValid() {
	// 	panic("custom class is invalid")
	// }

	return cls, nil
}

func (node *ClassMember) Evaluate(state types.State, class *types.Class) (types.Object, error) {
	if node.Variable != nil {
		return node.Variable.Evaluate(state)
	}

	if node.Method != nil {
		return node.Method.Evaluate(
			state,
			state.GetCurrentPackage().(*types.Package),
			// func(arguments []types.FunctionParameter, returnTypes []types.FunctionReturnType) error {
			// 	if err := checkMethod(class, arguments, returnTypes); err != nil {
			// 		return state.RuntimeError(err.Error(), node)
			// 	}
			//
			// 	if node.Method.Name == common.ConstructorName {
			// 		err := checkConstructor(arguments, returnTypes)
			// 		if err != nil {
			// 			return state.RuntimeError(err.Error(), node)
			// 		}
			// 	}
			//
			// 	return nil
			// },
		)
	}

	if node.Class != nil {
		return node.Class.Evaluate(state)
	}

	panic("unreachable")
}

// func checkMethod(class *types.Class, args []types.FunctionParameter, _ []types.FunctionReturnType) error {
// 	if len(args) == 0 {
// 		// TODO: ukr error text!
// 		return errors.New("not enough args, self required")
// 	}
//
// 	if args[0].Type != class {
// 		return errors.New(
// 			fmt.Sprintf(
// 				"перший параметер метода має бути типу '%s' отримано '%s'",
// 				class.GetTypeName(),
// 				args[0].GetTypeName(),
// 			),
// 		)
// 	}
//
// 	return nil
// }

// func checkConstructor(_ []types.FunctionParameter, returnTypes []types.FunctionReturnType) error {
// 	switch len(returnTypes) {
// 	case 0:
// 		// skip
// 	case 1:
// 		if returnTypes[0].Type != types.NilClass {
// 			return errors.New("конструктор має повертати 'нуль'")
// 		}
// 	default:
// 		return errors.New("ERROR")
// 	}
//
// 	return nil
// }
