package builtin

import (
	"os"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalEnv(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	argStr, err := (*args)[0].String(state)
	if err != nil {
		return nil, err
	}

	return types.StringInstance{Value: os.Getenv(argStr)}, nil
}

func makeEnvFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"середовище",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "ключ",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalEnv,
		[]types.FunctionReturnType{
			{
				Type:       types.String,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
