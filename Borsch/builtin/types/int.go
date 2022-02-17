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
	"math"
	"math/big"
	"strconv"
	"strings"
)

var IntType = ObjectType.NewType(
	"цілий", `
цілий(x=0) -> ціле число
цілий(x, база=10) -> ціле число

Перетворює число або рядок у ціле число, або повертає 0 якщо не було
передано жодних аргументів. Якщо x є числом, повертає x.__цілий__().
У дійсних чисел обрізає дробову частину.

Якщо x не є числом або якщо базу задано, тоді x екземпляром рядка,
який представляє ціле число у заданій базі. Літералу може
передувати '+' чи '-' та може бути огорнутий пропусками.
За замовчуванням база дорівнює 10. Коректними є бази 0 та 2-36.
База 0 означає, що потрібно інтерпретувати базу з рядка як ціле число.
>>> цілий('0b100', база=0)
4`, IntNew, nil,
)

type Int int64

const (
	// Maximum possible Int
	IntMax = math.MaxInt64
	// Minimum possible Int
	IntMin = math.MinInt64
	// The largest number such that sqrtIntMax**2 < IntMax
	sqrtIntMax = 3037000499
	// Go integer limits
	GoUintMax = ^uint(0)
	GoUintMin = 0
	GoIntMax  = int(GoUintMax >> 1)
	GoIntMin  = -GoIntMax - 1
)

// Type of this Int object
func (value Int) Type() *Type {
	return IntType
}

func IntNew(cls *Type, args Tuple, kwargs StringDict) (Object, error) {
	var xObj Object = Int(0)
	var baseObj Object
	base := 0
	err := ParseTupleAndKeywords(args, kwargs, "|OO:int", []string{"х", "база"}, &xObj, &baseObj)
	if err != nil {
		return nil, err
	}

	if baseObj != nil {
		base, err = MakeGoInt(baseObj)
		if err != nil {
			return nil, err
		}

		if base != 0 && (base < 2 || base > 36) {
			return nil, ErrorNewf(ValueError, "база для 'цілий()' має бути >= 2 та <= 36")
		}
	}

	// Special case converting string types
	switch x := xObj.(type) {
	// FIXME Bytearray
	// case Bytes:
	// 	return IntFromString(string(x), base)
	case String:
		return IntFromString(string(x), base)
	}

	if baseObj != nil {
		return nil, ErrorNewf(TypeError, "'цілий()' не може перетворити об'єкт не рядкового типу з явною базою")
	}

	return MakeInt(xObj)
}

// IntFromString Create an Int (or BigInt) from the string passed in.
func IntFromString(str string, base int) (Object, error) {
	var x *big.Int
	var ok bool
	s := str
	negative := false
	convertBase := base

	// Get rid of padding
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		goto error
	}

	// Get rid of sign
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			negative = true
		}
		s = s[1:]
		if len(s) == 0 {
			goto error
		}
	}

	if len(s) > 1 && s[0] == '0' {
		switch s[1] {
		case 'x', 'X':
			convertBase = 16
		case 'o', 'O':
			convertBase = 8
		case 'b', 'B':
			convertBase = 2
		default:
			goto nosigil
		}

		if base != 0 && base != convertBase {
			// int("0xFF", 10)
			// int("0b", 16)
			convertBase = base
			goto nosigil
		}

		s = s[2:]
		if len(s) == 0 {
			goto error
		}
	nosigil:
	}
	if convertBase == 0 {
		convertBase = 10
	}

	// Detect leading zeros which Borsch doesn't allow using base 0
	if base == 0 {
		if len(s) > 1 && s[0] == '0' && (s[1] >= '0' && s[1] <= '9') {
			goto error
		}
	}

	// Use int64 conversion for short strings since 12**36 < IntMax
	// and 10**18 < IntMax
	if len(s) <= 12 || (convertBase <= 10 && len(s) <= 18) {
		i, err := strconv.ParseInt(s, convertBase, 64)
		if err != nil {
			goto error
		}

		if negative {
			i = -i
		}

		return Int(i), nil
	}

	// The base argument must be 0 or a value from 2 through
	// 36. If the base is 0, the string prefix determines the
	// actual conversion base. A prefix of "0x" or "0X" selects
	// base 16; the "0" prefix selects base 8, and a “0b” or “0B”
	// prefix selects base 2. Otherwise, the selected base is 10.
	x, ok = new(big.Int).SetString(s, convertBase)
	if !ok {
		goto error
	}

	if negative {
		x.Neg(x)
	}

	return (*BigInt)(x).MaybeInt(), nil

error:
	return nil, ErrorNewf(ValueError, "некоректний літерал для 'цілий()' з базою %d: '%s'", convertBase, str)
}

// GoInt truncates to Go int.
//
// If it is outside the range of Go int it will return an error.
func (value Int) GoInt() (int, error) {
	r := int(value)
	if Int(r) != value {
		return 0, overflowErrorGo
	}

	return r, nil
}

