package types

import "fmt"

var AssertionErrorClass *Class

type AssertionError struct {
	message string
}

func (value *AssertionError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *AssertionError) Class() *Class {
	return AssertionErrorClass
}

func AssertionErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	e := &AssertionError{message: message}
	return e, nil
}

func NewAssertionError(text string) *AssertionError {
	return &AssertionError{message: text}
}

func NewAssertionErrorf(format string, args ...interface{}) *AssertionError {
	return &AssertionError{message: fmt.Sprintf(format, args...)}
}

func (value *AssertionError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *AssertionError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
