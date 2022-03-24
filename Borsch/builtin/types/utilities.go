package types

import "fmt"

func Represent(ctx Context, self Object) (Object, error) {
	if v, ok := self.(IRepresent); ok {
		return v.represent(ctx)
	}

	return String(fmt.Sprintf("<об'єкт %s з адресою %p>", self.Class().Name, self)), nil
}

func ToString(ctx Context, self Object) (Object, error) {
	if v, ok := self.(IString); ok {
		return v.string(ctx)
	}

	return Represent(ctx, self)
}

func ToGoString(ctx Context, self Object) (string, error) {
	s, err := ToString(ctx, self)
	if err != nil {
		return "", err
	}

	goString, ok := s.(String)
	if !ok {
		return "", ErrorNewf("результат виклику '__рядок__' має бути типу 'рядок', отримано '%s'", s.Class().Name)
	}

	return string(goString), nil
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
		return A.toInt(ctx)
	}

	// TODO: TypeError
	return nil, ErrorNewf("непідтримуваний тип операнда для 'цілий': '%s'", a.Class().Name)
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

	// TODO: TypeError
	return 0, ErrorNewf("об'єкт '%v' не може бути інтрпретований як ціле число", a.Class().Name)
}

func GetAttribute(ctx Context, self Object, name string) (Object, error) {
	if v, ok := self.(IGetAttribute); ok {
		return v.getAttribute(ctx, name)
	}

	if v, ok := self.(*Class); ok {
		if attr := v.GetAttributeOrNil(name); attr != nil {
			return attr, nil
		}
	}

	return nil, ErrorNewf("'%s' не містить атрибута '%s'", self.Class().Name, name)
}

func SetAttribute(ctx Context, self Object, name string, value Object) error {
	if v, ok := self.(ISetAttribute); ok {
		return v.setAttribute(ctx, name, value)
	}

	if v, ok := self.(*Class); ok {
		if attr := v.GetAttributeOrNil(name); attr != nil {
			if attr.Class() != value.Class() {
				return ErrorNewf(
					"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
					value.Class().Name,
					name,
					attr.Class().Name,
				)
			}
		}

		v.Dict[name] = value
		return nil
	}

	return ErrorNewf("'%s' не містить атрибута '%s'", self.Class().Name, name)
}

func DeleteAttribute(ctx Context, self Object, name string) (Object, error) {
	if v, ok := self.(IDeleteAttribute); ok {
		return v.deleteAttribute(ctx, name)
	}

	if v, ok := self.(*Class); ok {
		if attr := v.DeleteAttributeOrNil(name); attr != nil {
			return attr, nil
		}
	}

	return nil, ErrorNewf("'%s' не містить атрибута '%s'", self.Class().Name, name)
}

func Call(ctx Context, self Object, args Tuple) (Object, error) {
	if v, ok := self.(ICall); ok {
		return v.call(args)
	}

	return nil, ErrorNewf("неможливо застосувати оператор виклику до об'єкта з типом '%s'", self.Class().Name)
}
