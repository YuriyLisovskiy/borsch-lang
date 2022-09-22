package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
)

var ZeroDivisionErrorClass *Class

type ZeroDivisionError struct {
	message string
	dict    StringDict
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

	zeroDivErr := &ZeroDivisionError{message: message}
	zeroDivErr.allocate()
	return zeroDivErr, nil
}

func NewZeroDivisionError(text string) *ZeroDivisionError {
	err := &ZeroDivisionError{message: text}
	err.allocate()
	return err
}

func NewZeroDivisionErrorf(format string, args ...interface{}) *ZeroDivisionError {
	err := &ZeroDivisionError{message: fmt.Sprintf(format, args...)}
	err.allocate()
	return err
}

func (value *ZeroDivisionError) allocate() {
	allocate(value, &value.dict, value.Class())
}

func (value *ZeroDivisionError) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *ZeroDivisionError) string(_ Context) (Object, error) {
	return String(value.message), nil
}

func (value *ZeroDivisionError) getAttribute(_ Context, name string) (Object, error) {
	return getAttributeFrom(&value.dict, name, value.Class())
}

func (value *ZeroDivisionError) setAttribute(_ Context, name string, newValue Object) error {
	attr, ok := value.dict[name]
	if !ok {
		attr = value.Class().GetAttributeOrNil(name)
	}

	return setAttributeTo(value, &value.dict, attr, name, newValue)
}

func MakeZeroDivisionErrorClassAttributes(pkg *Package) map[string]Object {
	return map[string]Object{
		builtin.StringOperatorName:         makeStringMethod(pkg, ZeroDivisionErrorClass),
		builtin.RepresentationOperatorName: makeRepresentationMethod(pkg, ZeroDivisionErrorClass),
	}
}
