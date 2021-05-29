package builtin

import (
	"github.com/YuriyLisovskiy/borsch/src/util"
	"os"
)

func GetEnv(values ...ValueType) (ValueType, error) {
	if len(values) == 1 {
		return StringType{Value: os.Getenv(values[0].String())}, nil
	}

	return NoneType{}, util.RuntimeError("середовище() приймає лише один аргумент")
}
