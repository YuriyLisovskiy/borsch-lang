package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
)

func (f *FunctionDef) Evaluate(
	ctx common.Context,
	parentPackage *types.PackageInstance,
	check func([]types.FunctionArgument, []types.FunctionReturnType) error,
) (common.Type, error) {
	arguments, err := f.ParametersSet.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(ctx, f.ReturnTypes)
	if err != nil {
		return nil, err
	}

	if check != nil {
		if err := check(arguments, returnTypes); err != nil {
			return nil, err
		}
	}

	function := types.NewFunctionInstance(
		f.Name,
		arguments,
		func(ctx common.Context, _ *[]common.Type, kwargs *map[string]common.Type) (common.Type, error) {
			return f.Body.Evaluate(ctx)
		},
		returnTypes,
		parentPackage == nil,
		parentPackage,
		"", // TODO: add doc
	)
	return function, ctx.SetVar(f.Name, function)
}

func (p *ParametersSet) Evaluate(ctx common.Context) ([]types.FunctionArgument, error) {
	var arguments []types.FunctionArgument
	parameters := p.Parameters
	for _, parameter := range parameters {
		arg, err := parameter.Evaluate(ctx)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, *arg)
	}

	return arguments, nil
}

func (p *Parameter) Evaluate(ctx common.Context) (*types.FunctionArgument, error) {
	class, err := ctx.GetClass(p.Type)
	if err != nil {
		return nil, err
	}

	return &types.FunctionArgument{
		Type:       class.(*types.Class),
		Name:       p.Name,
		IsVariadic: false,
		IsNullable: p.IsNullable,
	}, nil
}

func (b *FunctionBody) Evaluate(ctx common.Context) (common.Type, error) {
	result := b.Stmts.Evaluate(ctx, true, false)
	return result.Value, result.Err
}

func (t *ReturnType) Evaluate(ctx common.Context) (*types.FunctionReturnType, error) {
	class, err := ctx.GetClass(t.Name)
	if err != nil {
		return nil, err
	}

	return &types.FunctionReturnType{
		Type:       class.(*types.Class),
		IsNullable: t.IsNullable,
	}, nil
}

func (s *ReturnStmt) Evaluate(ctx common.Context) (common.Type, error) {
	resultCount := len(s.Expressions)
	switch {
	case resultCount == 1:
		return s.Expressions[0].Evaluate(ctx, nil)
	case resultCount > 1:
		result := types.NewListInstance()
		for _, expression := range s.Expressions {
			value, err := expression.Evaluate(ctx, nil)
			if err != nil {
				return nil, err
			}

			result.Values = append(result.Values, value)
		}

		return result, nil
	}

	panic("unreachable")
}
