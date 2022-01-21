package grammar

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func (c *ClassDef) Evaluate(ctx common.Context, parentPackage *types.PackageInstance) (common.Type, error) {
	// TODO: add doc
	class := types.NewClass(c.Name, parentPackage, map[string]common.Type{}, "")
	err := ctx.SetVar(c.Name, class)
	if err != nil {
		return nil, err
	}

	classContext := ContextImpl{
		scopes:       []map[string]common.Type{{}},
		classContext: ctx,
	}

	for _, classMember := range c.Members {
		_, err := classMember.Evaluate(&classContext, class)
		if err != nil {
			return nil, err
		}
	}

	class.Attributes = classContext.PopScope()
	return class, nil
}

func (m *ClassMember) Evaluate(ctx common.Context, class *types.Class) (common.Type, error) {
	if m.Variable != nil {
		return m.Variable.Evaluate(ctx)
	}

	if m.Method != nil {
		return m.Method.Evaluate(
			ctx,
			nil,
			func(arguments []types.FunctionArgument, returnTypes []types.FunctionReturnType) error {
				if err := checkMethod(class, arguments, returnTypes); err != nil {
					return err
				}

				if m.Method.Name == ops.ConstructorName {
					return checkConstructor(arguments, returnTypes)
				}

				return nil
			},
		)
	}

	if m.Class != nil {
		return m.Class.Evaluate(ctx, nil)
	}

	panic("unreachable")
}

func checkMethod(class *types.Class, args []types.FunctionArgument, _ []types.FunctionReturnType) error {
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

func checkConstructor(_ []types.FunctionArgument, returnTypes []types.FunctionReturnType) error {
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
