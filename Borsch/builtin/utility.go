package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func deepCopy(object common.Type) (common.Type, error) {
	switch value := object.(type) {
	case *types.ClassInstance:
		copied := value.Copy()
		return copied, nil
	default:
		return value, nil
	}
}

func runUnaryOperator(state common.State, name string, object common.Type, expectedType *types.Class) (
	common.Type,
	error,
) {
	var args []common.Type
	result, err := types.CallByName(state, object, name, &args, nil, true)
	if err != nil {
		return nil, err
	}

	if result.(types.ObjectInstance).GetPrototype() != expectedType {
		return nil, util.RuntimeError(
			fmt.Sprintf(
				"'%s' має повертати значення з типом '%s', отримано '%s'",
				name, expectedType.GetTypeName(),
				result.GetTypeName(),
			),
		)
	}

	return result, nil
}

func MustBool(value common.Type, errFunc func(common.Type) error) (bool, error) {
	switch integer := value.(type) {
	case types.BoolInstance:
		return integer.Value, nil
	default:
		return false, errFunc(value)
	}
}
