package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

func (node *FunctionDef) Evaluate(
	state types.State,
	parentPackage *types.Package,
	// check func([]types.FunctionParameter, []types.FunctionReturnType) error,
) (types.Object, error) {
	// arguments, err := node.ParametersSet.Evaluate(state)
	// if err != nil {
	// 	return nil, err
	// }

	// returnTypes, err := evalReturnTypes(state, node.ReturnTypes)
	// if err != nil {
	// 	return nil, err
	// }

	// if check != nil {
	// 	if err := check(arguments, returnTypes); err != nil {
	// 		return nil, err
	// 	}
	// }

	function, err := types.NewMethod(
		node.Name.String(),
		func(state types.State, self types.Object, args types.Tuple) (types.Object, error) {
			return node.Body.Evaluate(state)
		},
		0,
		"", // TODO: add doc
	)
	if err != nil {
		return nil, err
	}

	function.Package = parentPackage
	return function, state.GetContext().SetVar(node.Name.String(), function)
}

// func (node *ParametersSet) Evaluate(state types.State) ([]types.FunctionParameter, error) {
// 	var arguments []types.FunctionParameter
// 	parameters := node.Parameters
// 	for _, parameter := range parameters {
// 		arg, err := parameter.Evaluate(state.GetContext())
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		arguments = append(arguments, *arg)
// 	}
//
// 	return arguments, nil
// }

// func (node *Parameter) Evaluate(ctx common.Context) (*types.FunctionParameter, error) {
// 	class, err := ctx.GetClass(node.TypeName.String())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &types.FunctionParameter{
// 		Type:       class.(*types.Class),
// 		Name:       node.Name.String(),
// 		IsVariadic: false,
// 		IsNullable: node.IsNullable,
// 	}, nil
// }

func (node *FunctionBody) Evaluate(state types.State) (types.Object, error) {
	result := node.Stmts.Evaluate(state, true, false)
	return result.Value, result.Err
}

// func (node *ReturnType) Evaluate(ctx common.Context) (*types.FunctionReturnType, error) {
// 	class, err := ctx.GetClass(node.Name.String())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &types.FunctionReturnType{
// 		Type:       class.(*types.Class),
// 		IsNullable: node.IsNullable,
// 	}, nil
// }

func (node *ReturnStmt) Evaluate(state types.State) (types.Object, error) {
	resultCount := len(node.Expressions)
	switch {
	case resultCount == 1:
		return node.Expressions[0].Evaluate(state, nil)
	case resultCount > 1:
		result := types.NewListSized(len(node.Expressions))
		for i, expression := range node.Expressions {
			value, err := expression.Evaluate(state, nil)
			if err != nil {
				return nil, err
			}

			result.Items[i] = value
		}

		return result, nil
	}

	panic("unreachable")
}
