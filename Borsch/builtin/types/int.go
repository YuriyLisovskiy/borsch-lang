package types

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var IntClass = ObjectClass.ClassNew("ціле", map[string]Object{}, true, IntNew, nil)

type Int int64

func io2bo(value Int) Bool {
	return value != 0
}

func (value Int) Class() *Class {
	return IntClass
}

func IntNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	var xObj Object = Int(0)
	var baseObj Object
	base := 0
	err := parseArgs(cls.Name, "oi|!!", args, 1, 2, &xObj, &baseObj)
	if err != nil {
		return nil, err
	}

	if baseObj != nil {
		base, err = ToGoInt(ctx, baseObj)
		if err != nil {
			return nil, err
		}

		if base != 0 && (base < 2 || base > 36) {
			return nil, NewErrorf("база для 'ціле()' має бути >= 2 та <= 36")
		}
	}

	if x, ok := xObj.(String); ok {
		return IntFromString(string(x), base)
	}

	if baseObj != nil {
		// TODO: TypeError
		return nil, NewErrorf("'ціле()' не може перетворити об'єкт не рядкового типу з явною базою")
	}

	return ToInt(ctx, xObj)
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

	// TODO:
	// return (*BigInt)(x).MaybeInt(), nil
	return nil, NewErrorf("overflow...")

error:
	// TODO: ValueError
	return nil, NewErrorf("некоректний літерал для 'ціле()' з базою %d: '%s'", convertBase, str)
}

func (value Int) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value Int) string(Context) (Object, error) {
	return String(fmt.Sprintf("%d", value)), nil
}

func (value Int) toBool(Context) (Object, error) {
	return Bool(value != 0), nil
}

func (value Int) toReal(Context) (Object, error) {
	return Real(value), nil
}

// func (value Int) toInt(Context) (Object, error) {
// 	return value, nil
// }

// GoInt truncates to Go int.
//
// If it is outside the range of Go int it will return an error.
func (value Int) toGoInt(Context) (int, error) {
	r := int(value)
	if Int(r) != value {
		return 0, overflowErrorGo
	}

	return r, nil
}

func (value Int) add(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value + otherValue, nil
	}

	if otherValue, ok := other.(Real); ok {
		return Real(value) + otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value + bo2io(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати додавання цілого числа до об'єкта '%s'", other.Class().Name)
}

func (value Int) reversedAdd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue + value, nil
	}

	if otherValue, ok := other.(Real); ok {
		return otherValue + Real(value), nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) + value, nil
	}

	return nil, NewErrorf("неможливо виконати додавання об'єкта '%s' до ціле число", other.Class().Name)
}

func (value Int) sub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value - otherValue, nil
	}

	if otherValue, ok := other.(Real); ok {
		return Real(value) - otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value - bo2io(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати віднімання цілого числа від об'єкта '%s'", other.Class().Name)
}

func (value Int) reversedSub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue - value, nil
	}

	if otherValue, ok := other.(Real); ok {
		return otherValue - Real(value), nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) - value, nil
	}

	return nil, NewErrorf("неможливо виконати віднімання об'єкта '%s' від цілого числа", other.Class().Name)
}

func (value Int) div(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		if otherValue == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return Real(value) / Real(otherValue), nil
	}

	if otherValue, ok := other.(Real); ok {
		if otherValue == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return Real(value) / otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return value, nil
	}

	return nil, NewErrorf("неможливо виконати ділення цілого числа на об'єкт '%s'", other.Class().Name)
}

func (value Int) reversedDiv(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return Real(otherValue) / Real(value), nil
	}

	if otherValue, ok := other.(Real); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return otherValue / Real(value), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		if !otherValue {
			return Real(0.0), nil
		}

		return 1.0 / Real(value), nil
	}

	return nil, NewErrorf("неможливо виконати ділення об'єкта '%s' на ціле число", other.Class().Name)
}

func (value Int) mul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value * otherValue, nil
	}

	if otherValue, ok := other.(Real); ok {
		return Real(value) * otherValue, nil
	}

	if otherValue, ok := other.(String); ok {
		result := String("")
		if value <= 0 {
			return result, nil
		}

		for i := int64(0); i < int64(value); i++ {
			result += otherValue
		}

		return result, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value * bo2io(otherValue), nil
	}

	// TODO: add multiplication for:
	//  int, ...

	return nil, NewErrorf("неможливо виконати множення цілого числа на об'єкт '%s'", other.Class().Name)
}

func (value Int) reversedMul(ctx Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue * value, nil
	}

	if otherValue, ok := other.(Real); ok {
		return otherValue * Real(value), nil
	}

	if otherValue, ok := other.(String); ok {
		return otherValue.mul(ctx, value)
	}

	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return Int(0), nil
		}

		return value, nil
	}

	// TODO: add multiplication for:
	//  ..., int

	return nil, NewErrorf("неможливо виконати множення об'єкта '%s' на ціле число", other.Class().Name)
}

