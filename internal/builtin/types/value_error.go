package types

import "fmt"

var ValueErrorClass *Class

type ValueError struct {
	message string
}

func (value *ValueError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *ValueError) Class() *Class {
	return ValueErrorClass
}

func ValueErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &ValueError{message: message}, nil
}

func NewValueError(text string) *ValueError {
	return &ValueError{message: text}
}

func NewValueErrorf(format string, args ...interface{}) *ValueError {
	return &ValueError{message: fmt.Sprintf(format, args...)}
}

func (value *ValueError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *ValueError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
