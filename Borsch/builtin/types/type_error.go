package types

import "fmt"

var TypeErrorClass = ErrorClass.ClassNew("ПомилкаТипу", map[string]Object{}, false, TypeErrorNew, nil)

type TypeError struct {
	message string
}

func (value *TypeError) Error() string {
	return fmt.Sprintf("%s: %s", value.Class().Name, value.message)
}

func (value *TypeError) Class() *Class {
	return TypeErrorClass
}

func TypeErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	return &TypeError{message: message}, nil
}

func NewTypeError(text string) *TypeError {
	return &TypeError{message: text}
}

func NewTypeErrorf(format string, args ...interface{}) *TypeError {
	return &TypeError{message: fmt.Sprintf(format, args...)}
}

func (value *TypeError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *TypeError) string(_ Context) (Object, error) {
	return String(value.message), nil
}
