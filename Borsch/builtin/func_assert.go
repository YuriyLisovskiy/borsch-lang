package builtin

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

func assert(state common.State, expected common.Value, actual common.Value, errorTemplate string) error {
	args := []common.Value{actual}
	result, err := types.CallByName(state, expected, common.EqualsOp.Name(), &args, nil, true)
	if err != nil {
		return err
	}

	success, err := mustBool(
		result, func(t common.Value) error {
			return errors.New(
				fmt.Sprintf(
					"результат порівняння має бути логічного типу, отримано %s",
					t.GetTypeName(),
				),
			)
		},
	)
	if err != nil {
		return err
	}

	if success {
		return nil
	}

	errMsg := ""
	if errorTemplate != "" {
		errMsg = errorTemplate
	} else {
		errMsg = "не вдалося підтвердити, що %s дорівнює %s"
	}

	expectedStr, err := expected.String(state)
	if err != nil {
		return err
	}

	actualStr, err := actual.String(state)
	if err != nil {
		return err
	}

	return utilities.RuntimeError(fmt.Sprintf(errMsg, expectedStr, actualStr))
}

func buildAssertMessage(state common.State, args *[]common.Value) (string, error) {
	message := ""
	if len(*args) > 2 {
		messageArgs := (*args)[2:]
		sz := len(messageArgs)
		for c := 0; c < sz; c++ {
			argStr, err := messageArgs[c].String(state)
			if err != nil {
				return "", err
			}

			message += argStr
			if c < sz-1 {
				message += " "
			}
		}
	}

	return message, nil
}

func evalAssert(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
	message, err := buildAssertMessage(state, args)
	if err != nil {
		return nil, err
	}

	return types.NewNilInstance(), assert(state, (*args)[0], (*args)[1], message)
}

func makeAssertFunction() *types.FunctionInstance {
	return types.NewFunctionInstance(
		"підтвердити",
		[]types.FunctionParameter{
			{
				Type:       types.Any,
				Name:       "очікуване",
				IsVariadic: false,
				IsNullable: true,
			},
			{
				Type:       types.Any,
				Name:       "фактичне",
				IsVariadic: false,
				IsNullable: true,
			},
			{
				Type:       types.String,
				Name:       "повідомлення_про_помилку",
				IsVariadic: true,
				IsNullable: false,
			},
		},
		evalAssert,
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
