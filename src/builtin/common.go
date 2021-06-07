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
		strArgs = append(strArgs, arg.Representation())
	}

	return types.NoneType{}, util.RuntimeError(strings.Join(strArgs, " "))
}

func GetEnv(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 1 {
		return types.StringType{Value: os.Getenv(args[0].Representation())}, nil
	}

	return types.NoneType{}, util.RuntimeError("функція 'середовище()' приймає лише один аргумент")
}

func Length(args ...types.ValueType) (types.ValueType, error) {
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case types.SequentialType:
			return types.IntegerType{Value: int64(arg.Length())}, nil
		}

		return types.NoneType{}, util.RuntimeError(fmt.Sprintf(
			"об'єкт типу '%s' не має довжини", args[0].TypeName(),
		))
	}

	return types.NoneType{}, util.RuntimeError("функція 'довжина()' приймає лише один аргумент")
}
