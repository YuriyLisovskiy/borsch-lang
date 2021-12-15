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

func AttributeNotFoundError(objTypeName, attrName string) error {
	return RuntimeError(fmt.Sprintf(
		"об'єкт типу '%s' не містить атрибута з назвою '%s'", objTypeName, attrName,
	))
}

func AttributeIsReadOnlyError(objTypeName, attrName string) error {
	return RuntimeError(fmt.Sprintf(
		"атрибут '%s' об'єкта типу '%s' призначений лише для читання", attrName, objTypeName,
	))
}

func OperatorError(opName, lType, rType string) error {
	return RuntimeError(fmt.Sprintf(
		"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
		opName, lType, rType,
	))
}
