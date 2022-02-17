// Copyright 2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Real objects.

package types

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var RealType = ObjectType.NewType(
	"дійсний",
	`дійсний(х) -> число з плаваючою комою

Перетворює рядок або число у число з плаваючою комою, якщо це можливо.`,
	RealNew,
	nil,
)

// Bits of precision in a float64
const (
	float64precision   = 53
	float64MaxExponent = 1023
)

type Real float64

// Type of this Real64 object
func (value Real) Type() *Type {
	return RealType
}

func RealNew(cls *Type, args Tuple, kwargs StringDict) (Object, error) {
	var xObj Object = Real(0)
	err := ParseTupleAndKeywords(args, kwargs, "|O", []string{"х"}, &xObj)
	if err != nil {
		return nil, err
	}

	// Special case converting string types
	switch x := xObj.(type) {
	// FIXME Bytearray
	// case Bytes:
	// 	return RealFromString(string(x))
	case String:
		return RealFromString(string(x))
	}

	return MakeReal(xObj)
}

func (value Real) __str__() (Object, error) {
	if i := int64(value); Real(i) == value {
		return String(fmt.Sprintf("%d.0", i)), nil
	}

	return String(fmt.Sprintf("%g", value)), nil
}

func (value Real) __represent__() (Object, error) {
	return value.__str__()
}

// RealFromString turns a string into a Real.
func RealFromString(str string) (Object, error) {
	str = strings.TrimSpace(str)
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			if numErr.Err == strconv.ErrRange {
				if str[0] == '-' {
					return Real(math.Inf(-1)), nil
				} else {
					return Real(math.Inf(1)), nil
				}
			}
		}

		return nil, ErrorNewf(ValueError, "некоректний літерал для дійсного числа: '%s'", str)
	}

	return Real(f), nil
}

var expectingReal = ErrorNewf(TypeError, "необхідне дійсне число")

// RealCheckExact returns the real value of obj if it is exactly a real.
func RealCheckExact(obj Object) (Real, error) {
	f, ok := obj.(Real)
	if !ok {
		return 0, expectingReal
	}

	return f, nil
}

// RealCheck returns the real value of obj if it is a real subclass.
func RealCheck(obj Object) (Real, error) {
	// FIXME: should be checking subclasses
	return RealCheckExact(obj)
}

func RealAsReal64(obj Object) (float64, error) {
	f, err := RealCheck(obj)
	if err == nil {
		return float64(f), nil
	}

	fObj, err := MakeReal(obj)
	if err != nil {
		return 0, err
	}

	f, err = RealCheck(fObj)
	if err == nil {
		return float64(f), nil
	}

	return float64(f), err
}

// Arithmetic

// Errors
var floatDivisionByZero = ErrorNewf(ZeroDivisionError, "ділення дійсного числа на нуль")

// Convert an Object to Real.
//
// Returns ok if the conversion worked or not.
func convertToReal(other Object) (Real, bool) {
	switch b := other.(type) {
	case Real:
		return b, true
	case Int:
		return Real(b), true
	case *BigInt:
		x, err := b.Real()
		return x, err == nil
	case Bool:
		if b {
			return Real(1), true
		} else {
			return Real(0), true
		}
	}

	return 0, false
}

func (value Real) __neg__() (Object, error) {
	return -value, nil
}

func (value Real) __pos__() (Object, error) {
	return value, nil
}

func (value Real) __add__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return value + b, nil
	}

	return NotImplemented, nil
}

func (value Real) __reversed_add__(other Object) (Object, error) {
	return value.__add__(other)
}

func (value Real) __in_place_add__(other Object) (Object, error) {
	return value.__add__(other)
}

func (value Real) __sub__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return value - b, nil
	}

	return NotImplemented, nil
}

func (value Real) __reversed_sub__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return b - value, nil
	}

	return NotImplemented, nil
}

func (value Real) __in_place_sub__(other Object) (Object, error) {
	return value.__sub__(other)
}

func (value Real) __mul__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return value * b, nil
	}

	return NotImplemented, nil
}

func (value Real) __reversed_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value Real) __in_place_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value Real) __div__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		if b == 0 {
			return nil, floatDivisionByZero
		}

		return value / b, nil
	}

	return NotImplemented, nil
}

func (value Real) __reversed_div__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		if value == 0 {
			return nil, floatDivisionByZero
		}

		return b / value, nil
	}

	return NotImplemented, nil
}

func (value Real) __in_place_div__(other Object) (Object, error) {
	return value.__div__(other)
}

// Does DivMod of two floating point numbers
func realDivMod(a, b Real) (Real, Real, error) {
	if b == 0 {
		return 0, 0, floatDivisionByZero
	}

	q := Real(math.Floor(float64(a / b)))
	r := a - q*b
	return q, r, nil
}

func (value Real) __mod__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		_, r, err := realDivMod(value, b)
		return r, err
	}

	return NotImplemented, nil
}

func (value Real) __reversed_mod__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		_, r, err := realDivMod(b, value)
		return r, err
	}

	return NotImplemented, nil
}

func (value Real) __in_place_mod__(other Object) (Object, error) {
	return value.__mod__(other)
}

func (value Real) __pow__(other, modulus Object) (Object, error) {
	if modulus != Nil {
		return NotImplemented, nil
	}

	if b, ok := convertToReal(other); ok {
		return Real(math.Pow(float64(value), float64(b))), nil
	}

	return NotImplemented, nil
}

func (value Real) __reversed_pow__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return Real(math.Pow(float64(b), float64(value))), nil
	}

	return NotImplemented, nil
}

func (value Real) __in_place_pow__(other, modulus Object) (Object, error) {
	return value.__pow__(other, modulus)
}

func (value Real) __bool__() (Object, error) {
	return NewBool(value != 0), nil
}

func (value Real) __int__() (Object, error) {
	if value >= IntMin && value <= IntMax {
		return Int(value), nil
	}

	frac, exp := math.Frexp(float64(value))          // x = frac << exp; 0.5 <= abs(x) < 1
	fracInt := int64(frac * (1 << float64precision)) // x = frac << (exp - float64precision)
	res := big.NewInt(fracInt)
	shift := exp - float64precision
	switch {
	case shift > 0:
		res.Lsh(res, uint(shift))
	case shift < 0:
		res.Rsh(res, uint(-shift))
	}

	return (*BigInt)(res), nil
}

func (value Real) __real__() (Object, error) {
	return value, nil
}

// Comparison

func (value Real) __less_than__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return NewBool(value < b), nil
	}

	return NotImplemented, nil
}

func (value Real) __less_or_equal__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return NewBool(value <= b), nil
	}

	return NotImplemented, nil
}

func (value Real) __equal__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return NewBool(value == b), nil
	}

	return NotImplemented, nil
}

func (value Real) __not_equal__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return NewBool(value != b), nil
	}

	return NotImplemented, nil
}

func (value Real) __greater_than__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return NewBool(value > b), nil
	}

	return NotImplemented, nil
}

func (value Real) __greater_or_equal__(other Object) (Object, error) {
	if b, ok := convertToReal(other); ok {
		return NewBool(value >= b), nil
	}

	return NotImplemented, nil
}

// Check interface is satisfied
var _ realArithmetic = Real(0)
var _ conversionBetweenTypes = Real(0)
var _ I__bool__ = Real(0)
var _ iComparison = Real(0)
