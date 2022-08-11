package types

var (
	BoolClass = ObjectClass.ClassNew("логічний", map[string]Object{}, true, BoolNew, nil)

	True  = Bool(true)
	False = Bool(false)
)

func goBoolToBoolObject(value bool) Object {
	if value {
		return True
	}

	return False
}

func gb2bo(value bool) Bool {
	if value {
		return True
	}

	return False
}

// bo2io converts Bool value to Int value.
func bo2io(value Bool) Int {
	if value {
		return 1
	}

	return 0
}

type Bool bool

func (value Bool) Class() *Class {
	return BoolClass
}

func NewBool(value bool) Bool {
	if value {
		return True
	}

	return False
}

func BoolNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	if len(args) != 1 {
		return nil, ErrorNewf("логічний() приймає 1 аргумент")
	}

	return ToBool(ctx, args[0])
}

func (value Bool) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value Bool) string(Context) (Object, error) {
	if value {
		return String("істина"), nil
	}

	return String("хиба"), nil
}

// func (value Bool) toBool(Context) (Object, error) {
// 	return value, nil
// }

func (value Bool) toReal(Context) (Object, error) {
	if value {
		return Real(1.0), nil
	}

	return Real(0.0), nil
}

func (value Bool) toInt(ctx Context) (Object, error) {
	return bo2io(value), nil
}

func (value Bool) add(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(value) + bo2io(otherValue), nil
	}

	if otherValue, ok := other.(Int); ok {
		return bo2io(value) + otherValue, nil
	}

	if otherValue, ok := other.(Real); ok {
		return bo2ro(value) + otherValue, nil
	}

	return nil, ErrorNewf("неможливо виконати додавання логічного значення до об'єкта '%s'", other.Class().Name)
}

func (value Bool) reversedAdd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) + bo2io(value), nil
	}

	if otherValue, ok := other.(Int); ok {
		return otherValue + bo2io(value), nil
	}

	if otherValue, ok := other.(Real); ok {
		return otherValue + bo2ro(value), nil
	}

	return nil, ErrorNewf("неможливо виконати додавання об'єкта '%s' до логічне значення", other.Class().Name)
}

func (value Bool) sub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(value) - bo2io(otherValue), nil
	}

	if otherValue, ok := other.(Int); ok {
		return bo2io(value) - otherValue, nil
	}

	if otherValue, ok := other.(Real); ok {
		return bo2ro(value) - otherValue, nil
	}

	return nil, ErrorNewf("неможливо виконати віднімання логічного значення від об'єкта '%s'", other.Class().Name)
}

func (value Bool) reversedSub(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) - bo2io(value), nil
	}

	if otherValue, ok := other.(Int); ok {
		return otherValue - bo2io(value), nil
	}

	if otherValue, ok := other.(Real); ok {
		return otherValue - bo2ro(value), nil
	}

	return nil, ErrorNewf("неможливо виконати віднімання об'єкта '%s' від логічне значення", other.Class().Name)
}

func (value Bool) div(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return bo2ro(value), nil
	}

	if otherValue, ok := other.(Int); ok {
		if otherValue == 0 {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return bo2ro(value) / Real(otherValue), nil
	}

	if otherValue, ok := other.(Real); ok {
		if otherValue == 0.0 {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return bo2ro(value) / otherValue, nil
	}

	return nil, ErrorNewf("неможливо виконати ділення логічного значення на об'єкт '%s'", other.Class().Name)
}

func (value Bool) reversedDiv(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		if !value {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return bo2ro(otherValue), nil
	}

	if otherValue, ok := other.(Int); ok {
		if !value {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return Real(otherValue), nil
	}

	if otherValue, ok := other.(Real); ok {
		if !value {
			return nil, ZeroDivisionErrorNewf("ділення на нуль")
		}

		return otherValue, nil
	}

	return nil, ErrorNewf("неможливо виконати ділення об'єкта '%s' на логічне значення", other.Class().Name)
}

func (value Bool) mul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(value && otherValue), nil
	}

	switch other.(type) {
	case Int, Real:
		if !value {
			return Int(0), nil
		}

		return other, nil
	}

	if otherValue, ok := other.(String); ok {
		if !value {
			return String(""), nil
		}

		return otherValue, nil
	}

	return nil, ErrorNewf("неможливо виконати множення логічного значення на об'єкт '%s'", other.Class().Name)
}

