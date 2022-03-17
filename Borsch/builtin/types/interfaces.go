package types

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/common"

// StringClass interfaces.

type I__represent__ interface {
	__represent__() (common.Value, error)
}

type I__str__ interface {
	__str__() (common.Value, error)
}

// Comparison operators.

type I__less_than__ interface {
	__less_than__(other common.Value) (common.Value, error)
}

type I__less_or_equal__ interface {
	__less_or_equal__(other common.Value) (common.Value, error)
}

type I__equal__ interface {
	__equal__(other common.Value) (common.Value, error)
}

type I__not_equal__ interface {
	__not_equal__(other common.Value) (common.Value, error)
}

type I__greater_than__ interface {
	__greater_than__(other common.Value) (common.Value, error)
}

type I__greater_or_equal__ interface {
	__greater_or_equal__(other common.Value) (common.Value, error)
}

type iComparison interface {
	I__less_than__
	I__less_or_equal__
	I__equal__
	I__not_equal__
	I__greater_than__
	I__greater_or_equal__
}

type I__bool__ interface {
	__bool__() (common.Value, error)
}

type I__get_attribute__ interface {
	__get_attribute__(ctx common.Context, name string) (common.Value, error)
}

type I__set_attribute__ interface {
	__set_attribute__(ctx common.Context, name string, value common.Value) error
}

type I__get__ interface {
	__get__(instance, owner common.Value) (common.Value, error)
}

type I__length__ interface {
	__length__() (common.Value, error)
}

type I__call__ interface {
	__call__(state common.State, args Tuple, kwargs map[string]common.Value) (common.Value, error)
}

type I__get_item__ interface {
	__get_item__(key common.Value) (common.Value, error)
}

type I__set_item__ interface {
	__set_item__(key, value common.Value) (common.Value, error)
}

type I__delete_item__ interface {
	__delete_item__(key common.Value) (common.Value, error)
}

type I__iter__ interface {
	__iter__() (common.Value, error)
}

type I__next__ interface {
	__next__() (common.Value, error)
}

type I_iterator interface {
	I__iter__
	I__next__
}

// I__contains__ called to implement membership test operators.
// Should return true if item is in self, false otherwise.
// For mapping objects, this should consider the keys of the
// mapping rather than the values or the key-item pairs.
type I__contains__ interface {
	__contains__(item common.Value) (common.Value, error)
}

// Arithmetic operators.

type I__add__ interface {
	__add__(ctx common.Context, other common.Value) (common.Value, error)
}

type I__sub__ interface {
	__sub__(other common.Value) (common.Value, error)
}

type I__mul__ interface {
	__mul__(other common.Value) (common.Value, error)
}

type I__div__ interface {
	__div__(other common.Value) (common.Value, error)
}

type I__mod__ interface {
	__mod__(other common.Value) (common.Value, error)
}

type I__pow__ interface {
	__pow__(other, modulo common.Value) (common.Value, error)
}

type I__left_shift__ interface {
	__left_shift__(other common.Value) (common.Value, error)
}

type I__right_shift__ interface {
	__right_shift__(other common.Value) (common.Value, error)
}

type I__and__ interface {
	__and__(other common.Value) (common.Value, error)
}

type I__xor__ interface {
	__xor__(other common.Value) (common.Value, error)
}

type I__or__ interface {
	__or__(other common.Value) (common.Value, error)
}

// Reversed arithmetic operators.

type I__reversed_add__ interface {
	__reversed_add__(other common.Value) (common.Value, error)
}

type I__reversed_sub__ interface {
	__reversed_sub__(other common.Value) (common.Value, error)
}

type I__reversed_mul__ interface {
	__reversed_mul__(other common.Value) (common.Value, error)
}

type I__reversed_div__ interface {
	__reversed_div__(other common.Value) (common.Value, error)
}

type I__reversed_mod__ interface {
	__reversed_mod__(other common.Value) (common.Value, error)
}

type I__reversed_pow__ interface {
	__reversed_pow__(other common.Value) (common.Value, error)
}

type I__reversed_left_shift__ interface {
	__reversed_left_shift__(other common.Value) (common.Value, error)
}

type I__reversed_right_shift__ interface {
	__reversed_right_shift__(other common.Value) (common.Value, error)
}

type I__reversed_and__ interface {
	__reversed_and__(other common.Value) (common.Value, error)
}

type I__reversed_xor__ interface {
	__reversed_xor__(other common.Value) (common.Value, error)
}

type I__reversed_or__ interface {
	__reversed_or__(other common.Value) (common.Value, error)
}

// In-place arithmetic operators.

type I__in_place_add__ interface {
	__in_place_add__(other common.Value) (common.Value, error)
}

type I__in_place_sub__ interface {
	__in_place_sub__(other common.Value) (common.Value, error)
}

type I__in_place_mul__ interface {
	__in_place_mul__(other common.Value) (common.Value, error)
}

type I__in_place_div__ interface {
	__in_place_div__(other common.Value) (common.Value, error)
}

type I__in_place_mod__ interface {
	__in_place_mod__(other common.Value) (common.Value, error)
}

type I__in_place_pow__ interface {
	__in_place_pow__(other, modulo common.Value) (common.Value, error)
}

type I__in_place_left_shift__ interface {
	__in_place_left_shift__(other common.Value) (common.Value, error)
}

type I__in_place_right_shift__ interface {
	__in_place_right_shift__(other common.Value) (common.Value, error)
}

type I__in_place_and__ interface {
	__in_place_and__(other common.Value) (common.Value, error)
}

type I__in_place_xor__ interface {
	__in_place_xor__(other common.Value) (common.Value, error)
}

type I__in_place_or__ interface {
	__in_place_or__(other common.Value) (common.Value, error)
}

// Called to implement the unary arithmetic operations (-, + and ~).

type I__neg__ interface {
	__neg__() (common.Value, error)
}

type I__pos__ interface {
	__pos__() (common.Value, error)
}

type I__invert__ interface {
	__invert__() (common.Value, error)
}

// Called to implement the built-in functions int() and float().
// Should return a value of the appropriate type.

type I__int__ interface {
	__int__() (common.Value, error)
}

type I__real__ interface {
	__real__() (common.Value, error)
}

type I__index__ interface {
	__index__() (Int, error)
}

// Int and Real should satisfy this.
type realArithmetic interface {
	I__neg__
	I__pos__
	I__add__
	I__sub__
	I__mul__
	I__div__
	I__mod__
	I__pow__
	I__reversed_add__
	I__reversed_sub__
	I__reversed_mul__
	I__reversed_div__
	I__reversed_mod__
	I__reversed_pow__
	I__in_place_add__
	I__in_place_sub__
	I__in_place_mul__
	I__in_place_div__
	I__in_place_mod__
	I__in_place_pow__
}

// Int should satisfy this
type booleanArithmetic interface {
	I__invert__
	I__left_shift__
	I__right_shift__
	I__and__
	I__xor__
	I__or__
	I__reversed_left_shift__
	I__reversed_right_shift__
	I__reversed_and__
	I__reversed_xor__
	I__reversed_or__
	I__in_place_left_shift__
	I__in_place_right_shift__
	I__in_place_and__
	I__in_place_xor__
	I__in_place_or__
}

// Real and Int should satisfy this.
type conversionBetweenTypes interface {
	I__int__
	I__real__
}

// StringClass, Tuple, List should satisfy this.
type sequenceArithmetic interface {
	I__add__
	I__mul__
	I__reversed_add__
	I__reversed_mul__
	I__in_place_add__
	I__in_place_mul__
}
