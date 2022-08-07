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
	aLen := len(args)
	if aLen > 1 {
		return nil, ErrorNewf("дійсний() приймає 1 аргумент, або не приймає жодного")
	}

	if aLen > 0 {
		xObj = args[0]
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

func (value Real) equals(_ Context, other Object) (Object, error) {
	if v, ok := other.(Real); ok {
		return goBoolToBoolObject(value == v), nil
	}

	if v, ok := other.(Int); ok {
		return goBoolToBoolObject(value == Real(v)), nil
	}

	return False, nil
}

func (value Real) negate(_ Context) (Object, error) {
	return -value, nil
}
