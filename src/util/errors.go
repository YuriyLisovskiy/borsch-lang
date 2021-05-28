package util

import (
	"errors"
	"fmt"
)

func RuntimeError(text string) error {
	return errors.New(fmt.Sprintf("Помилка виконання: %s", text))
}

func SyntaxError(text string) error {
	return errors.New(fmt.Sprintf("Синтаксичка помилка: %s", text))
}
