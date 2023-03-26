package methods

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"

func MakeLen(pkg *types.Package) *types.Method {
	return types.FunctionNew(
		"довжина", pkg, []types.MethodParameter{
			{
				Class:      types.ObjectClass,
				Name:       "о",
				IsNullable: false,
				IsVariadic: false,
			},
		},
		[]types.MethodReturnType{
			{
				Class:      types.IntClass,
				IsNullable: false,
			},
		},
		func(ctx types.Context, args types.Tuple, kwargs types.StringDict) (types.Object, error) {
			arg0 := args[0]
			if seq, ok := arg0.(types.ISequence); ok {
				return seq.Length(ctx)
			}

			return nil, types.NewTypeErrorf("функція 'довжина' не підтримує об'єкти з типом %s", arg0.Class().Name)
		},
	)
}
