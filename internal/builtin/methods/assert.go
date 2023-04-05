package methods

import (
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func MakeAssert(pkg *types2.Package) *types2.Method {
	return types2.FunctionNew(
		"переконатися", pkg, []types2.MethodParameter{
			{
				Class:      types2.BoolClass,
				Name:       "умова",
				IsNullable: false,
				IsVariadic: false,
			},
			{
				Class:      types2.StringClass,
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
			if args[0].(types2.Bool) {
				return types2.Nil, nil
			}

			return nil, types2.NewAssertionError(string(args[1].(types2.String)))
		},
	)
}
