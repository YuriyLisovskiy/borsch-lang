package builtin

import (
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
)

func Exit(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 0 {
		os.Exit(0)
	} else if len(args) == 1 {
		switch code := args[0].(type) {
		case types.IntegerType:
			os.Exit(int(code.Value))
		default:
			return nil, util.RuntimeError("код виходу має бути цілим числом")
		}
	} else {
		return nil, util.RuntimeError("функція 'вихід()' приймає лише один аргумент")
	}

	return nil, nil
}
