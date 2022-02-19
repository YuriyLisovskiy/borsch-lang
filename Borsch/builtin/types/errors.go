package types

import (
	"fmt"
	"log"
)

// Error object.
type Error struct {
	Base            *Class
	Args            Object
	Traceback       Object
	Context         Object
	Cause           Object
	SuppressContext bool
	Dict            Dict // anything else that we want to stuff in
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
	TypeError  = ErrorType.NewClass("ПомилкаТипу", "Невідповідний тип аргументу.", nil, nil)
	ValueError = ErrorType.NewClass(
		"ПомилкаЗначення",
		"Невідповідне значення аргументу (з коректним типом).",
		nil,
		nil,
	)
	RuntimeError        = ErrorType.NewClass("ПомилкаВиконання", "Невизначена помилка виконання.", nil, nil)
	NotImplementedError = RuntimeError.NewClass(
		"ПомилкаВідсутностіРеалізації",
		"Метод або функція не реалізована.",
		nil,
		nil,
	)
	StopIteration   = ErrorType.NewClass("ЗупинкаІтерування", "Сигнал зупинки з ітератор.__наступний__().", nil, nil)
	ArithmeticError = ErrorType.NewClass("АрифметичнаПомилка", "Базовий клас для аорифметичних помилок.", nil, nil)
	OverflowError   = ArithmeticError.NewClass(
		"ПомилкаПереповнення",
		"Результат є занадто великим, щоб його преставити.",
		nil,
		nil,
	)
	ZeroDivisionError = ArithmeticError.NewClass(
		"ПомилкаДіленняНаНуль",
		"Другий аргумент при діленні або обчисленні модуля був нулем.",
		nil,
		nil,
	)
	AttributeError = ErrorType.NewClass("ПомилкаАтрибута", "Атрибут не знайдено.", nil, nil)
	ImportError    = ErrorType.NewClass("ПомилкаІмпорту", "Імпорт не зміг знайти пакет.", nil, nil)
	LookupError    = ErrorType.NewClass("ПомилкаПошуку", "Базовий клас для помилок пошуку.", nil, nil)
	KeyError       = LookupError.NewClass("ПомилкаКлюча", "Ключ відображення не знайдено.", nil, nil)
	IndexError     = LookupError.NewClass("ПомилкаІндексу", "Індекс послідовності за межами її діапазону.", nil, nil)
	SystemError    = ErrorType.NewClass(
		"СистемнаПомилка",
		`Внутрішня помилка в інтерпретаторі Borsch.

Повідомте, будь-ласка, розробника інтерпретатора, включивши в повідомлення
стек помилкок, версію Borsch, а також платформу, де відбулася помилка, та її версію.`,
		nil,
		nil,
	)

	// Singleton exceptions

	NotImplemented Object
)

func init() {
	var err error
	NotImplemented, err = ErrorNew(NotImplementedError, nil)
	if err != nil {
		log.Fatalf("Failed to make NotImplemented")
	}
}

func (e *Error) Class() *Class {
	return e.Base
}

// Go error interface
func (e *Error) Error() string {
	message := e.Base.Name
	if args, ok := e.Args.(Tuple); ok {
		for i, arg := range args {
			if i == 0 {
				message += ": "
			} else {
				message += ", "
			}

			representation, err := RepresentAsString(arg)
			if err == nil {
				message += representation
			} else {
				message += "?"
			}
		}
	}

	if e.Dict["lineno"] != nil {
		message = fmt.Sprintf(
			"\n  Файл \"%v\", рядок %v, позиція %v\n    %s\n\n",
			e.Dict["filename"],
			e.Dict["lineno"],
			e.Dict["offset"],
			e.Dict["line"],
		) + message
	}

	return message
}

func ErrorNewf(cls *Class, format string, a ...interface{}) error {
	message := fmt.Sprintf(format, a...)
	return &Error{
		Base: cls,
		Args: Tuple{String(message)},
		Dict: Dict{},
	}
}

func ErrorNew(cls *Class, args Tuple) (Object, error) {
	return errorNew(cls, args), nil
}

func errorNew(cls *Class, args Tuple) *Error {
	return &Error{
		Base: cls,
		Args: args.Copy(),
		Dict: Dict{},
	}
}
