package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *FunctionDef) Evaluate(
	state common.State,
	parentPackage *types.PackageInstance,
	check func([]types.FunctionParameter, []types.FunctionReturnType) error,
) (common.Value, error) {
	arguments, err := node.ParametersSet.Evaluate(state)
	if err != nil {
		return nil, err
	}

	returnTypes, err := evalReturnTypes(state, node.ReturnTypes)
	if err != nil {
		return nil, err
	}

	if check != nil {
		if err := check(arguments, returnTypes); err != nil {
			return nil, err
		}
	}

	function := types.NewFunctionInstance(
		node.Name,
		arguments,
		func(state common.State, _ *[]common.Value, kwargs *map[string]common.Value) (common.Value, error) {
			return node.Body.Evaluate(state)
		},
		returnTypes,
		parentPackage == nil,
		parentPackage,
		"", // TODO: add doc
	)
	return function, state.GetContext().SetVar(node.Name, function)
}

func (node *ParametersSet) Evaluate(state common.State) ([]types.FunctionParameter, error) {
	var arguments []types.FunctionParameter
	parameters := node.Parameters
	for _, parameter := range parameters {
		arg, err := parameter.Evaluate(state.GetContext())
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, *arg)
	}

	return arguments, nil
}

func (node *Parameter) Evaluate(ctx common.Context) (*types.FunctionParameter, error) {
	class, err := ctx.GetClass(node.Type)
	if err != nil {
		return nil, err
	}

	return &types.FunctionParameter{
		Type:       class.(*types.Class),
		Name:       node.Name,
		IsVariadic: false,
		IsNullable: node.IsNullable,
	}, nil
}

func (b *FunctionBody) Evaluate(state common.State) (common.Value, error) {
	result := b.Stmts.Evaluate(state, true, false)
	return result.Value, result.Err
}

func (node *ReturnType) Evaluate(ctx common.Context) (*types.FunctionReturnType, error) {
	class, err := ctx.GetClass(node.Name)
	if err != nil {
		return nil, err
	}

	return &types.FunctionReturnType{
		Type:       class.(*types.Class),
		IsNullable: node.IsNullable,
	}, nil
}

func (node *ReturnStmt) Evaluate(state common.State) (common.Value, error) {
	resultCount := len(node.Expressions)
	switch {
	case resultCount == 1:
		return node.Expressions[0].Evaluate(state, nil)
	case resultCount > 1:
		result := types.NewListInstance()
		for _, expression := range node.Expressions {
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
