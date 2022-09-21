package methods

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

func MakeAdd(pkg *types.Package) *types.Method {
	return types.FunctionNew(
		"додати", pkg, []types.MethodParameter{
			{
				Classes:    []*types.Class{types.ListClass},
				Name:       "список_",
				IsNullable: false,
				IsVariadic: false,
			},
			{
				Class:      types.AnyClass,
				Name:       "елемент",
				IsVariadic: false,
			},
		},
		[]types.MethodReturnType{
			{
				Class:      types.ListClass,
				IsNullable: false,
			},
		},
		func(ctx types.Context, args types.Tuple, kwargs types.StringDict) (types.Object, error) {
			arg0 := args[0]
			switch container := arg0.(type) {
			case *types.List:
				return addToList(container, args[1]), nil
			}

			return nil, types.NewTypeErrorf("функція 'додати' не підтримує об'єкти з типом %s", arg0.Class().Name)
		},
	)
}

func addToList(list *types.List, item types.Object) *types.List {
	list.Values = append(list.Values, item)
	return list
}
