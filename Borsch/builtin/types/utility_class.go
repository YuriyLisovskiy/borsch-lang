package types

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

// GetAttribute - returns the result or an error to be raised if not found.
//
// If not found err will be an AttributeError.
func GetAttribute(ctx common.Context, self common.Value, key string) (res common.Value, err error) {
	if I, ok := self.(I__get_attribute__); ok {
		return I.__get_attribute__(ctx, key)
	}
	// else if res, ok, err = TypeCall1(self, "__отримати_атрибут__", common.Value(String(key))); ok {
	// 	return res, err
	// }

	// Look up any __special__ methods as M__special__ and return a bound method
	// if len(key) >= 5 && strings.HasPrefix(key, "__") && strings.HasSuffix(key, "__") {
	// 	objectValue := reflect.ValueOf(self)
	// 	methodValue := objectValue.MethodByName(key)
	// 	if methodValue.IsValid() {
	// 		return newBoundMethod(key, methodValue.Interface())
	// 	}
	// }

	// Look in the instance dictionary if it exists
	if I, ok := self.(IGetDict); ok {
		dict := I.GetDict()
		res, ok = dict[key]
		if ok {
			return res, err
		}
	}

	// Now look in type's dictionary etc
	t := self.(ObjectInstance).GetClass()
	res = t.NativeGetAttrOrNil(key)
	if res != nil {
		// Call __get__ which creates bound methods, reads properties etc
		if I, ok := res.(I__get__); ok {
			res, err = I.__get__(self, t)
		}

		return res, err
	}

	// Not found - return nil
	return nil, ErrorNewf(AttributeError, "'%s' has no attribute '%s'", self.(ObjectInstance).GetClass().Name, key)
}

// SetAttribute - returns nil or an error to be raised if not found.
//
// If not found err will be an AttributeError.
func SetAttribute(ctx common.Context, self common.Value, key string, value common.Value) error {
	// Call __set_attribute__ unconditionally if it exists
	if I, ok := self.(I__set_attribute__); ok {
		return I.__set_attribute__(ctx, key, value)
	}
	// else if _, ok, err := TypeCall2(self, "__встановити_attribute__", Object(String(key)), value); ok {
	// 	return err
	// }

	// Look in the instance dictionary if it exists
	// if I, ok := self.(IGetDict); ok {
	// 	dict := I.GetDict()
	// 	_, ok = dict[key]
	// 	if ok {
	// 		return nil
	// 	}
	// }

	// Not found - return nil
	return ErrorNewf(AttributeError, "'%s' has no attribute '%s'", self.(ObjectInstance).GetClass().Name, key)
}
