package builtin

import (
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

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
