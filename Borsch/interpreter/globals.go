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
		types.BoolClass.Name:   types.BoolClass,
		types.IntClass.Name:    types.IntClass,
		types.ListClass.Name:   types.ListClass,
		types.RealClass.Name:   types.RealClass,
		types.StringClass.Name: types.StringClass,
		types.TupleClass.Name:  types.TupleClass,

		types.ErrorClass.Name:                types.ErrorClass,
		types.RuntimeErrorClass.Name:         types.RuntimeErrorClass,
		types.TypeErrorClass.Name:            types.TypeErrorClass,
		types.AssertionErrorClass.Name:       types.AssertionErrorClass,
		types.ZeroDivisionErrorClass.Name:    types.ZeroDivisionErrorClass,
		types.IndexOutOfRangeErrorClass.Name: types.IndexOutOfRangeErrorClass,

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
		"переконатися": types.MethodNew(
			"переконатися", BuiltinPackage, []types.MethodParameter{
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
		),
	}
}
