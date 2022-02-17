package builtin

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalType(_ common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
	return types.GetTypeOfInstance((*args)[0])
}

func makeTypeFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"тип",
		[]types.FunctionParameter{
			{
				Type:       types.AnyClass,
				Name:       "значення",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalType,
		[]types.FunctionReturnType{
			{
				Type:       types.AnyClass,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
