package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalCopyright(common.State, *[]common.Value, *map[string]common.Value) (common.Value, error) {
	fmt.Printf("Copyright (c) %s %s.\nAll Rights Reserved.\n", build.Years, build.Author)
	return types.NewNilInstance(), nil
}

func makeCopyrightFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"автор",
		[]types.FunctionParameter{},
		evalCopyright,
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
