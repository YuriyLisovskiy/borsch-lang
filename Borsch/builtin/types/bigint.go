// Copyright 2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// BigInt objects.

package types

import (
	"fmt"
	"math"
	"math/big"
)

type BigInt big.Int

var BigIntClass = NewClass("великий_цілий", "Зберігає великі цілі числа")

func (value *BigInt) Class() *Class {
	return BigIntClass
}

func (value *BigInt) __str__() (Object, error) {
	return String(fmt.Sprintf("%d", (*big.Int)(value))), nil
}

func (value *BigInt) __represent__() (Object, error) {
	return value.__str__()
}

// Some common BigInt-s
var (
	bigInt0   = (*BigInt)(big.NewInt(0))
	bigInt1   = (*BigInt)(big.NewInt(1))
	bigInt10  = (*BigInt)(big.NewInt(10))
	bigIntMin = (*BigInt)(big.NewInt(IntMin))
	bigIntMax = (*BigInt)(big.NewInt(IntMax))
)

// Errors
var (
	overflowError      = ErrorNewf(OverflowError, "ціле число занадто велике, щоб перетворити його в int64")
	overflowErrorGo    = ErrorNewf(OverflowError, "ціле число занадто велике, щоб перетворити його в Go int")
	overflowErrorFloat = ErrorNewf(OverflowError, "ціле число занадто велике, щоб перетворити його в дійсне")
	expectingBigInt    = ErrorNewf(TypeError, "необхідне вилике ціле число")
)

// BigIntCheckExact checks that obj is exactly a BigInt and returns an error if not.
func BigIntCheckExact(obj Object) (*BigInt, error) {
	bigInt, ok := obj.(*BigInt)
	if !ok {
		return nil, expectingBigInt
	}
	return bigInt, nil
}

// BigIntCheck checks that obj is exactly a bigInd and returns an error if not.
func BigIntCheck(obj Object) (*BigInt, error) {
	// FIXME should be checking subclasses
	return BigIntCheckExact(obj)
}

// Arithmetic

// ConvertToBigInt converts an Object to an BigInt.
//
// Returns ok if the conversion worked or not.
func ConvertToBigInt(other Object) (*BigInt, bool) {
	switch b := other.(type) {
	case Int:
		return (*BigInt)(big.NewInt(int64(b))), true
	case *BigInt:
		return b, true
	case Bool:
		if b {
			return bigInt1, true
		} else {
			return bigInt0, true
		}
	}

	return nil, false
}

// Int truncates to Int.
//
// If it is outside the range of an Int it will return an error.
func (value *BigInt) Int() (Int, error) {
	if (*big.Int)(value).Cmp((*big.Int)(bigIntMax)) <= 0 && (*big.Int)(value).Cmp((*big.Int)(bigIntMin)) >= 0 {
		return Int((*big.Int)(value).Int64()), nil
	}

	return 0, overflowError
}

// MaybeInt truncates to Int if it can, otherwise returns the original BigInt.
func (value *BigInt) MaybeInt() Object {
	i, err := value.Int()
	if err != nil {
		return value
	}

	return i
}

// GoInt truncates to Go int.
//
// If it is outside the range of Go int it will return an error.
func (value *BigInt) GoInt() (int, error) {
	z, err := value.Int()
	if err != nil {
		return 0, err
	}

	return z.GoInt()
}

// GoInt64 truncates to Go int64.
//
// If it is outside the range of Go int64 it will return an error.
func (value *BigInt) GoInt64() (int64, error) {
	z, err := value.Int()
	if err != nil {
		return 0, err
	}

	return int64(z), nil
}

// Frexp produces frac and exp such that a ~= frac × 2**exp.
func (value *BigInt) Frexp() (frac float64, exp int) {
	aBig := (*big.Int)(value)
	bits := aBig.BitLen()
	exp = bits - 63
	t := new(big.Int).Set(aBig)
	switch {
	case exp > 0:
		t.Rsh(t, uint(exp))
	case exp < 0:
		t.Lsh(t, uint(-exp))
	}

	// t should now have 63 bits of the integer in and will fit in
	// an int64
	return float64(t.Int64()), exp
}

