package builtin

import (
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
