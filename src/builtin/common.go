package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
)

func Assert(expected, actual types.ValueType, errorTemplate string) error {
	leftV := expected
	rightV := actual
	if leftV.TypeHash() != rightV.TypeHash() {
		return util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор умови рівності до значень типів '%s' та '%s'",
				leftV.TypeName(), rightV.TypeName(),
			),
		)
	}

	errMsg := "не вдалося підтвердити, що %s дорівнює %s"
	if errorTemplate != "" {
		errMsg = errorTemplate
	}

	switch left := leftV.(type) {
	case types.NilType:
		return nil
	case types.RealType:
		right := rightV.(types.RealType)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
		}

		return nil
	case types.IntegerType:
		right := rightV.(types.IntegerType)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
		}

		return nil
	case types.StringType:
		right := rightV.(types.StringType)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
		}

		return nil
	case types.BoolType:
		right := rightV.(types.BoolType)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
		}

		return nil
	}

	return util.RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора умови рівності: '%s' і '%s'",
		leftV.TypeName(), rightV.TypeName(),
	))
}

func Help(word string) error {
	fmt.Println("Поки що не паше =(")
	return nil
}
