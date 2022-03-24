package types

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var IntClass = ObjectClass.ClassNew("цілий", map[string]Object{}, true, IntNew, nil)

type Int int64

func (value Int) Class() *Class {
	return IntClass
}

func IntNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	var xObj Object = Int(0)
	var baseObj Object
	base := 0
	aLen := len(args)
	if aLen > 2 {
		return nil, ErrorNewf("цілий() приймає 1 або 2 аргументи, або не приймає жодного")
	}

	if aLen > 0 {
		xObj = args[0]
		if aLen == 2 {
			baseObj = args[1]
		}
	}

	var err error
	if baseObj != nil {
		base, err = ToGoInt(ctx, baseObj)
		if err != nil {
			return nil, err
		}

		if base != 0 && (base < 2 || base > 36) {
			return nil, ErrorNewf("база для 'цілий()' має бути >= 2 та <= 36")
		}
	}

	if x, ok := xObj.(String); ok {
		return IntFromString(string(x), base)
	}

	if baseObj != nil {
		// TODO: TypeError
		return nil, ErrorNewf("'цілий()' не може перетворити об'єкт не рядкового типу з явною базою")
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
	return nil, ErrorNewf("overflow...")

error:
	// TODO: ValueError
	return nil, ErrorNewf("некоректний літерал для 'цілий()' з базою %d: '%s'", convertBase, str)
}

func (value Int) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value Int) string(Context) (Object, error) {
	return String(fmt.Sprintf("%d", value)), nil
}

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
