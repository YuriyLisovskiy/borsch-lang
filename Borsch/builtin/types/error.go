package types

import "fmt"

type LangException interface {
	Error() string
	Class() *Class
}

var ErrorClass *Class

type Error struct {
	message string
	dict    map[string]Object
}

func (value *Error) Error() string {
	return fmt.Sprintf("%s: %s", ErrorClass.Name, value.message)
}

func (value *Error) Class() *Class {
	return ErrorClass
}

func (value *Error) allocate() {
	if value.Class().Dict != nil {
		value.dict = map[string]Object{}
		for name, attr := range value.Class().Dict {
			if m, ok := attr.(*Method); ok && m.IsMethod() {
				value.dict[name] = &MethodWrapper{
					Method:   m,
					Instance: value,
				}
			}
		}
	}
}

func ErrorNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	message, err := errorMessageFromArgs(ctx, cls, args)
	if err != nil {
		return nil, err
	}

	e := &Error{message: message}
	e.allocate()
	return e, nil
}

func NewError(text string) *Error {
	e := &Error{message: text}
	e.allocate()
	return e
}

func NewErrorf(format string, args ...interface{}) *Error {
	e := &Error{message: fmt.Sprintf(format, args...)}
	e.allocate()
	return e
}

func (value *Error) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *Error) string(ctx Context) (Object, error) {
	return String(value.message), nil
}

func (value *Error) getAttribute(_ Context, name string) (Object, error) {
	if attr, ok := value.dict[name]; ok {
		return attr, nil
	}

	if attr := value.Class().GetAttributeOrNil(name); attr != nil {
		return attr, nil
	}

	return nil, NewErrorf("об'єкт '%s' не містить атрибута '%s'", value.Class().Name, name)
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
		"__рядок__": MethodNew(
			"__рядок__",
			pkg,
			[]MethodParameter{
				{
					Class:      ErrorClass,
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
				return args[0].(IString).string(ctx)
			},
		),
	}
}
