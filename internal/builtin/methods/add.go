package methods

import (
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
)

func MakeAdd(pkg *types2.Package) *types2.Method {
	return types2.FunctionNew(
		"додати", pkg, []types2.MethodParameter{
			{
				Classes: []*types2.Class{types2.ListClass},
				// Class:      types.ListClass,
				Name:       "список_",
				IsNullable: false,
				IsVariadic: false,
			},
			{
				Class:      types2.ObjectClass,
				Name:       "елемент",
				IsVariadic: false,
			},
		},
		[]types2.MethodReturnType{
			{
				Class:      types2.ListClass,
				IsNullable: false,
			},
		},
		func(ctx types2.Context, args types2.Tuple, kwargs types2.StringDict) (types2.Object, error) {
			arg0 := args[0]
			switch container := arg0.(type) {
			case *types2.List:
				return addToList(container, args[1]), nil
			}

			return nil, types2.NewTypeErrorf("функція 'додати' не підтримує об'єкти з типом %s", arg0.Class().Name)
		},
	)
}

func addToList(list *types2.List, item types2.Object) *types2.List {
	list.Values = append(list.Values, item)
	return list
}
