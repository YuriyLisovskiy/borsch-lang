package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
)

func Call(state State, fn Object, args Tuple) (Object, error) {
	if I, ok := fn.(I__call__); ok {
		return I.__call__(state, args)
	}

	return nil, ErrorNewf(TypeError, "об'єкт '%s' не може бути викликаний", fn.Class().Name)
}

// Represent calls __представлення__ on the object or returns a sensible default.
func Represent(self Object) (Object, error) {
	if I, ok := self.(I__represent__); ok {
		return I.__represent__()
	} else if res, ok, err := TypeCall0(self, builtin.RepresentOperatorName); ok {
		return res, err
	}

	return String(fmt.Sprintf("<%s instance at %p>", self.Class().Name, self)), nil
}

// Str calls common.StringOperator on the object and if not found
// calls common.RepresentOperator.
func Str(self Object) (Object, error) {
	if I, ok := self.(I__str__); ok {
		return I.__str__()
	} else if res, ok, err := TypeCall0(self, builtin.StringOperatorName); ok {
		return res, err
	}

	return Represent(self)
}

// StrAsString returns object as a string.
//
// Calls Str then makes sure the output is a string.
func StrAsString(self Object) (string, error) {
	res, err := Str(self)
	if err != nil {
		return "", err
	}

	str, ok := res.(String)
	if !ok {
		return "", ErrorNewf(TypeError, "result of __str__ must be string, not '%s'", res.Class().Name)
	}

	return string(str), nil
}

// Not return the result of not 'a'.
func Not(a Object) (Object, error) {
	b, err := MakeBool(a)
	if err != nil {
		return nil, err
	}

	switch b {
	case False:
		return True, nil
	case True:
		return False, nil
	}

	return nil, ErrorNewf(TypeError, "логічний() не повернув ні 'істина', ні 'хиба'")
}

// MakeBool is called to implement truth value testing and the built-in
// operation логічний(); should return False or True. When this method is
// not defined, __довжина__() is called, if it is defined, and the object
// is considered true if its result is nonzero. If a class defines
// neither __довжина__() nor __логічний__(), all its instances are considered
// true.
func MakeBool(a Object) (Object, error) {
	if _, ok := a.(Bool); ok {
		return a, nil
	}

	if A, ok := a.(I__bool__); ok {
		res, err := A.__bool__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	if B, ok := a.(I__length__); ok {
		res, err := B.__length__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return MakeBool(res)
		}
	}

	return True, nil
}

// MakeGoInt turns 'a' into Go int if possible.
func MakeGoInt(a Object) (int, error) {
	a, err := MakeInt(a)
	if err != nil {
		return 0, err
	}

	A, ok := a.(IGoInt)
	if ok {
		return A.GoInt()
	}

	return 0, ErrorNewf(TypeError, "об'єкт '%v' не може бути інтрпретований як ціле число", a.Class().Name)
}

// Index the Object returning an Int.
//
// Will raise TypeError if Index can't be run on this object
func Index(a Object) (Int, error) {
	if A, ok := a.(I__index__); ok {
		return A.__index__()
	}

	if A, ok, err := TypeCall0(a, builtin.IndexOperatorName); ok {
		if err != nil {
			return 0, err
		}

		if res, ok := A.(Int); ok {
			return res, nil
		}

		return 0, ErrorNewf(
			TypeError,
			"'%s' повернув не ціле число: (тип %s)",
			builtin.IndexOperatorName,
			A.Class().Name,
		)
	}

	return 0, ErrorNewf(
		TypeError,
		"непідтримуваний(і) тип(и) операнда для %s: '%s'",
		builtin.IndexOperatorName,
		a.Class().Name,
	)
}

// IndexInt the Object returning an int.
//
// Will raise TypeError if Index can't be run on this object
//
// or IndexError if the Int won't fit!
func IndexInt(a Object) (int, error) {
	i, err := Index(a)
	if err != nil {
		return 0, err
	}

	intI := int(i)

	// Int might not fit in an int
	if Int(intI) != i {
		return 0, ErrorNewf(IndexError, "cannot fit %d into an index-sized integer", i)
	}

	return intI, nil
}

