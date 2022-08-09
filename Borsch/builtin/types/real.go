package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var RealClass = ObjectClass.ClassNew("дійсний", map[string]Object{}, true, RealNew, nil)

type Real float64

func (value Real) Class() *Class {
	return RealClass
}

func RealNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	var xObj Object = Real(0)
	err := parseArgs(cls.Name, "o|!", args, 1, 1, &xObj)
	if err != nil {
		return nil, err
	}

	switch x := xObj.(type) {
	case String:
		return RealFromString(string(x))
	}

	return ToReal(ctx, xObj)
}

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

		return nil, ErrorNewf("invalid literal for real: '%s'", str)
	}

	return Real(f), nil
}

func (value Real) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value Real) string(Context) (Object, error) {
	if i := int64(value); Real(i) == value {
		return String(fmt.Sprintf("%d.0", i)), nil
	}

	return String(fmt.Sprintf("%g", value)), nil
}

func (value Real) toInt(ctx Context) (Object, error) {
	return Int(value), nil
}

func (value Real) add(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return value + otherValue, nil
	}

	if otherValue, ok := other.(Int); ok {
		return value + Real(otherValue), nil
	}

	return nil, ErrorNewf("неможливо виконати додавання дійсного числа до об'єкта '%s'", other.Class().Name)
}

func (value Real) reversedAdd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return otherValue + value, nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(otherValue) + value, nil
	}

	return nil, ErrorNewf("неможливо виконати додавання об'єкта '%s' до дійсне число", other.Class().Name)
}

func (value Real) sub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return value - otherValue, nil
	}

	if otherValue, ok := other.(Int); ok {
		return value - Real(otherValue), nil
	}

	return nil, ErrorNewf("неможливо виконати віднімання дійсного числа від об'єкта '%s'", other.Class().Name)
}

func (value Real) reversedSub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return otherValue - value, nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(otherValue) - value, nil
	}

	return nil, ErrorNewf("неможливо виконати віднімання об'єкта '%s' від дійсне число", other.Class().Name)
}

func (value Real) div(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		if otherValue == 0 {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return value / otherValue, nil
	}

	if otherValue, ok := other.(Int); ok {
		if otherValue == 0 {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return value / Real(otherValue), nil
	}

	return nil, ErrorNewf("неможливо виконати ділення дійсного числа на об'єкт '%s'", other.Class().Name)
}

func (value Real) reversedDiv(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		if value == 0 {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return otherValue / value, nil
	}

	if otherValue, ok := other.(Int); ok {
		if value == 0 {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return Real(otherValue) / value, nil
	}

	return nil, ErrorNewf("неможливо виконати ділення об'єкта '%s' на дійсне число", other.Class().Name)
}

func (value Real) mul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return value * otherValue, nil
	}

	if otherValue, ok := other.(Int); ok {
		return value * Real(otherValue), nil
	}

	return nil, ErrorNewf("неможливо виконати множення дійсного числа на об'єкт '%s'", other.Class().Name)
}

func (value Real) reversedMul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return otherValue * value, nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(otherValue) * value, nil
	}

	return nil, ErrorNewf("неможливо виконати множення об'єкта '%s' на дійсне число", other.Class().Name)
}

func (value Real) mod(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return Real(math.Mod(float64(value), float64(otherValue))), nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(math.Mod(float64(value), float64(otherValue))), nil
	}

	return nil, ErrorNewf("неможливо виконати модуль? дійсного числа  '%s'", other.Class().Name)
}

func (value Real) reversedMod(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return Real(math.Mod(float64(otherValue), float64(value))), nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(math.Mod(float64(otherValue), float64(value))), nil
	}

	return nil, ErrorNewf("неможливо виконати модуль? об'єкта '%s'  дійсне число", other.Class().Name)
}

func (value Real) pow(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return Real(math.Pow(float64(value), float64(otherValue))), nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(math.Pow(float64(value), float64(otherValue))), nil
	}

	return nil, ErrorNewf("неможливо виконати степінь? дійсного числа  '%s'", other.Class().Name)
}

func (value Real) reversedPow(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return Real(math.Pow(float64(otherValue), float64(value))), nil
	}

	if otherValue, ok := other.(Int); ok {
		result := math.Pow(float64(otherValue), float64(value))
		if otherValue < 0 {
			return Real(result), nil
		}

		return Int(result), nil
	}

	return nil, ErrorNewf("неможливо виконати степінь? об'єкта '%s'  дійсне число", other.Class().Name)
}

func (value Real) equals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value == v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value == Real(v)), nil
	}

	return False, nil
}

func (value Real) notEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value != v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value != Real(v)), nil
	}

	return False, nil
}

func (value Real) less(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value < v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value < Real(v)), nil
	}

	return False, nil
}

func (value Real) lessOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value <= v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value <= Real(v)), nil
	}

	return False, nil
}

func (value Real) greater(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value > v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value > Real(v)), nil
	}

	return False, nil
}

func (value Real) greaterOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value >= v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value >= Real(v)), nil
	}

	return False, nil
}

func (value Real) positive(_ Context) (Object, error) {
	return +value, nil
}

func (value Real) negate(_ Context) (Object, error) {
	return -value, nil
}
