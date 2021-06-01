package builtin

import (
	"fmt"
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
	"strings"
	"unicode/utf8"
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

func Length(args ...ValueType) (ValueType, error) {
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case StringType:
			return IntegerNumberType{
				Value: int64(utf8.RuneCountInString(arg.Representation())),
			}, nil
		}

		return NoneType{}, util.RuntimeError(fmt.Sprintf(
			"об'єкт типу '%s' не має довжини", args[0].TypeName(),
		))
	}

	return NoneType{}, util.RuntimeError("функція 'довжина()' приймає лише один аргумент")
}