func (value Bool) reversedMul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue && value), nil
	}

	switch other.(type) {
	case Int, Real:
		if !value {
			return Int(0), nil
		}

		return other, nil
	}

	if otherValue, ok := other.(String); ok {
		if !value {
			return String(""), nil
		}

		return otherValue, nil
	}

	return nil, ErrorNewf("неможливо виконати множення об'єкта '%s' на логічне значення", other.Class().Name)
}

func (value Bool) mod(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		if !otherValue {
			return nil, ZeroDivisionErrorNewf("цілочисельне ділення або за модулем на нуль")
		}

		return Int(0), nil
	}

	if otherValue, ok := other.(Int); ok {
		if otherValue == 0 {
			return nil, ZeroDivisionErrorNewf("цілочисельне ділення або за модулем на нуль")
		}

		// return Int(mod(float64(bo2ro(value)), float64(otherValue))), nil
		return Int(mod(bo2ro(value), Real(otherValue))), nil
	}

	if otherValue, ok := other.(Real); ok {
		if otherValue == 0.0 {
			return nil, ZeroDivisionErrorNewf("цілочисельне ділення або за модулем на нуль")
		}

		return mod(bo2ro(value), otherValue), nil
	}

	return nil, ErrorNewf("неможливо виконати модуль? логічного значення  '%s'", other.Class().Name)
}

func (value Bool) reversedMod(_ Context, other Object) (Object, error) {
	switch other.(type) {
	case Bool, Int:
		if !value {
			return nil, ZeroDivisionErrorNewf("цілочисельне ділення або за модулем на нуль")
		}

		return Int(0), nil
	case Real:
		if !value {
			return nil, ZeroDivisionErrorNewf("цілочисельне ділення або за модулем на нуль")
		}

		return Real(0.0), nil
	}

	return nil, ErrorNewf("неможливо виконати модуль? об'єкта '%s'  логічне значення", other.Class().Name)
}

func (value Bool) pow(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(!(!value && otherValue)), nil
	}

	switch otherValue := other.(type) {
	case Int:
		if value {
			if otherValue < 0 {
				return Real(1.0), nil
			}

			return Int(1), nil
		}

		if otherValue < 0.0 {
			// TODO: error
		}

		if otherValue == 0 {
			return Int(1), nil
		}

		return Int(0), nil
	case Real:
		if value {
			return Real(1.0), nil
		}

		if otherValue < 0.0 {
			// TODO: error
		}

		if otherValue == 0.0 {
			return Real(1.0), nil
		}

		return Real(0.0), nil
	}

	return nil, ErrorNewf("неможливо виконати обчислення логічного значення в степені '%s'", other.Class().Name)
}

func (value Bool) reversedPow(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(!(!otherValue && value)), nil
	}

	if value {
		switch other.(type) {
		case Int, Real:
			return other, nil
		}
	}

	switch other.(type) {
	case Int:
		return Int(1), nil
	case Real:
		return Real(1.0), nil
	}

	return nil, ErrorNewf("неможливо виконати обчислення об'єкта '%s' в степені логічного значення", other.Class().Name)
}

func (value Bool) equals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Bool); ok {
		return gb2bo(value == v), nil
	}

	if v, ok := other.(Int); ok {
		return gb2bo(bo2io(value) == v), nil
	}

	if v, ok := other.(Real); ok {
		return gb2bo(bo2ro(value) == v), nil
	}

	return False, nil
}

func (value Bool) notEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Bool); ok {
		return gb2bo(value != v), nil
	}

	if v, ok := other.(Int); ok {
		return gb2bo(bo2io(value) != v), nil
	}

	if v, ok := other.(Real); ok {
		return gb2bo(bo2ro(value) != v), nil
	}

	if _, ok := other.(String); ok {
		return True, nil
	}

	return False, nil
}

func (value Bool) less(_ Context, other Object) (Object, error) {
	if v, ok := other.(Bool); ok {
		return gb2bo(!bool(value) && bool(v)), nil
	}

	if v, ok := other.(Int); ok {
		return gb2bo(bo2io(value) < v), nil
	}

	if v, ok := other.(Real); ok {
		return gb2bo(bo2ro(value) < v), nil
	}

	return nil, OperatorNotSupportedErrorNew("<", value.Class().Name, other.Class().Name)
}

func (value Bool) lessOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Bool); ok {
		return gb2bo(!(bool(value) && !bool(v))), nil
	}

	if v, ok := other.(Int); ok {
		return gb2bo(bo2io(value) <= v), nil
	}

	if v, ok := other.(Real); ok {
		return gb2bo(bo2ro(value) <= v), nil
	}

	return nil, OperatorNotSupportedErrorNew("<=", value.Class().Name, other.Class().Name)
}

