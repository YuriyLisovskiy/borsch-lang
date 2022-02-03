package builtin

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func Assert(state common.State, expected common.Value, actual common.Value, errorTemplate string) error {
	args := []common.Value{actual}
	result, err := types.CallByName(state, expected, common.EqualsOp.Name(), &args, nil, true)
	if err != nil {
		return err
	}

	success, err := MustBool(
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

	return util.RuntimeError(fmt.Sprintf(errMsg, expectedStr, actualStr))
}

func Help(word string) error {
	fmt.Println("Поки що не паше =(")
	return nil
}
