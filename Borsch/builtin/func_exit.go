package builtin

import (
	"os"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalExit(_ common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	os.Exit(int((*args)[0].(types.IntegerInstance).Value))
	return types.NewNilInstance(), nil
}

func makeExitFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"вихід",
		[]types.FunctionParameter{
			{
				Type:       types.Integer,
				Name:       "код",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalExit,
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
