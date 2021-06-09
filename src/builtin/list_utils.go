package builtin

import (
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func AppendToList(args ...types.ValueType) (types.ValueType, error) {
	if len(args) < 2 {
		return nil, util.RuntimeError("функція 'додати()' приймає принаймні два аргументи")
	}

	switch list := args[0].(type) {
	case types.ListType:
		args = args[1:]
		for _, arg := range args {
			list.Values = append(list.Values, arg)
		}

		return list, nil
	default:
		return nil, util.RuntimeError("першим аргументом має бути об'єкт списку")
	}
}
