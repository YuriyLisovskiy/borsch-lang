package util

import (
	"errors"
	"fmt"
)

func RuntimeError(text string) error {
	return errors.New(fmt.Sprintf("Помилка виконання: %s", text))
}

func InternalError(text string) error {
	return errors.New(fmt.Sprintf("Внутрішня помилка: %s", text))
}

func AttributeError(objTypeName, attrName string) error {
	return RuntimeError(fmt.Sprintf(
		"об'єкт типу '%s' не містить атрибута з назвою '%s'", objTypeName, attrName,
	))
}
