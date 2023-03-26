package interpreter

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

func (node *Call) Evaluate(state State, variable types.Object) (types.Object, error) {
	args := types.Tuple{}
	if err := updateArgs(state, node.Arguments, &args); err != nil {
		return nil, err
	}

	return types.Call(state.Context(), variable, args)
}
