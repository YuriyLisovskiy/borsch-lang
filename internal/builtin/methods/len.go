package methods

import (
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func MakeLen(pkg *types2.Package) *types2.Method {
	return types2.FunctionNew(
		"довжина", pkg, []types2.MethodParameter{
			{
				Class:      types2.ObjectClass,
				Name:       "о",
				IsNullable: false,
				IsVariadic: false,
			},
		},
		[]types2.MethodReturnType{
			{
				Class:      types2.IntClass,
				IsNullable: false,
			},
		},
		func(ctx types2.Context, args types2.Tuple, kwargs types2.StringDict) (types2.Object, error) {
			arg0 := args[0]
			if seq, ok := arg0.(types2.ISequence); ok {
				return seq.Length(ctx)
			}

			return nil, types2.NewTypeErrorf("функція 'довжина' не підтримує об'єкти з типом %s", arg0.Class().Name)
		},
	)
}