// Real truncates to Real.
//
// If it is outside the range of Real it will return an error.
func (value *BigInt) Real() (Real, error) {
	frac, exp := value.Frexp()
	// FIXME: this is a bit approximate but errs on the low side so
	// we won't ever produce +Inf-s
	if exp > float64MaxExponent-63 {
		return 0, overflowErrorFloat
	}

	return Real(math.Ldexp(frac, exp)), nil
}

func (value *BigInt) __neg__() (Object, error) {
	return (*BigInt)(new(big.Int).Neg((*big.Int)(value))), nil
}

func (value *BigInt) __pos__() (Object, error) {
	return value, nil
}

func (value *BigInt) __invert__() (Object, error) {
	return (*BigInt)(new(big.Int).Not((*big.Int)(value))), nil
}

func (value *BigInt) __add__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Add((*big.Int)(value), (*big.Int)(b))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_add__(other Object) (Object, error) {
	return value.__add__(other)
}

func (value *BigInt) __in_place_add__(other Object) (Object, error) {
	return value.__add__(other)
}

func (value *BigInt) __sub__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Sub((*big.Int)(value), (*big.Int)(b))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_sub__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Sub((*big.Int)(b), (*big.Int)(value))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __in_place_sub__(other Object) (Object, error) {
	return value.__sub__(other)
}

func (value *BigInt) __mul__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Mul((*big.Int)(value), (*big.Int)(b))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value *BigInt) __in_place_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value *BigInt) __div__(other Object) (Object, error) {
	b, err := MakeReal(other)
	if err != nil {
		return nil, err
	}

	fa, err := value.Real()
	if err != nil {
		return nil, err
	}

	fb := b.(Real)
	if fb == 0 {
		return nil, divisionByZero
	}

	return fa / fb, nil
}

func (value *BigInt) __reversed_div__(other Object) (Object, error) {
	b, err := MakeReal(other)
	if err != nil {
		return nil, err
	}

	fa, err := value.Real()
	if err != nil {
		return nil, err
	}

	fb := b.(Real)
	if fa == 0 {
		return nil, divisionByZero
	}

	return fb / fa, nil
}

func (value *BigInt) __in_place_div__(other Object) (Object, error) {
	return value.__div__(other)
}

func (value *BigInt) __mod__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		_, result, err := value.divMod(b)
		return result, err
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_mod__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		_, result, err := b.divMod(value)
		return result, err
	}

	return NotImplemented, nil
}

func (value *BigInt) __in_place_mod__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		_, result, err := value.divMod(b)
		return result, err
	}

	return NotImplemented, nil
}

func (value *BigInt) divMod(b *BigInt) (Object, Object, error) {
	if (*big.Int)(b).Sign() == 0 {
		return nil, nil, divisionByZero
	}

	r := new(big.Int)
	q := new(big.Int)
	q.QuoRem((*big.Int)(value), (*big.Int)(b), r)

	// Implement floor division
	negativeResult := (*big.Int)(value).Sign() < 0
	if (*big.Int)(b).Sign() < 0 {
		negativeResult = !negativeResult
	}

	if negativeResult && r.Sign() != 0 {
		q.Sub(q, (*big.Int)(bigInt1))
		r.Add(r, (*big.Int)(b))
	}

	return (*BigInt)(q).MaybeInt(), (*BigInt)(r).MaybeInt(), nil
}

// Raise to the power a**b or if m != nil, a**b mod m
func (value *BigInt) pow(b, m *BigInt) (Object, error) {
	// -ve power => make real
	if (*big.Int)(b).Sign() < 0 {
		if m != nil {
			return nil, ErrorNewf(
				TypeError,
				"Другий аргумент оператора '**' не може бути негативний, якщо визначено третій аргумент",
			)
		}

		fa, err := value.Real()
		if err != nil {
			return nil, err
		}

		fb, err := b.Real()
		if err != nil {
			return nil, err
		}

		return fa.__pow__(fb, Nil)
	}

	return (*BigInt)(new(big.Int).Exp((*big.Int)(value), (*big.Int)(b), (*big.Int)(m))).MaybeInt(), nil
}

func (value *BigInt) __pow__(other, modulus Object) (Object, error) {
	var m *BigInt
	if modulus != Nil {
		var ok bool
		if m, ok = ConvertToBigInt(modulus); !ok {
			return NotImplemented, nil
		}
	}

	if b, ok := ConvertToBigInt(other); ok {
		return value.pow(b, m)
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_pow__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return b.pow(value, nil)
	}

	return NotImplemented, nil
}

