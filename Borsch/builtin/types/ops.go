package types

import (
	"fmt"
	"math"
	"strings"
)

func mod(l, r Real) Real {
	a := float64(l)
	b := float64(r)
	return Real(math.Mod(b+math.Mod(a, b), b))
}

func Represent(ctx Context, self Object) (Object, error) {
	if v, ok := self.(IRepresent); ok {
		res, err := v.represent(ctx)
		if err != nil {
			return nil, err
		}

		if _, ok := res.(String); !ok {
			return nil, NewTypeErrorf(
				"результат виклику '__представлення__' має бути типу 'рядок', отримано '%s'",
				res.Class().Name,
			)
		}

		return res, nil
	}

	return String(fmt.Sprintf("<об'єкт '%s' з адресою %p>", self.Class().Name, self)), nil
}

func ToString(ctx Context, self Object) (Object, error) {
	if v, ok := self.(IString); ok {
		res, err := v.string(ctx)
		if err != nil {
			return nil, err
		}

		if _, ok := res.(String); !ok {
			return nil, NewTypeErrorf(
				"результат виклику '__рядок__' має бути типу 'рядок', отримано '%s'",
				res.Class().Name,
			)
		}

		return res, nil
	}

	return Represent(ctx, self)
}

func ToGoString(ctx Context, self Object) (string, error) {
	s, err := ToString(ctx, self)
	if err != nil {
		return "", err
	}

	return string(s.(String)), nil
}

func ToBool(ctx Context, self Object) (Object, error) {
	if _, ok := self.(Bool); ok {
		return self, nil
	}

	if b, ok := self.(IBool); ok {
		result, err := b.toBool(ctx)
		if err != nil {
			return nil, err
		}

		return ToBool(ctx, result)
	}

	return True, nil
}

