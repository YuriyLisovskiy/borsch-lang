package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
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
	value.dict["повідомлення"] = String(value.message)
}

func ErrorConstruct(ctx Context, self Object, args Tuple) error {
	message, err := errorMessageFromArgs(ctx, nil, args)
	if err != nil {
		return err
	}

	return SetAttribute(ctx, self, "повідомлення", String(message))
}

func NewError(text string) *Error {
	e := &Error{message: text}
	initInstance(e, &e.dict, e.Class())
	e.init()
	return e
}

func NewErrorf(format string, args ...interface{}) *Error {
	e := &Error{message: fmt.Sprintf(format, args...)}
	initInstance(e, &e.dict, e.Class())
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

func MakeErrorClassOperators(pkg *Package) map[common.OperatorHash]*Method {
	return map[common.OperatorHash]*Method{
		common.StringOp:         makeStringMethod(pkg, ErrorClass),
		common.RepresentationOp: makeRepresentationMethod(pkg, ErrorClass),
	}
}

func MakeErrorClassMethods(pkg *Package) StringDict {
	instanceMethod := MethodNew(
		InitializeMethodName, pkg, []MethodParameter{
			{
				Class:      ErrorClass,
				Classes:    nil,
				Name:       "я",
				IsNullable: false,
				IsVariadic: false,
			},
			{
				Class:      ObjectClass,
				Classes:    nil,
				Name:       "повідомлення",
				IsNullable: false,
				IsVariadic: true,
			},
		},
		[]MethodReturnType{
			{
				Class:      NilClass,
				IsNullable: true,
			},
		},
		func(ctx Context, args Tuple, kwargs StringDict) (Object, error) {
			return Nil, ErrorConstruct(ctx, args[0], args[1:])
		},
	)

	return StringDict{
		instanceMethod.Name: instanceMethod,
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
			if instance, ok := args[0].(*Class); ok && instance.IsInstance() {
				message, err := GetAttribute(ctx, instance, "повідомлення")
				if err != nil {
					return nil, err
				}

				if v, ok := message.(IString); ok {
					str, err := v.string(ctx)
					if err != nil {
						return nil, err
					}

					return String(fmt.Sprintf("%s: %s", instance.Class().Name, str)), nil
				}
			}

			panic("unreachable")
		},
	)
}

// TODO: fix an infinite recursion
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