func (value Bool) greater(_ Context, other Object) (Object, error) {
	if v, ok := other.(Bool); ok {
		return gb2bo(bool(value) && !bool(v)), nil
	}

	if v, ok := other.(Int); ok {
		return gb2bo(bo2io(value) > v), nil
	}

	if v, ok := other.(Real); ok {
		return gb2bo(bo2ro(value) > v), nil
	}

	return nil, OperatorNotSupportedErrorNew(">", value.Class().Name, other.Class().Name)
}

func (value Bool) greaterOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Bool); ok {
		return gb2bo(bool(value) || !bool(v)), nil
	}

	if v, ok := other.(Int); ok {
		return gb2bo(bo2io(value) >= v), nil
	}

	if v, ok := other.(Real); ok {
		return gb2bo(bo2ro(value) >= v), nil
	}

	return nil, OperatorNotSupportedErrorNew(">=", value.Class().Name, other.Class().Name)
}

func (value Bool) shiftLeft(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(value) << bo2io(otherValue), nil
	}

	if otherValue, ok := other.(Int); ok {
		return bo2io(value) << otherValue, nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітовий зсув ліворуч логічного значення на об'єкт '%s'",
		other.Class().Name,
	)
}

func (value Bool) reversedShiftLeft(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) << bo2io(value), nil
	}

	if otherValue, ok := other.(Int); ok {
		return otherValue << bo2io(value), nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітовий зсув ліворуч об'єкта '%s' на логічне значення",
		other.Class().Name,
	)
}

func (value Bool) shiftRight(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(value) >> bo2io(otherValue), nil
	}

	if otherValue, ok := other.(Int); ok {
		return bo2io(value) >> otherValue, nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітовий зсув праворуч логічного значення на об'єкт '%s'",
		other.Class().Name,
	)
}

func (value Bool) reversedShiftRight(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return bo2io(otherValue) >> bo2io(value), nil
	}

	if otherValue, ok := other.(Int); ok {
		return otherValue >> bo2io(value), nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітовий зсув праворуч об'єкта '%s' на логічного значення",
		other.Class().Name,
	)
}

func (value Bool) bitwiseOr(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return io2bo(bo2io(value) | bo2io(otherValue)), nil
	}
	if otherValue, ok := other.(Int); ok {
		return bo2io(value) | otherValue, nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітову диз'юнкцію логічного значення та об'єкта '%s'",
		other.Class().Name,
	)
}

func (value Bool) reversedBitwiseOr(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return io2bo(bo2io(otherValue) | bo2io(value)), nil
	}

	if otherValue, ok := other.(Int); ok {
		return otherValue | bo2io(value), nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітову диз'юнкцію об'єкта '%s' та логічного значення",
		other.Class().Name,
	)
}

func (value Bool) bitwiseXor(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return io2bo(bo2io(value) ^ bo2io(otherValue)), nil
	}

	if otherValue, ok := other.(Int); ok {
		return bo2io(value) ^ otherValue, nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітову виняткову диз'юнкцію логічного значення та об'єкта '%s'",
		other.Class().Name,
	)
}

func (value Bool) reversedBitwiseXor(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return io2bo(bo2io(otherValue) ^ bo2io(value)), nil
	}

	if otherValue, ok := other.(Int); ok {
		return otherValue ^ bo2io(value), nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітову виняткову диз'юнкцію об'єкта '%s' та логічного значення",
		other.Class().Name,
	)
}

func (value Bool) bitwiseAnd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return io2bo(bo2io(value) & bo2io(otherValue)), nil
	}

	if otherValue, ok := other.(Int); ok {
		return bo2io(value) & otherValue, nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітову кон'юнкцію логічного значення та об'єкта '%s'",
		other.Class().Name,
	)
}

func (value Bool) reversedBitwiseAnd(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Bool); ok {
		return io2bo(bo2io(otherValue) & bo2io(value)), nil
	}

	if otherValue, ok := other.(Int); ok {
		return otherValue & bo2io(value), nil
	}

	return nil, ErrorNewf(
		"неможливо виконати побітову кон'юнкцію об'єкта '%s' та логічного значення",
		other.Class().Name,
	)
}

func (value Bool) positive(_ Context) (Object, error) {
	return +bo2io(value), nil
}

func (value Bool) negate(_ Context) (Object, error) {
	return -bo2io(value), nil
}

func (value Bool) invert(_ Context) (Object, error) {
	return ^bo2io(value), nil
}
