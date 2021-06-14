package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/builtin/types"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
	"strings"
)

func Panic(args ...types.ValueType) (types.ValueType, error) {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, arg.String())
	}

	return types.NoneType{}, util.RuntimeError(strings.Join(strArgs, " "))
}

func GetEnv(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 1 {
		return types.StringType{Value: os.Getenv(args[0].String())}, nil
	}

	return types.NoneType{}, util.RuntimeError("функція 'середовище()' приймає лише один аргумент")
}

func Assert(args ...types.ValueType) (types.ValueType, error) {
	if len(args) >= 2 && len(args) <= 3 {
		leftV := args[0]
		rightV := args[1]
		if leftV.TypeHash() != rightV.TypeHash() {
			return nil, util.RuntimeError(
				fmt.Sprintf(
					"неможливо застосувати оператор умови рівності до значень типів '%s' та '%s'",
					leftV.TypeName(), rightV.TypeName(),
				),
			)
		}

		errMsg := "не вдалося підтвердити, що %s дорівнює %s"
		if len(args) == 3 {
			errMsg = args[2].String()
		}

		switch left := leftV.(type) {
		case types.NoneType:
			return nil, nil
		case types.RealType:
			right := rightV.(types.RealType)
			if left.Value != right.Value {
				return nil, util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
			}

			return nil, nil
		case types.IntegerType:
			right := rightV.(types.IntegerType)
			if left.Value != right.Value {
				return nil, util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
			}

			return nil, nil
		case types.StringType:
			right := rightV.(types.StringType)
			if left.Value != right.Value {
				return nil, util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
			}

			return nil, nil
		case types.BoolType:
			right := rightV.(types.BoolType)
			if left.Value != right.Value {
				return nil, util.RuntimeError(fmt.Sprintf(errMsg, left.String(), right.String()))
			}

			return nil, nil
		}

		return nil, util.RuntimeError(fmt.Sprintf(
			"непідтримувані типи операндів для оператора умови рівності: '%s' і '%s'",
			leftV.TypeName(), rightV.TypeName(),
		))
	}

	return nil, util.RuntimeError("функція 'підтвердити()' приймає два, або три аргументи")
}
