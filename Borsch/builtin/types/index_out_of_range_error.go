package types

import "fmt"

var IndexOutOfRangeErrorClass *Class

type IndexOutOfRangeError struct {
	message string
}

func (value *IndexOutOfRangeError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *IndexOutOfRangeError) Class() *Class {
	return IndexOutOfRangeErrorClass
}

func IndexOutOfRangeErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &IndexOutOfRangeError{message: message}, nil
}

func NewIndexOutOfRangeError(text string) *IndexOutOfRangeError {
	return &IndexOutOfRangeError{message: text}
}

func NewIndexOutOfRangeErrorf(format string, args ...interface{}) *IndexOutOfRangeError {
	return &IndexOutOfRangeError{message: fmt.Sprintf(format, args...)}
}

func (value *IndexOutOfRangeError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *IndexOutOfRangeError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