// GoInt64 truncates to Go int64.
//
// If it is outside the range of Go int64 it will return an error
func (value Int) GoInt64() (int64, error) {
	return int64(value), nil
}

func (value Int) __str__() (Object, error) {
	return String(fmt.Sprintf("%d", value)), nil
}

func (value Int) __represent__() (Object, error) {
	return value.__str__()
}

// Arithmetic

// Errors
var (
	divisionByZero     = ErrorNewf(ZeroDivisionError, "ділення на нуль")
	negativeShiftCount = ErrorNewf(ValueError, "від'ємна кількість біт для побітового зсуву")
)

// Constructs a TypeError
func cantConvert(a Object, to string) (Object, error) {
	return nil, ErrorNewf(TypeError, "неможливо перетворити %s у %s", a.Type().Name, to)
}

// Converts an Object to an Int.
//
// Returns ok if the conversion worked or not.
func convertToInt(other Object) (Int, bool) {
	switch b := other.(type) {
	case Int:
		return b, true
	case Bool:
		if b {
			return Int(1), true
		} else {
			return Int(0), true
		}
		// case Real:
		// 	ib := Int(b)
		// 	if Real(ib) == b {
		// 		return ib, true
		// 	}
	}

	return 0, false
}

func (value Int) __neg__() (Object, error) {
	if value == IntMin {
		// Up-convert overflowing case
		r := big.NewInt(IntMin)
		r.Neg(r)
		return (*BigInt)(r), nil
	}

	return -value, nil
}

func (value Int) __pos__() (Object, error) {
	return value, nil
}

func (value Int) __invert__() (Object, error) {
	return ^value, nil
}

// Integer add with overflow detection.
func intAdd(a, b Int) Object {
	if a >= 0 {
		// Overflow when a + b > IntMax
		// b > IntMax - a
		// IntMax - a can't overflow since
		// IntMax = 7FFF, a = 0..7FFF
		if b > IntMax-a {
			goto overflow
		}
	} else {
		// Underflow when a + b < IntMin
		// => b < IntMin-a
		// IntMin-a can't overflow since
		// IntMin=-8000, a = -8000..-1
		if b < IntMin-a {
			goto overflow
		}
	}

	return a + b

overflow:
	aBig := big.NewInt(int64(a))
	bBig := big.NewInt(int64(b))
	aBig.Add(aBig, bBig)
	return (*BigInt)(aBig).MaybeInt()
}

// Integer subtract with overflow detection.
func intSub(a, b Int) Object {
	if b >= 0 {
		// Underflow when a - b < IntMin
		// a < IntMin + b
		// IntMin + b can't overflow since
		// IntMin = -8000, b 0..7FFF
		if a < IntMin+b {
			goto overflow
		}
	} else {
		// Overflow when a - b > IntMax
		// a < IntMax + b
		// IntMax + b can't overflow since
		// IntMax=7FFF, b = -8000..-1, IntMax + b = -1..0x7FFE
		if a < IntMax+b {
			goto overflow
		}
	}

	return a - b

overflow:
	aBig := big.NewInt(int64(a))
	bBig := big.NewInt(int64(b))
	aBig.Sub(aBig, bBig)
	return (*BigInt)(aBig).MaybeInt()
}

// Integer multiplication with overflow detection.
func intMul(a, b Int) Object {
	absA := a
	if a < 0 {
		absA = -a
	}
	absB := b
	if b < 0 {
		absB = -b
	}

	// A crude but effective test!
	if absA <= sqrtIntMax && absB <= sqrtIntMax {
		return a * b
	}

	aBig := big.NewInt(int64(a))
	bBig := big.NewInt(int64(b))
	aBig.Mul(aBig, bBig)
	return (*BigInt)(aBig).MaybeInt()
}

// Left shift a << b.
func intLeftShift(a, b Int) (Object, error) {
	if b < 0 {
		return nil, negativeShiftCount
	}

	shift := uint(b)
	r := a << shift
	if r>>shift != a {
		aBig := big.NewInt(int64(a))
		aBig.Lsh(aBig, shift)
		return (*BigInt)(aBig), nil
	}

	return r, nil
}

func (value Int) __add__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intAdd(value, b), nil
	}

	return NotImplemented, nil
}

func (value Int) __reversed_add__(other Object) (Object, error) {
	return value.__add__(other)
}

func (value Int) __in_place_add__(other Object) (Object, error) {
	return value.__add__(other)
}

func (value Int) __sub__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intSub(value, b), nil
	}

	return NotImplemented, nil
}

func (value Int) __reversed_sub__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intSub(b, value), nil
	}

	return NotImplemented, nil
}

func (value Int) __in_place_sub__(other Object) (Object, error) {
	return value.__sub__(other)
}

func (value Int) __mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intMul(value, b), nil
	}

	return NotImplemented, nil
}

