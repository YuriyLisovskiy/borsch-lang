package types

import (
	"fmt"
	"log"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type Error struct {
	Base    *Class
	Dict    Dict
	Message string
}

func (e *Error) String(common.State) (string, error) {
	return e.Error(), nil
}

func (e *Error) Representation(state common.State) (string, error) {
	return e.String(state)
}

func (e *Error) AsBool(state common.State) (bool, error) {
	return true, nil
}

func (e *Error) GetTypeName() string {
	return e.GetClass().Name
}

func (e *Error) GetOperator(s string) (common.Value, error) {
	// TODO implement me
	panic("implement me")
}

func (e *Error) GetAttribute(s string) (common.Value, error) {
	// TODO implement me
	panic("implement me")
}

func (e *Error) SetAttribute(s string, value common.Value) error {
	// TODO implement me
	panic("implement me")
}

func (e *Error) HasAttribute(s string) bool {
	// TODO implement me
	panic("implement me")
}

func (e *Error) GetClass() *Class {
	return e.Base
}

func (e *Error) Error() string {
	message := e.Base.Name + ": " + e.Message

	// if e.Dict["lineno"] != nil {
	// 	message = fmt.Sprintf(
	// 		"\n  Файл \"%v\", рядок %v, позиція %v\n    %s\n\n",
	// 		e.Dict["filename"],
	// 		e.Dict["lineno"],
	// 		e.Dict["offset"],
	// 		e.Dict["line"],
	// 	) + message
	// }

	return message
}

func ErrorNewf(cls *Class, format string, a ...interface{}) error {
	return &Error{
		Base:    cls,
		Message: fmt.Sprintf(format, a...),
		Dict:    Dict{},
	}
}

var (
	BaseError = ObjectClass.NewClass(
		"БазоваПомилка",
		"Спільний базовий клас для усіх помилок",
		ErrorNew,
		nil,
	)
	ErrorType = BaseError.NewClass(
		"Помилка",
		"Спільний базовий клас для усіх помилок, визначених розробником",
		nil,
		nil,
	)
	AttributeError      = ErrorType.NewClass("ПомилкаАтрибута", "Атрибут не знайдено.", nil, nil)
	TypeError           = ErrorType.NewClass("ПомилкаТипу", "Невідповідний тип аргументу.", nil, nil)
	RuntimeError        = ErrorType.NewClass("ПомилкаВиконання", "Невизначена помилка виконання.", nil, nil)
	NotImplementedError = RuntimeError.NewClass(
		"ПомилкаВідсутностіРеалізації",
		"Метод або функція не реалізована.",
		nil,
		nil,
	)

	// Singleton exceptions

	NotImplemented common.Value
)

func init() {
	var err error
	NotImplemented = &Error{
		Base:    NotImplementedError,
		Dict:    Dict{},
		Message: "",
	}
	if err != nil {
		log.Fatalf("Failed to make NotImplemented")
	}
}