func (value *BigInt) __in_place_pow__(other, modulus Object) (Object, error) {
	return value.__pow__(other, modulus)
}

func (value *BigInt) __left_shift__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		bb, err := b.GoInt()
		if err != nil {
			return nil, err
		}

		if bb < 0 {
			return nil, negativeShiftCount
		}

		return (*BigInt)(new(big.Int).Lsh((*big.Int)(value), uint(bb))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_left_shift__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		aa, err := value.GoInt()
		if err != nil {
			return nil, err
		}

		if aa < 0 {
			return nil, negativeShiftCount
		}

		return (*BigInt)(new(big.Int).Lsh((*big.Int)(b), uint(aa))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __in_place_left_shift__(other Object) (Object, error) {
	return value.__left_shift__(other)
}

func (value *BigInt) __right_shift__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		bb, err := b.GoInt()
		if err != nil {
			return nil, err
		}

		if bb < 0 {
			return nil, negativeShiftCount
		}

		return (*BigInt)(new(big.Int).Rsh((*big.Int)(value), uint(bb))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_right_shift__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		aa, err := value.GoInt()
		if err != nil {
			return nil, err
		}

		if aa < 0 {
			return nil, negativeShiftCount
		}

		return (*BigInt)(new(big.Int).Rsh((*big.Int)(b), uint(aa))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __in_place_right_shift__(other Object) (Object, error) {
	return value.__right_shift__(other)
}

func (value *BigInt) __and__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).And((*big.Int)(value), (*big.Int)(b))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_and__(other Object) (Object, error) {
	return value.__and__(other)
}

func (value *BigInt) __in_place_and__(other Object) (Object, error) {
	return value.__and__(other)
}

func (value *BigInt) __xor__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Xor((*big.Int)(value), (*big.Int)(b))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_xor__(other Object) (Object, error) {
	return value.__xor__(other)
}

func (value *BigInt) __in_place_xor__(other Object) (Object, error) {
	return value.__xor__(other)
}

func (value *BigInt) __or__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return (*BigInt)(new(big.Int).Or((*big.Int)(value), (*big.Int)(b))).MaybeInt(), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __reversed_or__(other Object) (Object, error) {
	return value.__or__(other)
}

func (value *BigInt) __in_place_or__(other Object) (Object, error) {
	return value.__or__(other)
}

func (value *BigInt) __bool__() (Object, error) {
	return NewBool((*big.Int)(value).Sign() != 0), nil
}

func (value *BigInt) __index__() (Int, error) {
	return value.Int()
}

func (value *BigInt) __int__() (Object, error) {
	return value, nil
}

func (value *BigInt) __real__() (Object, error) {
	return value.Real()
}

// Comparison

func (value *BigInt) __less_than__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return NewBool((*big.Int)(value).Cmp((*big.Int)(b)) < 0), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __less_or_equal__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return NewBool((*big.Int)(value).Cmp((*big.Int)(b)) <= 0), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __equal__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return NewBool((*big.Int)(value).Cmp((*big.Int)(b)) == 0), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __not_equal__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return NewBool((*big.Int)(value).Cmp((*big.Int)(b)) != 0), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __greater_than__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return NewBool((*big.Int)(value).Cmp((*big.Int)(b)) > 0), nil
	}

	return NotImplemented, nil
}

func (value *BigInt) __greater_or_equal__(other Object) (Object, error) {
	if b, ok := ConvertToBigInt(other); ok {
		return NewBool((*big.Int)(value).Cmp((*big.Int)(b)) >= 0), nil
	}

	return NotImplemented, nil
}

// Check interface is satisfied
var _ Object = (*BigInt)(nil)
var _ realArithmetic = (*BigInt)(nil)
var _ booleanArithmetic = (*BigInt)(nil)
var _ conversionBetweenTypes = (*BigInt)(nil)
var _ I__bool__ = (*BigInt)(nil)
var _ I__index__ = (*BigInt)(nil)
var _ iComparison = (*BigInt)(nil)
var _ IGoInt = (*BigInt)(nil)
var _ IGoInt64 = (*BigInt)(nil)
