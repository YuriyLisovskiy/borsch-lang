package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalHelp(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
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
				Type:       types.StringClass,
				Name:       "слово",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalHelp,
		[]types.FunctionReturnType{
			{
				Type:       types.NilClass,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
