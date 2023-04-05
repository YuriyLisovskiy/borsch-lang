package interpreter

import (
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func (node *FunctionDef) Evaluate(
	state State,
	parentPackage *types2.Package,
	isClassMember bool,
	check func([]types2.MethodParameter, []types2.MethodReturnType) error,
) (types2.Object, error) {
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

	methodF := func(ctx types2.Context, _ types2.Tuple, kwargs types2.StringDict) (types2.Object, error) {
		return node.Body.Evaluate(state.NewChild().WithContext(ctx))
	}

	var method *types2.Method
	if isClassMember {
		method = types2.MethodNew(node.Name.String(), parentPackage, arguments, returnTypes, methodF)
	} else {
		method = types2.FunctionNew(node.Name.String(), parentPackage, arguments, returnTypes, methodF)
	}

	return method, state.Context().SetVar(node.Name.String(), method)
}

func (node *ParametersSet) Evaluate(state State) ([]types2.MethodParameter, error) {
	var arguments []types2.MethodParameter
	parameters := node.Parameters
	for _, parameter := range parameters {
		arg, err := parameter.Evaluate(state.Context())
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, *arg)
	}

	return arguments, nil
}

func (node *Parameter) Evaluate(ctx types2.Context) (*types2.MethodParameter, error) {
	class, err := ctx.GetClass(node.TypeName.String())
	if err != nil {
		return nil, err
	}

	return &types2.MethodParameter{
		Class:      class.(*types2.Class),
		Name:       node.Name.String(),
		IsVariadic: false,
		IsNullable: node.IsNullable,
	}, nil
}

func (node *FunctionBody) Evaluate(state State) (types2.Object, error) {
	result := node.Stmts.Evaluate(state, true, false)
	return result.Value, result.Err
}

func (node *ReturnType) Evaluate(ctx types2.Context) (*types2.MethodReturnType, error) {
	class, err := ctx.GetClass(node.Name.String())
	if err != nil {
		return nil, err
	}

	return &types2.MethodReturnType{
		Class:      class.(*types2.Class),
		IsNullable: node.IsNullable,
	}, nil
}

func (node *ReturnStmt) Evaluate(state State) (types2.Object, error) {
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
