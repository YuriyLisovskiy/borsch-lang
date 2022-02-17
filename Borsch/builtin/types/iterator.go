// Copyright 2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Iterator objects.

package types

// Iterator is a Borsch iterator object.
type Iterator struct {
	Pos  int
	Objs []Object
}

var IteratorType = NewType("ітератор", "тип ітератора")

// Type of this object
func (value *Iterator) Type() *Type {
	return IteratorType
}

// NewIterator defines a new iterator.
func NewIterator(Objs []Object) *Iterator {
	return &Iterator{
		Pos:  0,
		Objs: Objs,
	}
}

func (value *Iterator) __iter__() (Object, error) {
	return value, nil
}

// __next__ returns next one from the iteration.
func (value *Iterator) __next__() (Object, error) {
	if value.Pos >= len(value.Objs) {
		return nil, StopIteration
	}

	r := value.Objs[value.Pos]
	value.Pos++
	return r, nil
}

// Check interface is satisfied
var _ I_iterator = (*Iterator)(nil)
