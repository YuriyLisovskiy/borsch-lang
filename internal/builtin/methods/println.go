package methods

import (
	"fmt"

	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func MakePrintln(pkg *types2.Package) *types2.Method {
	return types2.FunctionNew(
		"друкр", pkg, []types2.MethodParameter{
			{
				Class:      types2.ObjectClass,
				Name:       "повідомлення",
				IsNullable: false,
				IsVariadic: false,
			},
		},
		[]types2.MethodReturnType{
			{
				Class:      types2.NilClass,
				IsNullable: true,
			},
		},
		func(ctx types2.Context, args types2.Tuple, kwargs types2.StringDict) (types2.Object, error) {
			message, err := types2.ToGoString(ctx, args[0])
			if err != nil {
				return nil, err
			}

			fmt.Println(message)
			return types2.Nil, nil
		},
	)
}
