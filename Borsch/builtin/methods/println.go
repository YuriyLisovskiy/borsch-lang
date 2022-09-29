package methods

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

func MakePrintln(pkg *types.Package) *types.Method {
	return types.FunctionNew(
		"друкр", pkg, []types.MethodParameter{
			{
				Class:      types.ObjectClass,
				Name:       "повідомлення",
				IsNullable: false,
				IsVariadic: false,
			},
		},
		[]types.MethodReturnType{
			{
				Class:      types.NilClass,
				IsNullable: true,
			},
		},
		func(ctx types.Context, args types.Tuple, kwargs types.StringDict) (types.Object, error) {
			message, err := types.ToGoString(ctx, args[0])
			if err != nil {
				return nil, err
			}

			fmt.Println(message)
			return types.Nil, nil
		},
	)
}
