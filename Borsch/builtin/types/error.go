package types

import "fmt"

type LangException interface {
	Error() string
	Class() *Class
}

var ErrorClass = ObjectClass.ClassNew("Помилка", map[string]Object{}, false, ErrorNew, nil)

type Error struct {
	message string
}

func (value *Error) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *Error) Class() *Class {
	return ErrorClass
}

func ErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &Error{message: message}, nil
}

func NewError(text string) *Error {
	return &Error{message: text}
}

func NewErrorf(format string, args ...interface{}) *Error {
	return &Error{message: fmt.Sprintf(format, args...)}
}

func (value *Error) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *Error) string(ctx Context) (Object, error) {
	return String(value.message), nil
}

func errorMessageFromArgs(ctx Context, cls *Class, args Tuple) (string, error) {
	message := ""
	for _, arg := range args {
		sArg, err := ToGoString(ctx, arg)
		if err != nil {
			return "", err
		}

		message += sArg
	}

	return message, nil
}

func OperatorNotSupportedErrorNew(operator, lType, rType string) error {
	return NewTypeErrorf(
		"екземпляри типів '%s' і '%s' не підтримують оператор '%s'",
		lType,
		rType,
		operator,
	)
}
