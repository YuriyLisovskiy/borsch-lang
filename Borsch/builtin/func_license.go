package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalLicense(common.State, *[]common.Object, *map[string]common.Object) (common.Object, error) {
	fmt.Println(build.License)
	return types.NewNilInstance(), nil
}

func makeLicenseFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"ліцензія",
		[]types.FunctionParameter{},
		evalLicense,
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
