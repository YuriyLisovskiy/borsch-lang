package builtin

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalAddToList(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	list := (*args)[0].(types.ListInstance)
	values := (*args)[1:]
	for _, value := range values {
		list.Values = append(list.Values, value)
	}

	return list, nil
}

func makeAddToListFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"додати",
		[]types.FunctionParameter{
			{
				Type:       types.List,
				Name:       "вхідний_список",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       types.Any,
				Name:       "елементи",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		evalAddToList,
		[]types.FunctionReturnType{
			{
				Type:       types.List,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
