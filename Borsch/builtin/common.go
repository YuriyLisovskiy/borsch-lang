package builtin

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func Assert(ctx common.Context, expected common.Type, actual common.Type, errorTemplate string) error {
	leftV := expected
	rightV := actual
	if leftV.GetTypeHash() != rightV.GetTypeHash() {
		return util.RuntimeError(
			fmt.Sprintf(
				"неможливо застосувати оператор умови рівності до значень типів '%s' та '%s'",
				leftV.GetTypeName(), rightV.GetTypeName(),
			),
		)
	}

	errMsg := "не вдалося підтвердити, що %s дорівнює %s"
	if errorTemplate != "" {
		errMsg = errorTemplate
	}

	switch left := leftV.(type) {
	case types.NilInstance:
		return nil
	case types.RealInstance:
		right := rightV.(types.RealInstance)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(ctx), right.String(ctx)))
		}

		return nil
	case types.IntegerInstance:
		right := rightV.(types.IntegerInstance)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(ctx), right.String(ctx)))
		}

		return nil
	case types.StringInstance:
		right := rightV.(types.StringInstance)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(ctx), right.String(ctx)))
		}

		return nil
	case types.BoolInstance:
		right := rightV.(types.BoolInstance)
		if left.Value != right.Value {
			return util.RuntimeError(fmt.Sprintf(errMsg, left.String(ctx), right.String(ctx)))
		}

		return nil
	}

	return util.RuntimeError(
		fmt.Sprintf(
			"непідтримувані типи операндів для оператора умови рівності: '%s' і '%s'",
			leftV.GetTypeName(), rightV.GetTypeName(),
		),
	)
}

func Help(word string) error {
	fmt.Println("Поки що не паше =(")
	return nil
}
