package methods

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

func MakeAssert(pkg *types.Package) *types.Method {
	return types.MethodNew(
		"переконатися", pkg, []types.MethodParameter{
			{
				Class:      types.BoolClass,
				Name:       "умова",
				IsNullable: false,
				IsVariadic: false,
			},
			{
				Class:      types.StringClass,
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
			if args[0].(types.Bool) {
				return nil, nil
			}

			return nil, types.NewAssertionError(string(args[1].(types.String)))
		},
	)
}
