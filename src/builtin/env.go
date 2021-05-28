package builtin

import (
	"errors"
	"os"
)

func GetEnv(values ...string) (string, error) {
	if len(values) == 1 {
		return os.Getenv(values[0]), nil
	}

	return "", errors.New("функція 'середовище' приймає один аргумент")
}
