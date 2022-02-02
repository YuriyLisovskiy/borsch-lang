package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
)

func (f *FunctionDef) Evaluate(
	state common.State,
	parentPackage *types.PackageInstance,
	check func([]types.FunctionParameter, []types.FunctionReturnType) error,
) (common.Type, error) {
	arguments, err := f.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, f.ReturnTypes)
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
		func(state common.State, _ *[]common.Type, kwargs *map[string]common.Type) (common.Type, error) {
			return f.Body.Evaluate(state)
		},
		returnTypes,
		parentPackage == nil,
		parentPackage,
		"", // TODO: add doc
	)
	return function, state.GetContext().SetVar(f.Name, function)
}

func (p *ParametersSet) Evaluate(state common.State) ([]types.FunctionParameter, error) {
	var arguments []types.FunctionParameter
	parameters := p.Parameters
	for _, parameter := range parameters {
		arg, err := parameter.Evaluate(state.GetContext())
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, *arg)
	}

	return arguments, nil
}

func (p *Parameter) Evaluate(ctx common.Context) (*types.FunctionParameter, error) {
	class, err := ctx.GetClass(p.Type)
	if err != nil {
		return nil, err
	}

	return &types.FunctionParameter{
		Type:       class.(*types.Class),
		Name:       p.Name,
		IsVariadic: false,
		IsNullable: p.IsNullable,
	}, nil
}

func (b *FunctionBody) Evaluate(state common.State) (common.Type, error) {
	result := b.Stmts.Evaluate(state, true, false)
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

func (s *ReturnStmt) Evaluate(state common.State) (common.Type, error) {
	resultCount := len(s.Expressions)
	switch {
	case resultCount == 1:
		return s.Expressions[0].Evaluate(state, nil)
	case resultCount > 1:
		result := types.NewListInstance()
		for _, expression := range s.Expressions {
			value, err := expression.Evaluate(state, nil)
			if err != nil {
				return nil, err
			}

			result.Values = append(result.Values, value)
		}

		return result, nil
	}

	panic("unreachable")
}
