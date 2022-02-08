package builtin

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalPanic(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	self := (*args)[0]
	msg, err := self.String(state)
	if err != nil {
		return nil, err
	}

	return types.NewNilInstance(), errors.New(fmt.Sprintf("%s: %s", self.GetTypeName(), msg))
}

func makePanicFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"панікувати",
		[]types.FunctionParameter{
			{
				Type:       ErrorClass,
				Name:       "помилка",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalPanic,
		[]types.FunctionReturnType{
			{
				Type:       types.Nil,
				IsNullable: true,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
