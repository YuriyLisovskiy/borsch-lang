// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package types

type Bool bool

var (
	BoolType = NewClass(
		"логічний",
		`логічний(x) -> логічний

Повертає 'істина', якщо аргумент x є істиною, інакше - 'хиба'.
The builtins True and False are the only two instances of the class bool.
The class bool is a subclass of the class int, and cannot be subclassed.`,
	)

	False = Bool(false)
	True  = Bool(true)
)

func (b Bool) Class() *Class {
	return BoolType
}

// NewBool returns the canonical True and False values.
func NewBool(t bool) Bool {
	if t {
		return True
	}

	return False
}

func (b Bool) __bool__() (Object, error) {
	return b, nil
}

func (b Bool) __index__() (Int, error) {
	if b {
		return Int(1), nil
	}

	return Int(0), nil
}

func (b Bool) __str__() (Object, error) {
	return b.__represent__()
}

func (b Bool) __represent__() (Object, error) {
	if b {
		return String("істина"), nil
	}

	return String("хиба"), nil
}

// Convert an Object to Bool.
//
// Returns ok if the conversion worked or not.
func convertToBool(other Object) (Bool, bool) {
	switch b := other.(type) {
	case Bool:
		return b, true
	case Int:
		switch b {
		case 0:
			return False, true
		case 1:
			return True, true
		default:
			return False, false
		}
	case Real:
		switch b {
		case 0:
			return False, true
		case 1:
			return True, true
		default:
			return False, false
		}
	}

	return False, false
}

func (b Bool) __equal__(other Object) (Object, error) {
	if o, ok := convertToBool(other); ok {
		return NewBool(b == o), nil
	}

	return False, nil
}

func (b Bool) __not_equal__(other Object) (Object, error) {
	if o, ok := convertToBool(other); ok {
		return NewBool(b != o), nil
	}

	return True, nil
}

func notEq(eq Object, err error) (Object, error) {
	if err != nil {
		return nil, err
	}

	if eq == NotImplemented {
		return eq, nil
	}

	return Not(eq)
}

// Check interface is satisfied
var _ I__bool__ = Bool(false)
var _ I__index__ = Bool(false)
var _ I__str__ = Bool(false)
var _ I__represent__ = Bool(false)
var _ I__equal__ = Bool(false)
var _ I__not_equal__ = Bool(false)
