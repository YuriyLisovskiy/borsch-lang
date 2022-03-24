package interpreter

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

var (
	BuiltinPackage = types.PackageNew(
		"builtin", nil, &ContextImpl{
			scopes:        []map[string]types.Object{},
			parentContext: nil,
		},
	)

	GlobalScope map[string]types.Object
)

func init() {
	GlobalScope = map[string]types.Object{
		"друкр": types.MethodNew(
			"друкр", BuiltinPackage, []types.MethodParameter{
				{
					Class:      types.AnyClass,
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
				return nil, nil
			},
		),
	}
}
