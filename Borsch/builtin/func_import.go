package builtin

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalImport(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	return state.GetInterpreter().Import(state, (*args)[0].(types.StringInstance).Value)
}

func makeImportFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"імпорт",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "шлях",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalImport,
		[]types.FunctionReturnType{
			{
				Type:       types.Package,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
