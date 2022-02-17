// Automatically generated - DO NOT EDIT
// Regenerate with: go generate

// Arithmetic operations

package types

// Negate the Object returning an Object.
//
// Will raise TypeError if Negate can't be run on this object.
func Negate(a Object) (Object, error) {

	if A, ok := a.(I__neg__); ok {
		res, err := A.__neg__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримуваний тип операнда для -: '%s'", a.Type().Name)
}

// MakePositive the Object returning an Object.
//
// Will raise TypeError if MakePositive can't be run on this object.
func MakePositive(a Object) (Object, error) {

	if A, ok := a.(I__pos__); ok {
		res, err := A.__pos__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримуваний тип операнда для +: '%s'", a.Type().Name)
}

// Invert the Object returning an Object.
//
// Will raise TypeError if Invert can't be run on this object.
func Invert(a Object) (Object, error) {

	if A, ok := a.(I__invert__); ok {
		res, err := A.__invert__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримуваний тип операнда для ~: '%s'", a.Type().Name)
}

// MakeInt the Object returning an Object.
//
// Will raise TypeError if MakeInt can't be run on this object.
func MakeInt(a Object) (Object, error) {

	if _, ok := a.(Int); ok {
		return a, nil
	}

	if A, ok := a.(I__int__); ok {
		res, err := A.__int__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримуваний тип операнда для int: '%s'", a.Type().Name)
}

// MakeReal the Object returning an Object.
//
// Will raise TypeError if MakeReal can't be run on this object.
func MakeReal(a Object) (Object, error) {

	if _, ok := a.(Real); ok {
		return a, nil
	}

	if A, ok := a.(I__real__); ok {
		res, err := A.__real__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримуваний тип операнда для real: '%s'", a.Type().Name)
}

// Iter the Object returning an Object.
//
// Will raise TypeError if Iter can't be run on this object.
func Iter(a Object) (Object, error) {

	if A, ok := a.(I__iter__); ok {
		res, err := A.__iter__()
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримуваний тип операнда для iter: '%s'", a.Type().Name)
}

// Add two objects together returning an Object.
//
// Will raise TypeError if add can't be run on these objects.
func Add(a, b Object) (Object, error) {
	// Try using a to add
	if A, ok := a.(I__add__); ok {
		res, err := A.__add__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_add if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_add__); ok {
			res, err := B.__reversed_add__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для +: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceAdd(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_add__); ok {
		res, err := A.__in_place_add__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Add(a, b)
}

// Subtract two objects together returning an Object.
//
// Will raise TypeError if sub can't be run on these objects.
func Subtract(a, b Object) (Object, error) {
	// Try using a to sub
	if A, ok := a.(I__sub__); ok {
		res, err := A.__sub__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_sub if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_sub__); ok {
			res, err := B.__reversed_sub__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для -: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceSubtract(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_sub__); ok {
		res, err := A.__in_place_sub__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Subtract(a, b)
}

// Multiply two objects together returning an Object.
//
// Will raise TypeError if mul can't be run on these objects.
func Multiply(a, b Object) (Object, error) {
	// Try using a to mul
	if A, ok := a.(I__mul__); ok {
		res, err := A.__mul__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_mul if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_mul__); ok {
			res, err := B.__reversed_mul__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для *: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceMultiply(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_mul__); ok {
		res, err := A.__in_place_mul__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Multiply(a, b)
}

// Divide two objects together returning an Object.
//
// Will raise TypeError if div can't be run on these objects.
func Divide(a, b Object) (Object, error) {
	// Try using a to div
	if A, ok := a.(I__div__); ok {
		res, err := A.__div__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_div if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_div__); ok {
			res, err := B.__reversed_div__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для /: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceDivide(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_div__); ok {
		res, err := A.__in_place_div__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Divide(a, b)
}

// Mod two objects together returning an Object.
//
// Will raise TypeError if mod can't be run on these objects.
func Mod(a, b Object) (Object, error) {
	// Try using a to mod
	if A, ok := a.(I__mod__); ok {
		res, err := A.__mod__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_mod if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_mod__); ok {
			res, err := B.__reversed_mod__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для %%: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceMod(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_mod__); ok {
		res, err := A.__in_place_mod__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Mod(a, b)
}

// LeftShift two objects together returning an Object.
//
// Will raise TypeError if left_shift can't be run on these objects.
func LeftShift(a, b Object) (Object, error) {
	// Try using a to left_shift
	if A, ok := a.(I__left_shift__); ok {
		res, err := A.__left_shift__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_left_shift if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_left_shift__); ok {
			res, err := B.__reversed_left_shift__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для <<: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceLeftShift(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_left_shift__); ok {
		res, err := A.__in_place_left_shift__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return LeftShift(a, b)
}

// RightShift two objects together returning an Object.
//
// Will raise TypeError if right_shift can't be run on these objects.
func RightShift(a, b Object) (Object, error) {
	// Try using a to right_shift
	if A, ok := a.(I__right_shift__); ok {
		res, err := A.__right_shift__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_right_shift if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_right_shift__); ok {
			res, err := B.__reversed_right_shift__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для >>: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceRightShift(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_right_shift__); ok {
		res, err := A.__in_place_right_shift__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return RightShift(a, b)
}

// And two objects together returning an Object.
//
// Will raise TypeError if and can't be run on these objects.
func And(a, b Object) (Object, error) {
	// Try using a to and
	if A, ok := a.(I__and__); ok {
		res, err := A.__and__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_and if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_and__); ok {
			res, err := B.__reversed_and__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для &: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceAnd(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_and__); ok {
		res, err := A.__in_place_and__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return And(a, b)
}

// Xor two objects together returning an Object.
//
// Will raise TypeError if xor can't be run on these objects.
func Xor(a, b Object) (Object, error) {
	// Try using a to xor
	if A, ok := a.(I__xor__); ok {
		res, err := A.__xor__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_xor if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_xor__); ok {
			res, err := B.__reversed_xor__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для ^: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceXor(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_xor__); ok {
		res, err := A.__in_place_xor__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Xor(a, b)
}

// Or two objects together returning an Object.
//
// Will raise TypeError if or can't be run on these objects.
func Or(a, b Object) (Object, error) {
	// Try using a to or
	if A, ok := a.(I__or__); ok {
		res, err := A.__or__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_or if different in type to a
	if a.Type() != b.Type() {
		if B, ok := b.(I__reversed_or__); ok {
			res, err := B.__reversed_or__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для |: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlaceOr(a, b Object) (Object, error) {
	if A, ok := a.(I__in_place_or__); ok {
		res, err := A.__in_place_or__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Or(a, b)
}

// Pow three objects together returning an Object.
//
// If c != Nil then it won't attempt to call __reversed_pow__
//
// Will raise TypeError if pow can't be run on these objects.
func Pow(a, b, c Object) (Object, error) {
	// Try using a to pow
	if A, ok := a.(I__pow__); ok {
		res, err := A.__pow__(b, c)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Now using b to reversed_pow if different in type to a
	if c == Nil && a.Type() != b.Type() {
		if B, ok := b.(I__reversed_pow__); ok {
			res, err := B.__reversed_pow__(a)
			if err != nil {
				return nil, err
			}

			if res != NotImplemented {
				return res, nil
			}
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для **: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

func InPlacePow(a, b, c Object) (Object, error) {
	if A, ok := a.(I__in_place_pow__); ok {
		res, err := A.__in_place_pow__(b, c)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	return Pow(a, b, c)
}

// Greater two objects returning a boolean result.
//
// Will raise TypeError if Greater can't be run on this object.
func Greater(a Object, b Object) (Object, error) {
	// Try using a to greater_than.
	if A, ok := a.(I__greater_than__); ok {
		res, err := A.__greater_than__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to less_than with reversed parameters.
	if B, ok := b.(I__less_than__); ok {
		res, err := B.__less_than__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для >: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

// GreaterOrEqual two objects returning a boolean result.
//
// Will raise TypeError if GreaterOrEqual can't be run on this object.
func GreaterOrEqual(a Object, b Object) (Object, error) {
	// Try using a to greater_or_equal.
	if A, ok := a.(I__greater_or_equal__); ok {
		res, err := A.__greater_or_equal__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to less_or_equal with reversed parameters.
	if B, ok := b.(I__less_or_equal__); ok {
		res, err := B.__less_or_equal__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для >=: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

// LessThan two objects returning a boolean result.
//
// Will raise TypeError if LessThan can't be run on this object.
func LessThan(a Object, b Object) (Object, error) {
	// Try using a to less_than.
	if A, ok := a.(I__less_than__); ok {
		res, err := A.__less_than__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to greater_than with reversed parameters.
	if B, ok := b.(I__greater_than__); ok {
		res, err := B.__greater_than__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для <: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

// LessOrEqual two objects returning a boolean result.
//
// Will raise TypeError if LessOrEqual can't be run on this object.
func LessOrEqual(a Object, b Object) (Object, error) {
	// Try using a to less_or_equal.
	if A, ok := a.(I__less_or_equal__); ok {
		res, err := A.__less_or_equal__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to greater_or_equal with reversed parameters.
	if B, ok := b.(I__greater_or_equal__); ok {
		res, err := B.__greater_or_equal__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для <=: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

// Equal two objects returning a boolean result.
//
// Will raise TypeError if Equal can't be run on this object.
func Equal(a Object, b Object) (Object, error) {
	// Try using a to equal.
	if A, ok := a.(I__equal__); ok {
		res, err := A.__equal__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to equal with reversed parameters.
	if B, ok := b.(I__equal__); ok {
		res, err := B.__equal__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	if a.Type() != b.Type() {
		return False, nil
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для ==: '%s' та '%s'", a.Type().Name, b.Type().Name)
}

// NotEqual two objects returning a boolean result.
//
// Will raise TypeError if NotEqual can't be run on this object.
func NotEqual(a Object, b Object) (Object, error) {
	// Try using a to not_equal.
	if A, ok := a.(I__not_equal__); ok {
		res, err := A.__not_equal__(b)
		if err != nil {
			return nil, err
		}

		if res != NotImplemented {
			return res, nil
		}
	}

	// Try using b to not_equal with reversed parameters.
	if B, ok := b.(I__not_equal__); ok {
		res, err := B.__not_equal__(a)
		if err != nil {
			return nil, err
		}
		if res != NotImplemented {
			return res, nil
		}
	}

	if a.Type() != b.Type() {
		return True, nil
	}

	return nil, ErrorNewf(TypeError, "непідтримувані типи операндів для !=: '%s' та '%s'", a.Type().Name, b.Type().Name)
}
