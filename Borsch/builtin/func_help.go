package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalHelp(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	argStr, err := (*args)[0].String(state)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Поки що не паше =( - %s\n", argStr)
	return types.NewNilInstance(), nil
}

func makeHelpFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"допомога",
		[]types.FunctionParameter{
			{
				Type:       types.String,
				Name:       "слово",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalHelp,
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
