// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.
//
// Method objects
//
// This is about the type 'builtin_function_or_method', not Python
// methods in user-defined classes.  See class.go for the latter.

package types

import (
	"fmt"
)

// Types for methods

// FunctionNoArgs called with self only.
type FunctionNoArgs func(State, Object) (Object, error)

const (
	// These two constants are not used to indicate the calling convention
	// but the binding when use with methods of classes. These may not be
	// used for functions defined for modules. At most one of these flags
	// may be set for any given method.

	// The method will be passed the type object as the first parameter
	// rather than an instance of the type. This is used to create class
	// methods, similar to what is created when using the classmethod()
	// built-in function.
	METH_CLASS = 0x0010

	// The method will be passed NULL as the first parameter rather than
	// an instance of the type. This is used to create static methods,
	// similar to what is created when using the staticmethod() built-in
	// function.
	METH_STATIC = 0x0020

	// One other constant controls whether a method is loaded in
	// place of another definition with the same method name.

	// The method will be loaded in place of existing definitions. Without
	// METH_COEXIST, the default is to skip repeated definitions. Since
	// slot wrappers are loaded before the method table, the existence of
	// a sq_contains slot, for example, would generate a wrapped method
	// named __contains__() and preclude the loading of a corresponding
	// Function with the same name. With the flag defined, the
	// Function will be loaded in place of the wrapper object and will
	// co-exist with the slot. This is helpful because calls to
	// PyCFunctions are optimized more than wrapper object calls.
	METH_COEXIST = 0x0040
)

// A python Method object
// TODO: add method args and expected return type!
type Method struct {
	// Name of this function
	Name string
	// Doc string
	Doc string
	// Flags - see METH_* flags
	Flags int
	// Go function implementation
	method interface{}
	// Parent module of this method
	Package *Package
}

// Internal method types implemented within eval.go
type InternalMethod int

const (
	InternalMethodNone InternalMethod = iota
	InternalMethodGlobals
	InternalMethodLocals
	InternalMethodImport
	InternalMethodEval
	InternalMethodExec
)

var MethodType = NewClass("метод", "об_єкт методу")

func (value *Method) Class() *Class {
	return MethodType
}

// Define a new method
func NewMethod(name string, method interface{}, flags int, doc string) (*Method, error) {
	// have to write out the function arguments - can't use the
	// type aliases as they are different types :-(
	switch method.(type) {
	case func(state State, self Object, args Tuple) (Object, error):
	case FunctionNoArgs:
	case func(State, Object, Object) (Object, error):
	case InternalMethod:
	default:
		return nil, ErrorNewf(SystemError, "Unknown function type for NewMethod %q, %T", name, method)
	}

	return &Method{
		Name:   name,
		Doc:    doc,
		Flags:  flags,
		method: method,
	}, nil
}

// MustNewMethod As NewMethod but panics on error.
func MustNewMethod(name string, method interface{}, flags int, doc string) *Method {
	m, err := NewMethod(name, method, flags, doc)
	if err != nil {
		panic(err)
	}

	return m
}

// Internal returns the InternalMethod type of this method.
func (value *Method) Internal() InternalMethod {
	if internalMethod, ok := value.method.(InternalMethod); ok {
		return internalMethod
	}

	return InternalMethodNone
}

// Call the method with the given arguments
func (value *Method) Call(state State, self Object, args Tuple) (Object, error) {
	// TODO: check method args

	state = state.WithContext(value.Package.Context)
	switch f := value.method.(type) {
	case func(state State, self Object, args Tuple) (Object, error):
		return f(state, self, args)
	case FunctionNoArgs:
		if len(args) != 0 {
			return nil, ErrorNewf(TypeError, "%s() не приймає аргументів (отримано %d)", value.Name, len(args))
		}

		return f(state, self)
	case func(State, Object, Object) (Object, error):
		if len(args) != 1 {
			return nil, ErrorNewf(TypeError, "%s() приймає точно 1 аргумент (отримано %d)", value.Name, len(args))
		}

		return f(state, self, args[0])
	}

	// TODO: check method return type if it matches with actual result

	panic(fmt.Sprintf("Unknown method type: %T", value.method))
}

// Return a new Method with the bound method passed in, or an error
//
// This needs to convert the methods into internally callable python
// methods
func newBoundMethod(name string, fn interface{}) (Object, error) {
	m := &Method{
		Name: name,
	}
	switch f := fn.(type) {
	case func(args Tuple) (Object, error):
		m.method = func(_ State, _ Object, args Tuple) (Object, error) {
			return f(args)
		}
	// __str__() (Object, error)
	case func() (Object, error):
		m.method = func(_ State, _ Object) (Object, error) {
			return f()
		}
	// __add__(other Object) (Object, error)
	case func(Object) (Object, error):
		m.method = func(_ State, _ Object, other Object) (Object, error) {
			return f(other)
		}
	// __get_attr__(name string) (Object, error)
	case func(string) (Object, error):
		m.method = func(_ State, _ Object, stringObject Object) (Object, error) {
			name, err := StrAsString(stringObject)
			if err != nil {
				return nil, err
			}

			return f(name)
		}
	// __get__(instance, owner Object) (Object, error)
	case func(Object, Object) (Object, error):
		m.method = func(_ State, _ Object, args Tuple) (Object, error) {
			var a, b Object
			err := UnpackTuple(args, name, 2, 2, &a, &b)
			if err != nil {
				return nil, err
			}

			return f(a, b)
		}
	// __new__(cls, args, kwargs Object) (Object, error)
	case func(Object, Object, Object) (Object, error):
		m.method = func(_ State, _ Object, args Tuple) (Object, error) {
			var a, b, c Object
			err := UnpackTuple(args, name, 3, 3, &a, &b, &c)
			if err != nil {
				return nil, err
			}

			return f(a, b, c)
		}
	default:
		return nil, fmt.Errorf("unknown bound method type for %q: %T", name, fn)
	}

	return m, nil
}

// Call a method
func (value *Method) __call__(state State, args Tuple) (Object, error) {
	self := Object(value.Package)
	return value.Call(state, self, args)
}

// Read a method from a class which makes a bound method
func (value *Method) __get__(instance, owner Object) (Object, error) {
	panic("unreachable")
	// TODO:
	// if instance != Nil {
	// 	return NewBoundMethod(instance, m), nil
	// }
	//
	// return m, nil
}

// FIXME this should be the default?
func (value *Method) __equal__(other Object) (Object, error) {
	if otherMethod, ok := other.(*Method); ok && value == otherMethod {
		return True, nil
	}

	return False, nil
}

// FIXME this should be the default?
func (value *Method) __not_equal__(other Object) (Object, error) {
	if otherMethod, ok := other.(*Method); ok && value == otherMethod {
		return False, nil
	}

	return True, nil
}

// Make sure it satisfies the interface
var _ Object = (*Method)(nil)
var _ I__call__ = (*Method)(nil)

// var _ I__get__ = (*Method)(nil)
var _ I__equal__ = (*Method)(nil)
var _ I__not_equal__ = (*Method)(nil)
