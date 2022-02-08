package builtin

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

func evalLength(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	sequence := (*args)[0]
	if !sequence.HasAttribute(common.LengthOperatorName) {
		return nil, errors.New(fmt.Sprintf("об'єкт типу '%s' не має довжини", sequence.GetTypeName()))
	}

	return runUnaryOperator(state, common.LengthOperatorName, sequence, types.Integer)
}

func makeLengthFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"довжина",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "послідовність",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		evalLength,
		[]types.FunctionReturnType{
			{
				Type:       types.Integer,
				IsNullable: false,
			},
		},
		false,
		types.BuiltinPackage,
		"", // TODO: add doc
	)
}
