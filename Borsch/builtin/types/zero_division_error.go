package types

import "fmt"

var ZeroDivisionErrorClass = ErrorClass.ClassNew("ПомилкаДіленняНаНуль", map[string]Object{}, false, ZeroDivisionErrorNew, nil)

type ZeroDivisionError struct {
	message string
}

func (value *ZeroDivisionError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *ZeroDivisionError) Class() *Class {
	return ZeroDivisionErrorClass
}

func ZeroDivisionErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &ZeroDivisionError{message: message}, nil
}

func NewZeroDivisionError(text string) *ZeroDivisionError {
	return &ZeroDivisionError{message: text}
}

func NewZeroDivisionErrorf(format string, args ...interface{}) *ZeroDivisionError {
	return &ZeroDivisionError{message: fmt.Sprintf(format, args...)}
}

func (value *ZeroDivisionError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *ZeroDivisionError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
