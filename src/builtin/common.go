package builtin

import (
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
	"strings"
)

func Panic(args ...ValueType) (ValueType, error) {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, arg.Representation())
	}

	return NoneType{}, util.RuntimeError(strings.Join(strArgs, " "))
}

func GetEnv(args ...ValueType) (ValueType, error) {
	if len(args) == 1 {
		return StringType{Value: os.Getenv(args[0].Representation())}, nil
	}

	return NoneType{}, util.RuntimeError("функція 'середовище()' приймає лише один аргумент")
}
