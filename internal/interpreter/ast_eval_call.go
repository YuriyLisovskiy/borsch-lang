package interpreter

import (
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func (node *Call) Evaluate(state State, variable types2.Object) (types2.Object, error) {
	args := types2.Tuple{}
	if err := updateArgs(state, node.Arguments, &args); err != nil {
		return nil, err
	}

	return types2.Call(state.Context(), variable, args)
}
