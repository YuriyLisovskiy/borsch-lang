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
