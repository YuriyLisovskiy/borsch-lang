package util

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/alecthomas/participle/v2/lexer"
)

func RuntimeError(text string) error {
	return errors.New(fmt.Sprintf("Помилка виконання: %s", text))
}

func InternalError(text string) error {
	return errors.New(fmt.Sprintf("Внутрішня помилка: %s", text))
}

func ParseError(pos lexer.Position, unexpected string, err string) error {
	return errors.New(
		fmt.Sprintf(
			"  Файл \"%s\", рядок %d, позиція %d,\n    %s\n    %s\nСинтаксична помилка: %s",
			pos.Filename,
			pos.Line,
			pos.Column,
			unexpected,
			strings.Repeat(" ", utf8.RuneCountInString(unexpected))+"^",
			err,
		),
	)
}

func AttributeNotFoundError(objTypeName, attrName string) error {
	return RuntimeError(fmt.Sprintf("об'єкт типу '%s' не містить атрибута з назвою '%s'", objTypeName, attrName))
}

func OperatorNotFoundError(objTypeName, opName string) error {
	return RuntimeError(
		fmt.Sprintf(
			"об'єкт типу '%s' не містить оператора з назвою '%s'", objTypeName, opName,
		),
	)
}

func CantSetAttributeOfBuiltinTypeError(objTypeName string) error {
	return RuntimeError(
		fmt.Sprintf(
			"неможливо встановлювати атрибути для вбудованого типу '%s'", objTypeName,
		),
	)
}

func AttributeIsReadOnlyError(objTypeName, attrName string) error {
	return RuntimeError(
		fmt.Sprintf(
			"атрибут '%s' об'єкта типу '%s' призначений лише для читання", attrName, objTypeName,
		),
	)
}

func OperatorNotSupportedError(operator common.Operator, leftType, rightType string) error {
	return RuntimeError(
		fmt.Sprintf(
			"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
			operator.Sign(), leftType, rightType,
		),
	)
}

func OperandsNotSupportedError(operator common.Operator, leftType, rightType string) error {
	return RuntimeError(
		fmt.Sprintf(
			"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
			operator.Sign(), leftType, rightType,
		),
	)
}

func ObjectIsNotCallable(objectName, typeName string) error {
	if objectName != "" {
		objectName = fmt.Sprintf(" '%s'", objectName)
	}

	return RuntimeError(
		fmt.Sprintf(
			"неможливо застосувати оператор виклику до об'єкта%s з типом '%s'", objectName, typeName,
		),
	)
}

type InterpreterError struct {
	message string
}

func (e InterpreterError) Error() string {
	return fmt.Sprintf("InterpreterError: %s", e.message)
}

func IncorrectUseOfFunctionError(functionName string) error {
	return InterpreterError{message: fmt.Sprintf("incorrect use of '%s' func", functionName)}
}
