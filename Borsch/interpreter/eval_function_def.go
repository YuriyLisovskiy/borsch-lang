package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func (node *FunctionDef) Evaluate(
	state common.State,
	parentPackage *types.Package,
	check func([]types.MethodParameter, []types.MethodReturnType) error,
) (types.Object, error) {
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

	function := types.MethodNew(
		node.Name.String(),
		parentPackage,
		arguments,
		returnTypes,
		func(ctx types.Context, _ types.Tuple, kwargs types.StringDict) (types.Object, error) {
			return node.Body.Evaluate(state.WithContext(ctx))
		},
	)
	return function, state.GetContext().SetVar(node.Name.String(), function)
}

func (node *ParametersSet) Evaluate(state common.State) ([]types.MethodParameter, error) {
	var arguments []types.MethodParameter
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

func (node *Parameter) Evaluate(ctx types.Context) (*types.MethodParameter, error) {
	class, err := ctx.GetClass(node.TypeName.String())
	if err != nil {
		return nil, err
	}

	return &types.MethodParameter{
		Class:      class.(*types.Class),
		Name:       node.Name.String(),
		IsVariadic: false,
		IsNullable: node.IsNullable,
	}, nil
}

func (node *FunctionBody) Evaluate(state common.State) (types.Object, error) {
	result := node.Stmts.Evaluate(state, true, false)
	return result.Value, result.Err
}

func (node *ReturnType) Evaluate(ctx types.Context) (*types.MethodReturnType, error) {
	class, err := ctx.GetClass(node.Name.String())
	if err != nil {
		return nil, err
	}

	return &types.MethodReturnType{
		Class:      class.(*types.Class),
		IsNullable: node.IsNullable,
	}, nil
}

func (node *ReturnStmt) Evaluate(state common.State) (types.Object, error) {
	resultCount := len(node.Expressions)
	switch {
	case resultCount == 1:
		return node.Expressions[0].Evaluate(state, nil)
	case resultCount > 1:
		// result := types.NewListInstance()
		// for _, expression := range node.Expressions {
		// 	value, err := expression.Evaluate(state, nil)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		//
		// 	result.Values = append(result.Values, value)
		// }
		//
		// return result, nil
	}

	panic("unreachable")
}
