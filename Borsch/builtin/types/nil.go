// Copyright 2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Nil objects

package types

type NilType struct{}

var (
	NilTypeClass = NewClass("нульовий", "")

	Nil = NilType(struct{}{})
)

func (value NilType) Class() *Class {
	return NilTypeClass
}

func (value NilType) __bool__() (Object, error) {
	return False, nil
}

func (value NilType) __str__() (Object, error) {
	return value.__represent__()
}

func (value NilType) __represent__() (Object, error) {
	return String("нуль"), nil
}

// Convert an Object to an NilType.
//
// Returns ok if the conversion worked or not.
func convertToNilType(other Object) (NilType, bool) {
	switch b := other.(type) {
	case NilType:
		return b, true
	}

	return Nil, false
}

func (value NilType) __equal__(other Object) (Object, error) {
	if _, ok := convertToNilType(other); ok {
		return True, nil
	}

	return False, nil
}

func (value NilType) __not_equal__(other Object) (Object, error) {
	if _, ok := convertToNilType(other); ok {
		return False, nil
	}

	return True, nil
}

// Check interface is satisfied
var _ I__bool__ = Nil
var _ I__str__ = Nil
var _ I__represent__ = Nil
var _ I__equal__ = Nil
var _ I__not_equal__ = Nil
