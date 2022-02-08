package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (node *ClassDef) Evaluate(state common.State) (common.Value, error) {
	ctx := state.GetContext()

	// TODO: add doc
	cls := &types.Class{
		Name:    node.Name,
		IsFinal: node.IsFinal,
		Class:   nil,
		Parent:  state.GetCurrentPackage(),
	}

	for _, name := range node.Bases {
		base, err := ctx.GetClass(name)
		if err != nil {
			return nil, err
		}

		baseClass := base.(*types.Class)
		if baseClass.IsFinalClass() {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"клас '%s' є закритим для розширення, не наслідуйте цей клас",
					name,
				),
			)
		}

		cls.Bases = append(cls.Bases, baseClass)
	}

	cls.GetEmptyInstance = func() (common.Value, error) {
		return types.NewClassInstance(cls, map[string]common.Value{}), nil
	}

	err := ctx.SetVar(node.Name, cls)
	if err != nil {
		return nil, err
	}

	classContext := ContextImpl{
		scopes:       []map[string]common.Value{{}},
		classContext: ctx,
	}

	for _, classMember := range node.Members {
		_, err := classMember.Evaluate(state.WithContext(&classContext), cls)
		if err != nil {
			return nil, err
		}
	}

	cls.SetAttributes(classContext.PopScope())
	cls.Setup()
	if !cls.IsValid() {
		panic("custom class is invalid")
	}

	return cls, nil
}

func (m *ClassMember) Evaluate(state common.State, class *types.Class) (common.Value, error) {
	if m.Variable != nil {
		return m.Variable.Evaluate(state)
	}

	if m.Method != nil {
		return m.Method.Evaluate(
			state,
			state.GetCurrentPackage().(*types.PackageInstance),
			func(arguments []types.FunctionParameter, returnTypes []types.FunctionReturnType) error {
				if err := checkMethod(class, arguments, returnTypes); err != nil {
					return err
				}

				if m.Method.Name == common.ConstructorName {
					return checkConstructor(arguments, returnTypes)
				}

				return nil
			},
		)
	}

	if m.Class != nil {
		return m.Class.Evaluate(state)
	}

	panic("unreachable")
}

func checkMethod(class *types.Class, args []types.FunctionParameter, _ []types.FunctionReturnType) error {
	if len(args) == 0 {
		// TODO: ukr error text!
		return util.RuntimeError("not enough args, self required")
	}

	if args[0].Type != class {
		return util.RuntimeError(
			fmt.Sprintf(
				"перший параметер метода має бути типу '%s' отримано '%s'",
				class.GetTypeName(),
				args[0].GetTypeName(),
			),
		)
	}

	return nil
}

func checkConstructor(_ []types.FunctionParameter, returnTypes []types.FunctionReturnType) error {
	switch len(returnTypes) {
	case 0:
		// skip
	case 1:
		if returnTypes[0].Type != types.Nil {
			return util.RuntimeError("конструктор має повертати 'нуль'")
		}
	default:
		return util.RuntimeError("ERROR")
	}

	return nil
}
