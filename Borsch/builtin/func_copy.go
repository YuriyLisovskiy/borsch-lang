package builtin

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func deepCopy(object common.Value) (common.Value, error) {
	switch value := object.(type) {
	case *types.ClassInstance:
		copied := value.Copy()
		return copied, nil
	default:
		return value, nil
	}
}

func evalDeepCopy(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	return deepCopy((*args)[0])
}

func makeDeepCopyFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"копіювати",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "значення",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalDeepCopy,
		[]types.FunctionReturnType{
			{
				Type:       types.Any,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
