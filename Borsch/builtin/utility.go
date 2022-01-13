package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func DeepCopy(object types.Type) (types.Type, error) {
	switch value := object.(type) {
	case *types.ClassInstance:
		copied := value.Copy()
		return copied, nil
	default:
		return value, nil
	}
}

func runOperator(name string, object types.Type, expectedTypeName string, expectedTypeHash uint64) (types.Type, error) {
	attribute, err := object.GetAttribute(name)
	if err != nil {
		return nil, util.RuntimeError(fmt.Sprintf("об'єкт типу '%s' не має довжини", object.GetTypeName()))
	}

	switch operator := attribute.(type) {
	case *types.FunctionInstance:
		args := []types.Type{object}
		kwargs := map[string]types.Type{operator.Arguments[0].Name: object}
		if err := types.CheckFunctionArguments(operator, &args, &kwargs); err != nil {
			return nil, err
		}

		result, err := operator.Call(&args, &kwargs)
		if err != nil {
			return nil, err
		}

		if err := types.CheckResult(result, operator); err != nil {
			return nil, err
		}

		if result.GetTypeHash() != expectedTypeHash {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"'%s' має повертати значення з типом '%s', отримано '%s'",
					name, expectedTypeName, result.GetTypeName(),
				),
			)
		}

		return result, nil
	default:
		return nil, util.ObjectIsNotCallable(name, attribute.GetTypeName())
	}
}
