package types

import "fmt"

var RuntimeErrorClass *Class

type RuntimeError struct {
	message string
}

func (value *RuntimeError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *RuntimeError) Class() *Class {
	return RuntimeErrorClass
}

func RuntimeErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &RuntimeError{message: message}, nil
}

func NewRuntimeError(text string) *RuntimeError {
	return &RuntimeError{message: text}
}

func NewRuntimeErrorf(format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{message: fmt.Sprintf(format, args...)}
}

func (value *RuntimeError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *RuntimeError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
