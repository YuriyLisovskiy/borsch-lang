// Copyright 2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package types

import (
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

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

	return 0, ErrorNewf(TypeError, "об'єкт '%v' не може бути інтрпретований як ціле число", a.Type().Name)
}

// Index the Object returning an Int.
//
// Will raise TypeError if Index can't be run on this object
func Index(a Object) (Int, error) {
	if A, ok := a.(I__index__); ok {
		return A.__index__()
	}

	if A, ok, err := TypeCall0(a, common.IndexOperator); ok {
		if err != nil {
			return 0, err
		}

		if res, ok := A.(Int); ok {
			return res, nil
		}

		return 0, ErrorNewf(TypeError, "'%s' повернув не ціле число: (тип %s)", common.IndexOperator, A.Type().Name)
	}

	return 0, ErrorNewf(
		TypeError,
		"непідтримуваний(і) тип(и) операнда для %s: '%s'",
		common.IndexOperator,
		a.Type().Name,
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

// Call calls function fnObj with args and kwargs.
//
// kwargs should be nil if not required.
//
// fnObj must be a callable type such as *Method or *Function.
//
// The result is returned.
func Call(fn Object, args Tuple, kwargs StringDict) (Object, error) {
	if I, ok := fn.(I__call__); ok {
		return I.__call__(args, kwargs)
	}

	return nil, ErrorNewf(TypeError, "об'єкт '%s' не може бути викликаний", fn.Type().Name)
}

// Represent calls __представлення__ on the object or returns a sensible default.
func Represent(self Object) (Object, error) {
	if I, ok := self.(I__represent__); ok {
		return I.__represent__()
	} else if res, ok, err := TypeCall0(self, common.RepresentOperator); ok {
		return res, err
	}

	return String(fmt.Sprintf("<%s instance at %p>", self.Type().Name, self)), nil
}

// Str calls common.StringOperator on the object and if not found
// calls common.RepresentOperator.
func Str(self Object) (Object, error) {
	if I, ok := self.(I__str__); ok {
		return I.__str__()
	} else if res, ok, err := TypeCall0(self, common.StringOperator); ok {
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
		return "", ErrorNewf(
			TypeError,
			"результат '%s' має бути рядком, отримано '%s'",
			common.StringOperator,
			res.Type().Name,
		)
	}

	return string(str), nil
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
			common.RepresentOperator,
			res.Type().Name,
		)
	}

	return string(str), nil
}
