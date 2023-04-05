package types

import "fmt"

var AttributeErrorClass *Class

type AttributeError struct {
	message string
}

func (value *AttributeError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *AttributeError) Class() *Class {
	return AttributeErrorClass
}

func AttributeErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &AttributeError{message: message}, nil
}

func NewAttributeError(text string) *AttributeError {
	return &AttributeError{message: text}
}

func NewAttributeErrorf(format string, args ...interface{}) *AttributeError {
	return &AttributeError{message: fmt.Sprintf(format, args...)}
}

func (value *AttributeError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *AttributeError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
