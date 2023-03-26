package types

import "fmt"

var IdentifierErrorClass *Class

type IdentifierError struct {
	message string
}

func (value *IdentifierError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *IdentifierError) Class() *Class {
	return IdentifierErrorClass
}

func IdentifierErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &IdentifierError{message: message}, nil
}

func NewIdentifierError(text string) *IdentifierError {
	return &IdentifierError{message: text}
}

func NewIdentifierErrorf(format string, args ...interface{}) *IdentifierError {
	return &IdentifierError{message: fmt.Sprintf(format, args...)}
}

func (value *IdentifierError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *IdentifierError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
