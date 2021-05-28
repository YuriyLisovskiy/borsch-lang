package builtin

import (
	"errors"
	"os"
)

func GetEnv(values ...ValueType) (ValueType, error) {
	if len(values) == 1 {
		return StringType{Value: os.Getenv(values[0].String())}, nil
	}

	return NoneType{}, errors.New("функція 'середовище' приймає один аргумент")
}
