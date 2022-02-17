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
	Dict            StringDict // anything else that we want to stuff in
}

var (
	BaseError = ObjectClass.NewType(
		"БазоваПомилка",
		"Спільний базовий клас для усіх помилок",
		ErrorNew,
		nil,
	)
	ErrorType = BaseError.NewType(
		"Помилка",
		"Спільний базовий клас для усіх помилок, визначених розробником",
		nil,
		nil,
	)
	TypeError  = ErrorType.NewType("ПомилкаТипу", "Невідповідний тип аргументу.", nil, nil)
	ValueError = ErrorType.NewType(
		"ПомилкаЗначення",
		"Невідповідне значення аргументу (з коректним типом).",
		nil,
		nil,
	)
	RuntimeError        = ErrorType.NewType("ПомилкаВиконання", "Невизначена помилка виконання.", nil, nil)
	NotImplementedError = RuntimeError.NewType(
		"ПомилкаВідсутностіРеалізації",
		"Метод або функція не реалізована.",
		nil,
		nil,
	)
	StopIteration   = ErrorType.NewType("ЗупинкаІтерування", "Сигнал зупинки з ітератор.__наступний__().", nil, nil)
	ArithmeticError = ErrorType.NewType("АрифметичнаПомилка", "Базовий клас для аорифметичних помилок.", nil, nil)
	OverflowError   = ArithmeticError.NewType(
		"ПомилкаПереповнення",
		"Результат є занадто великим, щоб його преставити.",
		nil,
		nil,
	)
	ZeroDivisionError = ArithmeticError.NewType(
		"ПомилкаДіленняНаНуль",
		"Другий аргумент при діленні або обчисленні модуля був нулем.",
		nil,
		nil,
	)
	LookupError = ErrorType.NewType("ПомилкаПошуку", "Базовий клас для помилок пошуку.", nil, nil)
	KeyError    = LookupError.NewType("ПомилкаКлюча", "Ключ відображення не знайдено.", nil, nil)
	IndexError  = LookupError.NewType("ПомилкаІндексу", "Індекс послідовності за межами її діапазону.", nil, nil)
	SystemError = ErrorType.NewType(
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
	NotImplemented, err = ErrorNew(NotImplementedError, nil, nil)
	if err != nil {
		log.Fatalf("Failed to make NotImplemented")
	}
}

// Type of this object
func (e *Error) Type() *Class {
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
		Dict: StringDict{},
	}
}

func ErrorNew(cls *Class, args Tuple) (Object, error) {
	return errorNew(cls, args), nil
}

func errorNew(cls *Class, args Tuple) *Error {
	return &Error{
		Base: cls,
		Args: args.Copy(),
		Dict: StringDict{},
	}
}