func (value Int) __reversed_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value Int) __in_place_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value Int) __div__(other Object) (Object, error) {
	b, err := MakeReal(other)
	if err != nil {
		return nil, err
	}
	fa := Real(value)
	if err != nil {
		return nil, err
	}

	fb := b.(Real)
	if fb == 0 {
		return nil, divisionByZero
	}

	return Real(fa / fb), nil
}

func (value Int) __reversed_div__(other Object) (Object, error) {
	b, err := MakeReal(other)
	if err != nil {
		return nil, err
	}

	fa := Real(value)
	if err != nil {
		return nil, err
	}

	fb := b.(Real)
	if fa == 0 {
		return nil, divisionByZero
	}

	return Real(fb / fa), nil
}

func (value Int) __in_place_div__(other Object) (Object, error) {
	return value.__div__(other)
}

func (value Int) __mod__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		_, result, err := value.divMod(b)
		return result, err
	}

	return NotImplemented, nil
}

func (value Int) __reversed_mod__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		_, result, err := b.divMod(value)
		return result, err
	}

	return NotImplemented, nil
}

func (value Int) __in_place_mod__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		_, result, err := value.divMod(b)
		return result, err
	}

	return NotImplemented, nil
}

func (value Int) divMod(b Int) (Object, Object, error) {
	if b == 0 {
		return nil, nil, divisionByZero
	}

	// Can't overflow
	result, remainder := value/b, value%b

	// Implement floor division
	negativeResult := value < 0
	if b < 0 {
		negativeResult = !negativeResult
	}

	if negativeResult && remainder != 0 {
		result -= 1
		remainder += b
	}

	return result, remainder, nil
}

func (value Int) __pow__(other, modulus Object) (Object, error) {
	return (*BigInt)(big.NewInt(int64(value))).__pow__(other, modulus)
}

func (value Int) __reversed_pow__(other Object) (Object, error) {
	return (*BigInt)(big.NewInt(int64(value))).__reversed_pow__(other)
}

func (value Int) __in_place_pow__(other, modulus Object) (Object, error) {
	return value.__pow__(other, modulus)
}

func (value Int) __left_shift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intLeftShift(value, b)
	}

	return NotImplemented, nil
}

func (value Int) __reversed_left_shift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return intLeftShift(b, value)
	}

	return NotImplemented, nil
}

func (value Int) __in_place_left_shift__(other Object) (Object, error) {
	return value.__left_shift__(other)
}

func (value Int) __right_shift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			return nil, negativeShiftCount
		}

		// Can't overflow
		return value >> uint64(b), nil
	}

	return NotImplemented, nil
}

func (value Int) __reversed_right_shift__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			return nil, negativeShiftCount
		}

		// Can't overflow
		return b >> uint64(value), nil
	}

	return NotImplemented, nil
}

func (value Int) __in_place_right_shift__(other Object) (Object, error) {
	return value.__right_shift__(other)
}

func (value Int) __and__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return value & b, nil
	}

	return NotImplemented, nil
}

func (value Int) __reversed_and__(other Object) (Object, error) {
	return value.__and__(other)
}

func (value Int) __in_place_and__(other Object) (Object, error) {
	return value.__and__(other)
}

func (value Int) __xor__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return value ^ b, nil
	}

	return NotImplemented, nil
}

func (value Int) __reversed_xor__(other Object) (Object, error) {
	return value.__xor__(other)
}

func (value Int) __in_place_xor__(other Object) (Object, error) {
	return value.__xor__(other)
}

func (value Int) __or__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return value | b, nil
	}

	return NotImplemented, nil
}

func (value Int) __reversed_or__(other Object) (Object, error) {
	return value.__or__(other)
}

func (value Int) __in_place_or__(other Object) (Object, error) {
	return value.__or__(other)
}

func (value Int) __bool__() (Object, error) {
	return NewBool(value != 0), nil
}

func (value Int) __index__() (Int, error) {
	return value, nil
}

func (value Int) __int__() (Object, error) {
	return value, nil
}

func (value Int) __real__() (Object, error) {
	if r, ok := convertToReal(value); ok {
		return r, nil
	}

	return cantConvert(value, "дійсний")
}

// Rich comparison

func (value Int) __less_than__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(value < b), nil
	}

	return NotImplemented, nil
}

func (value Int) __less_or_equal__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(value <= b), nil
	}

	return NotImplemented, nil
}

func (value Int) __equal__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(value == b), nil
	}

	return NotImplemented, nil
}

func (value Int) __not_equal__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(value != b), nil
	}

	return NotImplemented, nil
}

func (value Int) __greater_than__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(value > b), nil
	}

	return NotImplemented, nil
}

func (value Int) __greater_or_equal__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		return NewBool(value >= b), nil
	}

	return NotImplemented, nil
}

// Check interface is satisfied
var _ realArithmetic = Int(0)
var _ booleanArithmetic = Int(0)
var _ conversionBetweenTypes = Int(0)
var _ I__bool__ = Int(0)
var _ I__index__ = Int(0)
var _ iComparison = Int(0)
var _ IGoInt = Int(0)
var _ IGoInt64 = Int(0)
