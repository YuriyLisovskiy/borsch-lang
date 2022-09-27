package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
)

type LangException interface {
	Error() string
	Class() *Class
}

var ErrorClass *Class

type Error struct {
	message string
	dict    StringDict
}

func (value *Error) Error() string {
	return fmt.Sprintf("%s: %s", ErrorClass.Name, value.message)
}

func (value *Error) Class() *Class {
	return ErrorClass
}

func (value *Error) init() {
	initInstance(value, &value.dict, value.Class())
}

func ErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	e := &Error{message: message}
	e.init()
	return e, nil
}

func NewError(text string) *Error {
	e := &Error{message: text}
	e.init()
	return e
}

func NewErrorf(format string, args ...interface{}) *Error {
	e := &Error{message: fmt.Sprintf(format, args...)}
	e.init()
	return e
}

func (value *Error) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *Error) string(ctx Context) (Object, error) {
	return String(value.message), nil
}

func (value *Error) getAttribute(_ Context, name string) (Object, error) {
	return getAttributeFrom(&value.dict, name, value.Class())
}

func (value *Error) setAttribute(_ Context, name string, newValue Object) error {
	attr, ok := value.dict[name]
	if !ok {
		attr = value.Class().GetAttributeOrNil(name)
	}

	return setAttributeTo(value, &value.dict, attr, name, newValue)
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

func MakeErrorClassAttributes(pkg *Package) map[string]Object {
	return map[string]Object{
		builtin.StringOperatorName:         makeStringMethod(pkg, ErrorClass),
		builtin.RepresentationOperatorName: makeRepresentationMethod(pkg, ErrorClass),
	}
}

func makeStringMethod(pkg *Package, cls *Class) *Method {
	return MethodNew(
		builtin.StringOperatorName,
		pkg,
		[]MethodParameter{
			{
				Class:      cls,
				Classes:    nil,
				Name:       "я",
				IsNullable: false,
				IsVariadic: false,
			},
		},
		[]MethodReturnType{
			{
				Class:      StringClass,
				IsNullable: false,
			},
		},
		func(ctx Context, args Tuple, kwargs StringDict) (Object, error) {
			if v, ok := args[0].(IString); ok {
				return v.string(ctx)
			}

			panic("unreachable")
		},
	)
}

func makeRepresentationMethod(pkg *Package, cls *Class) *Method {
	return MethodNew(
		builtin.RepresentationOperatorName,
		pkg,
		[]MethodParameter{
			{
				Class:      cls,
				Classes:    nil,
				Name:       "я",
				IsNullable: false,
				IsVariadic: false,
			},
		},
		[]MethodReturnType{
			{
				Class:      StringClass,
				IsNullable: false,
			},
		},
		func(ctx Context, args Tuple, kwargs StringDict) (Object, error) {
			if v, ok := args[0].(IRepresent); ok {
				return v.represent(ctx)
			}

			panic("unreachable")
		},
	)
}