// IndexIntCheck as IndexInt but if index is -ve addresses it from the end.
//
// If index is out of range throws IndexError
func IndexIntCheck(a Object, max int) (int, error) {
	i, err := IndexInt(a)
	if err != nil {
		return 0, err
	}

	if i < 0 {
		i += max
	}

	if i < 0 || i >= max {
		return 0, ErrorNewf(IndexError, "індекс за межами діапазону")
	}

	return i, nil
}

// RepresentAsString returns object as a string.
//
// Calls Represent then makes sure the output is a string.
func RepresentAsString(self Object) (string, error) {
	res, err := Represent(self)
	if err != nil {
		return "", err
	}

	str, ok := res.(String)
	if !ok {
		return "", ErrorNewf(
			TypeError,
			"результат '%s' має бути рядком, отримано '%s'",
			builtin.RepresentOperatorName,
			res.Class().Name,
		)
	}

	return string(str), nil
}

// GetAttrString - returns the result or an error to be raised if not found.
//
// If not found err will be an AttributeError.
func GetAttrString(self Object, key string) (res Object, err error) {
	// Call __get_attribute__ unconditionally if it exists
	if I, ok := self.(I__get_attribute__); ok {
		return I.__get_attribute__(key)
	} else if res, ok, err = TypeCall1(self, "__get_attribute__", Object(String(key))); ok {
		return res, err
	}

	// Look up any __special__ methods as M__special__ and return a bound method
	if len(key) >= 5 && strings.HasPrefix(key, "__") && strings.HasSuffix(key, "__") {
		objectValue := reflect.ValueOf(self)
		methodValue := objectValue.MethodByName(key)
		if methodValue.IsValid() {
			return newBoundMethod(key, methodValue.Interface())
		}
	}

	// Look in the instance dictionary if it exists
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		res, ok = dict[key]
		if ok {
			return res, err
		}
	}

	// Now look in type's dictionary etc
	t := self.Class()
	res = t.NativeGetAttrOrNil(key)
	if res != nil {
		// Call __get__ which creates bound methods, reads properties etc
		if I, ok := res.(I__get__); ok {
			res, err = I.__get__(self, t)
		}

		return res, err
	}

	// And now only if not found call __getattr__
	if I, ok := self.(I__get_attr__); ok {
		return I.__get_attr__(key)
	} else if res, ok, err = TypeCall1(self, "__get_attr__", Object(String(key))); ok {
		return res, err
	}

	// Not found - return nil
	return nil, ErrorNewf(AttributeError, "'%s' has no attribute '%s'", self.Class().Name, key)
}

// SetAttrString - returns nil or an error to be raised if not found.
//
// If not found err will be an AttributeError.
func SetAttrString(self Object, key string, value Object) error {
	// Call __set_attribute__ unconditionally if it exists
	if I, ok := self.(I__set_attribute__); ok {
		return I.__set_attribute__(key, value)
	} else if _, ok, err := TypeCall2(self, "__set_attribute__", Object(String(key)), value); ok {
		return err
	}

	// Look in the instance dictionary if it exists
	// if I, ok := self.(IGetDict); ok {
	// 	dict := I.GetDict()
	// 	_, ok = dict[key]
	// 	if ok {
	// 		return nil
	// 	}
	// }

	// And now only if not found call __set_attr__
	if I, ok := self.(I__set_attr__); ok {
		return I.__set_attr__(key, value)
	} else if _, ok, err := TypeCall2(self, "__set_attr__", Object(String(key)), value); ok {
		return err
	}

	// Not found - return nil
	return ErrorNewf(AttributeError, "'%s' has no attribute '%s'", self.Class().Name, key)
}

func GetItem(self Object, key Object) (Object, error) {
	if s, ok := self.(I__get_item__); ok {
		return s.__get_item__(key)
	}

	return Nil, ErrorNewf(TypeError, "'%s' object is not subscriptable", self.Class().Name)
}

func SetItem(self Object, key Object, value Object) (Object, error) {
	if s, ok := self.(I__set_item__); ok {
		return s.__set_item__(key, value)
	}

	return Nil, ErrorNewf(TypeError, "'%s' object is not subscriptable", self.Class().Name)
}
