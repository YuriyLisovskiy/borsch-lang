package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func Length(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case types.SequentialType:
			return types.IntegerType{Value: arg.Length()}, nil
		case types.DictionaryType:
			return types.IntegerType{Value: arg.Length()}, nil
		}

		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"об'єкт типу '%s' не має довжини", args[0].TypeName(),
		))
	}

	return types.NoneType{}, util.RuntimeError("функція 'довжина()' приймає лише один аргумент")
}

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

func RemoveFromDictionary(args ...types.ValueType) (types.ValueType, error) {
	if len(args) != 2 {
		return nil, util.RuntimeError("функція 'вилучити()' приймає лише два аргументи")
	}

	switch container := args[0].(type) {
	case types.DictionaryType:
		err := container.RemoveElement(args[1])
		if err != nil {
			return nil, util.RuntimeError(err.Error())
		}

		return container, nil
	default:
		return nil, util.RuntimeError("першим аргументом має бути об'єкт з типом 'словник'")
	}
}
