package types

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var RealClass = ObjectClass.ClassNew("дійсне", map[string]Object{}, true, RealNew, nil)

type Real float64

func bo2ro(value Bool) Real {
	if value {
		return 1.0
	}

	return 0.0
}

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

		return nil, NewErrorf("invalid literal for real: '%s'", str)
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

func (value Real) toBool(Context) (Object, error) {
	return Bool(value != 0.0), nil
}

// func (value Real) toReal(Context) (Object, error) {
// 	return value, nil
// }

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

	if otherValue, ok := other.(Bool); ok {
		return value + bo2ro(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати додавання дійсного числа до об'єкта '%s'", other.Class().Name)
}

func (value Real) reversedAdd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return otherValue + value, nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(otherValue) + value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2ro(otherValue) + value, nil
	}

	return nil, NewErrorf("неможливо виконати додавання об'єкта '%s' до дійсне число", other.Class().Name)
}

func (value Real) sub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return value - otherValue, nil
	}

	if otherValue, ok := other.(Int); ok {
		return value - Real(otherValue), nil
	}

	if otherValue, ok := other.(Bool); ok {
		return value - bo2ro(otherValue), nil
	}

	return nil, NewErrorf("неможливо виконати віднімання дійсного числа від об'єкта '%s'", other.Class().Name)
}

func (value Real) reversedSub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return otherValue - value, nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(otherValue) - value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		return bo2ro(otherValue) - value, nil
	}

	return nil, NewErrorf("неможливо виконати віднімання об'єкта '%s' від дійсне число", other.Class().Name)
}

func (value Real) div(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		if otherValue == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return value / otherValue, nil
	}

	if otherValue, ok := other.(Int); ok {
		if otherValue == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return value / Real(otherValue), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return value, nil
	}

	return nil, NewErrorf("неможливо виконати ділення дійсного числа на об'єкт '%s'", other.Class().Name)
}

func (value Real) reversedDiv(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return otherValue / value, nil
	}

	if otherValue, ok := other.(Int); ok {
		if value == 0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		return Real(otherValue) / value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		if value == 0.0 {
			return nil, NewZeroDivisionError("ділення на нуль")
		}

		if !otherValue {
			return Real(0.0), nil
		}

		return 1.0 / value, nil
	}

	return nil, NewErrorf("неможливо виконати ділення об'єкта '%s' на дійсне число", other.Class().Name)
}

func (value Real) mul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return value * otherValue, nil
	}

	if otherValue, ok := other.(Int); ok {
		return value * Real(otherValue), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if otherValue {
			return value, nil
		}

		return Real(0.0), nil
	}

	return nil, NewErrorf("неможливо виконати множення дійсного числа на об'єкт '%s'", other.Class().Name)
}

func (value Real) reversedMul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return otherValue * value, nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(otherValue) * value, nil
	}

	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return Real(0.0), nil
		}

		return value, nil
	}

	return nil, NewErrorf("неможливо виконати множення об'єкта '%s' на дійсне число", other.Class().Name)
}

func (value Real) mod(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return mod(value, otherValue), nil
	}

	if otherValue, ok := other.(Int); ok {
		return mod(value, Real(otherValue)), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return mod(value, bo2ro(otherValue)), nil
	}

	return nil, NewErrorf("неможливо виконати модуль? дійсного числа  '%s'", other.Class().Name)
}

func (value Real) reversedMod(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return mod(otherValue, value), nil
	}

	if otherValue, ok := other.(Int); ok {
		return mod(Real(otherValue), value), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if value == 0.0 {
			return nil, NewZeroDivisionError("цілочисельне ділення або за модулем на нуль")
		}

		return mod(bo2ro(otherValue), value), nil
	}

	return nil, NewErrorf("неможливо виконати модуль? об'єкта '%s'  дійсне число", other.Class().Name)
}

func (value Real) pow(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Real); ok {
		return Real(math.Pow(float64(value), float64(otherValue))), nil
	}

	if otherValue, ok := other.(Int); ok {
		return Real(math.Pow(float64(value), float64(otherValue))), nil
	}

	if otherValue, ok := other.(Bool); ok {
		if otherValue {
			return value, nil
		}

		return Real(1.0), nil
	}

	return nil, NewErrorf("неможливо виконати степінь? дійсного числа  '%s'", other.Class().Name)
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

	if otherValue, ok := other.(Bool); ok {
		if otherValue {
			return Real(1.0), nil
		}

		if value < 0.0 {
			return nil, NewZeroDivisionError("неможливо піднести 0.0 до від'ємного степеня")
		}

		if value == 0.0 {
			return Real(1.0), nil
		}

		return Real(0.0), nil
	}

	return nil, NewErrorf("неможливо виконати степінь? об'єкта '%s'  дійсне число", other.Class().Name)
}

func (value Real) equals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value == v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value == Real(v)), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value == bo2ro(v)), nil
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

	if v, ok := other.(Bool); ok {
		return gb2bo(value != bo2ro(v)), nil
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

	if v, ok := other.(Bool); ok {
		return gb2bo(value < bo2ro(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew("<", value.Class().Name, other.Class().Name)
}

func (value Real) lessOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value <= v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value <= Real(v)), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value <= bo2ro(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew("<=", value.Class().Name, other.Class().Name)
}

func (value Real) greater(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value > v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value > Real(v)), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value > bo2ro(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew(">", value.Class().Name, other.Class().Name)
}

func (value Real) greaterOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value >= v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value >= Real(v)), nil
	}

	if v, ok := other.(Bool); ok {
		return gb2bo(value >= bo2ro(v)), nil
	}

	return nil, OperatorNotSupportedErrorNew(">=", value.Class().Name, other.Class().Name)
}

func (value Real) positive(_ Context) (Object, error) {
	return +value, nil
}

func (value Real) negate(_ Context) (Object, error) {
	return -value, nil
}