func (value Int) mod(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		if otherValue == 0 {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return value % otherValue, nil
	}

	if otherValue, ok := other.(Real); ok {
		if otherValue == 0 {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return Real(math.Mod(float64(value), float64(otherValue))), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return Int(mod(Real(value), bo2ro(otherValue))), nil
	}

	return nil, NewErrorf("неможливо виконати модуль? цілого числа  '%s'", other.Class().Name)
}

func (value Int) reversedMod(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return Int(mod(Real(otherValue), Real(value))), nil
	}

	if otherValue, ok := other.(Real); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return mod(otherValue, Real(value)), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return Int(mod(bo2ro(otherValue), Real(value))), nil
	}

	return nil, NewErrorf("неможливо виконати модуль? об'єкта '%s'  ціле число", other.Class().Name)
}

func (value Int) pow(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		result := math.Pow(float64(value), float64(otherValue))
		if otherValue < 0 {
			return Real(result), nil
		}

		return Int(result), nil
	}

	if otherValue, ok := other.(Real); ok {
		return Real(math.Pow(float64(value), float64(otherValue))), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if otherValue {
			return value, nil
		}

		return Int(1), nil
	}

	// TODO:
	return nil, NewErrorf("неможливо виконати ??? '%s'", other.Class().Name)
}

func (value Int) reversedPow(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		result := math.Pow(float64(otherValue), float64(value))
		if otherValue < 0 {
			return Real(result), nil
		}

		return Int(result), nil
	}

	if otherValue, ok := other.(Real); ok {
		return Real(math.Pow(float64(otherValue), float64(value))), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if otherValue {
			if value < 0 {
				return Real(1.0), nil
			}

			return Int(1), nil
		}

		if value < 0.0 {
			return nil, NewZeroDivisionError("неможливо піднести 0.0 до від'ємного степеня")
		}

		if value == 0 {
			return Int(1), nil
		}

		return Int(0), nil
	}

	// TODO:
	return nil, NewErrorf("неможливо виконати ??? '%s' ???", other.Class().Name)
}

func (value Int) equals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Int); ok {
		return gb2bo(value == v), nil
	}

	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(Real(value) == v), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value == bo2io(v)), nil
	}

	return False, nil
}

func (value Int) notEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value != v), nil
	}

	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(Real(value) != v), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value != bo2io(v)), nil
	}

	return False, nil
}

func (value Int) less(_ Context, other Object) (Object, error) {
	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value < v), nil
	}

	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(Real(value) < v), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value < bo2io(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew("<", value.Class().Name, other.Class().Name)
}

func (value Int) lessOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value <= v), nil
	}

	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(Real(value) <= v), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value <= bo2io(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew("<=", value.Class().Name, other.Class().Name)
}

func (value Int) greater(_ Context, other Object) (Object, error) {
	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value > v), nil
	}

	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(Real(value) > v), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value > bo2io(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew(">", value.Class().Name, other.Class().Name)
}

func (value Int) greaterOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value >= v), nil
	}

	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(Real(value) >= v), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value >= bo2io(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew(">=", value.Class().Name, other.Class().Name)
}

func (value Int) shiftLeft(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value << otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value << bo2io(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати побітовий зсув ліворуч цілого числа на об'єкт '%s'", other.Class().Name)
}

func (value Int) reversedShiftLeft(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue << value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) << value, nil
	}

	return nil, NewErrorf("неможливо виконати побітовий зсув ліворуч об'єкта '%s' на ціле число", other.Class().Name)
}

func (value Int) shiftRight(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value >> otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value >> bo2io(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати побітовий зсув праворуч цілого числа на об'єкт '%s'", other.Class().Name)
}

func (value Int) reversedShiftRight(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue >> value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) >> value, nil
	}

	return nil, NewErrorf("неможливо виконати побітовий зсув праворуч об'єкта '%s' на ціле число", other.Class().Name)
}

func (value Int) bitwiseOr(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value | otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value | bo2io(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати побітову диз'юнкцію цілого числа та об'єкта '%s'", other.Class().Name)
}

func (value Int) reversedBitwiseOr(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue | value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) | value, nil
	}

	return nil, NewErrorf("неможливо виконати побітову диз'юнкцію об'єкта '%s' та ціле число", other.Class().Name)
}

func (value Int) bitwiseXor(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value ^ otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value ^ bo2io(otherValue), nil
	}

	return nil, NewErrorf(
		"неможливо виконати побітову виняткову диз'юнкцію цілого числа та об'єкта '%s'",
		other.Class().Name,
	)
}

func (value Int) reversedBitwiseXor(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue ^ value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) ^ value, nil
	}

	return nil, NewErrorf(
		"неможливо виконати побітову виняткову диз'юнкцію об'єкта '%s' та ціле число",
		other.Class().Name,
	)
}

func (value Int) bitwiseAnd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return value & otherValue, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value & bo2io(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати побітову кон'юнкцію цілого числа та об'єкта '%s'", other.Class().Name)
}

func (value Int) reversedBitwiseAnd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue & value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) & value, nil
	}

	return nil, NewErrorf("неможливо виконати побітову кон'юнкцію об'єкта '%s' та ціле число", other.Class().Name)
}

func (value Int) positive(_ Context) (Object, error) {
	return +value, nil
}

func (value Int) negate(_ Context) (Object, error) {
	return -value, nil
}

func (value Int) invert(_ Context) (Object, error) {
	return ^value, nil
}
