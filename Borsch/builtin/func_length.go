package builtin

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalLength(state common.State, args *[]common.Object, _ *map[string]common.Object) (common.Object, error) {
	sequence := (*args)[0]
	if !sequence.HasAttribute(common.LengthOperatorName) {
		return nil, errors.New(fmt.Sprintf("об'єкт типу '%s' не має довжини", sequence.GetTypeName()))
	}

	return runUnaryOperator(state, common.LengthOperatorName, sequence, types.IntClass)
}

func makeLengthFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"довжина",
		[]types.FunctionParameter{
			{
				Type:       types.AnyClass,
				Name:       "послідовність",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalLength,
		[]types.FunctionReturnType{
			{
				Type:       types.IntClass,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