// ToInt the Object returning an Object.
//
// Will raise TypeError if MakeInt can't be run on this object.
func ToInt(ctx Context, a Object) (Object, error) {
	if _, ok := a.(Int); ok {
		return a, nil
	}

	if A, ok := a.(IInt); ok {
		result, err := A.toInt(ctx)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewTypeErrorf("непідтримуваний тип операнда для 'ціле': '%s'", a.Class().Name)
}

func ToReal(ctx Context, a Object) (Object, error) {
	if _, ok := a.(Real); ok {
		return a, nil
	}

	if A, ok := a.(IReal); ok {
		return A.toReal(ctx)
	}

	return nil, NewTypeErrorf("непідтримуваний тип операнда для 'дійсне': '%s'", a.Class().Name)
}

// ToGoInt turns 'a' into Go int if possible.
func ToGoInt(ctx Context, a Object) (int, error) {
	a, err := ToInt(ctx, a)
	if err != nil {
		return 0, err
	}

	if v, ok := a.(IGoInt); ok {
		return v.toGoInt(ctx)
	}

	return 0, NewTypeErrorf("об'єкт '%v' не може бути інтерпретований як ціле число", a.Class().Name)
}

func GetAttribute(ctx Context, self Object, name string) (Object, error) {
	if v, ok := self.(IGetAttribute); ok {
		return v.getAttribute(ctx, name)
	}

	return nil, NewErrorf("'%s' не містить атрибута '%s'", self.Class().Name, name)
}

func SetAttribute(ctx Context, self Object, name string, value Object) error {
	if v, ok := self.(ISetAttribute); ok {
		return v.setAttribute(ctx, name, value)
	}

	return NewAttributeErrorf("'%s' не містить атрибута '%s'", self.Class().Name, name)
}

func DeleteAttribute(ctx Context, self Object, name string) (Object, error) {
	if v, ok := self.(IDeleteAttribute); ok {
		return v.deleteAttribute(ctx, name)
	}

	return nil, NewErrorf("'%s' не містить атрибута '%s'", self.Class().Name, name)
}

func Call(ctx Context, self Object, args Tuple) (Object, error) {
	if v, ok := self.(ICall); ok {
		return v.call(ctx, args)
	}

	return nil, NewErrorf("неможливо застосувати оператор виклику до об'єкта з типом '%s'", self.Class().Name)
}

func parseArgs(name, format string, args Tuple, argsMin, argsMax int, results ...*Object) error {
	typesFormat, nullablesFormat, err := parseFormat(format)
	if len(typesFormat) != len(results) {
		return NewErrorf("Internal Error: supply the same number of results and types in format")
	}

	if err = checkNumberOfArgs(name, len(args), len(results), argsMin, argsMax); err != nil {
		return err
	}

	for i, arg := range args {
		result := results[i]
		isNullable := nullablesFormat[i] == '?'
		if arg.Class() == NilClass {
			if !isNullable {
				return NewErrorf( /*TypeError,*/ "%s() аргумент виклику %d не може бути нульовим", name, i+1)
			}
		} else {
			extra := ""
			if isNullable {
				extra = "або 'нульове'"
			}

			t := typesFormat[i]
			switch t {
			case 'b':
				if _, ok := arg.(Bool); !ok {
					return NewErrorf(
						"%s() аргумент %d має бути типу 'логічне'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
					)
				}
			case 'i':
				if _, ok := arg.(Int); !ok {
					return NewErrorf(
						"%s() аргумент %d має бути типу 'ціле'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
					)
				}
			// case 'l':
			// 	if _, ok := arg.(List); !ok {
			// 		return NewErrorf(
			// 			"%s() аргумент %d має бути типу 'список'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
			// 		)
			// 	}
			// case 'm':
			// 	if _, ok := arg.(Method); !ok {
			// 		return NewErrorf(
			// 			"%s() аргумент %d має бути типу 'метод'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
			// 		)
			// 	}
			case 'n':
				if arg != Nil {
					return NewErrorf(
						"%s() аргумент %d має бути типу 'нульове', а не '%s'", name, i+1, arg.Class().Name,
					)
				}
			// case 'p':
			// 	if _, ok := arg.(Package); !ok {
			// 		return NewErrorf(
			// 			"%s() аргумент %d має бути типу 'пакет'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
			// 		)
			// 	}
			case 'r':
				if _, ok := arg.(Real); !ok {
					return NewErrorf(
						"%s() аргумент %d має бути типу 'дійсне'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
					)
				}
			case 's':
				if _, ok := arg.(String); !ok {
					return NewErrorf(
						"%s() аргумент %d має бути типу 'рядок'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
					)
				}
			case 't':
				if _, ok := arg.(*Tuple); !ok {
					return NewErrorf(
						"%s() аргумент %d має бути типу 'кортеж'%s, а не '%s'", name, i+1, extra, arg.Class().Name,
					)
				}
			case 'o':
			default:
				return NewErrorf("Internal Error: unknown type to parse from format")
			}
		}

		*result = arg
	}

	return nil
}

func checkNumberOfArgs(name string, argsN, resultsN, argsMin, argsMax int) error {
	if argsMin == argsMax {
		if argsN != argsMax {
			return NewErrorf( /*TypeError, */ "%s() takes exactly %d arguments (%d given)", name, argsMax, argsN)
		}
	} else {
		if argsN > argsMax {
			return NewErrorf( /*TypeError, */ "%s() takes at most %d arguments (%d given)", name, argsMax, argsN)
		}
		if argsN < argsMin {
			return NewErrorf( /*TypeError, */ "%s() takes at least %d arguments (%d given)", name, argsMin, argsN)
		}
	}

	if argsN > resultsN {
		return NewErrorf( /*TypeError, */ "Internal error: not enough arguments supplied to Unpack*/Parse*")
	}

	return nil
}

// Format has three parts: types|nullables.
//
//  Example of format: we need to parse int, real and string
//  in this order and real arg is nullable:
//      irs|.?.
func parseFormat(format string) (string, string, error) {
	parts := strings.Split(format, "|")
	if len(parts) != 2 {
		return "", "", NewErrorf("Internal Error: provide nullables in format")
	}

	typesFormat := parts[0]
	nullablesFormat := parts[1]
	if len(typesFormat) != len(nullablesFormat) {
		return "", "", NewErrorf("Internal Error: supply the same number of nullables and types in format")
	}

	return typesFormat, nullablesFormat, nil
}

func getAttributeFrom(dict *StringDict, name string, cls *Class) (Object, error) {
	if attr, ok := (*dict)[name]; ok {
		return attr, nil
	}

	if attr := cls.GetAttributeOrNil(name); attr != nil {
		return attr, nil
	}

	return nil, NewErrorf("об'єкт '%s' не містить атрибута '%s'", cls.Name, name)
}

func setAttributeTo(instance Object, dict *StringDict, attr Object, name string, value Object) error {
	if attr != nil && !accepts(attr.Class(), value.Class()) {
		if attr.Class() == MethodWrapperClass {
			switch value.Class() {
			case MethodClass, FunctionClass, LambdaClass:
				(*dict)[name] = wrap(instance, value.(*Method))
				return nil
			}
		}

		return NewTypeErrorf(
			"неможливо записати значення типу '%s' в атрибут '%s' з типом '%s'",
			value.Class().Name,
			name,
			attr.Class().Name,
		)
	}

	(*dict)[name] = value
	return nil
}

func wrapMethod(instance, obj Object) (Object, bool) {
	if method, ok := obj.(*Method); ok && method.IsMethod() {
		return &MethodWrapper{
			Method:   method,
			Instance: instance,
		}, true
	}

	return obj, false
}

func wrap(instance Object, method *Method) Object {
	return &MethodWrapper{
		Method:   method,
		Instance: instance,
	}
}

// initInstance prepares object's attributes using the base class
// and performs any other common operations for object with
// dict.
func initInstance(instance Object, dict *StringDict, cls *Class) {
	if cls.Dict != nil {
		*dict = StringDict{}
		for name, attr := range cls.Dict {
			if m, ok := attr.(*Method); ok && m.IsMethod() {
				(*dict)[name] = &MethodWrapper{
					Method:   m,
					Instance: instance,
				}
			}
		}
	}
}

func accepts(a, b *Class) bool {
	return a == b || a.IsBaseOf(b)
}
