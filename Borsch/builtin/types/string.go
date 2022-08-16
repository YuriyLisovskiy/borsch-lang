package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

var StringClass = ObjectClass.ClassNew("рядок", map[string]Object{}, true, StringNew, nil)

type String string

func (value String) Class() *Class {
	return StringClass
}

func (value String) represent(Context) (Object, error) {
	out, err := StringEscape(value, false)
	return String(out), err
}

func (value String) string(Context) (Object, error) {
	return value, nil
}

func (value String) toBool(Context) (Object, error) {
	return Bool(value != ""), nil
}

func (value String) add(_ Context, other Object) (Object, error) {
	if s, ok := other.(String); ok {
		return value + s, nil
	}

	return nil, NewErrorf("неможливо виконати конкатенацію рядка з об'єктом '%s'", other.Class().Name)
}

func (value String) reversedAdd(_ Context, other Object) (Object, error) {
	if s, ok := other.(String); ok {
		return s + value, nil
	}

	return nil, NewErrorf("неможливо виконати конкатенацію об'єкта '%s' з рядком", other.Class().Name)
}

func (value String) mul(_ Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		result := String("")
		if otherValue <= 0 {
			return result, nil
		}

		for i := int64(0); i < int64(otherValue); i++ {
			result += value
		}

		return result, nil
	}

	if otherValue, ok := other.(Bool); ok {
		if otherValue {
			return value, nil
		}

		return String(""), nil
	}

	return nil, NewErrorf("неможливо виконати множення рядка на об'єкт '%s'", other.Class().Name)
}

func (value String) reversedMul(ctx Context, other Object) (Object, error) {
	if otherValue, ok := other.(Int); ok {
		return otherValue.mul(ctx, value)
	}

	if otherValue, ok := other.(Bool); ok {
		if otherValue {
			return value, nil
		}

		return String(""), nil
	}

	return nil, NewErrorf("неможливо виконати множення об'єкта '%s' на рядок", other.Class().Name)
}

func (value String) equals(_ Context, other Object) (Object, error) {
	if s, ok := other.(String); ok {
		return goBoolToBoolObject(value == s), nil
	}

	if _, ok := other.(Bool); ok {
		return False, nil
	}

	return False, nil
}

func (value String) notEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(String); ok {
		return goBoolToBoolObject(value != v), nil
	}

	if _, ok := other.(Bool); ok {
		return True, nil
	}

	return False, nil
}

func (value String) less(_ Context, other Object) (Object, error) {
	if v, ok := other.(String); ok {
		return goBoolToBoolObject(value < v), nil
	}

	return False, nil
}

func (value String) lessOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(String); ok {
		return goBoolToBoolObject(value <= v), nil
	}

	return False, nil
}

func (value String) greater(_ Context, other Object) (Object, error) {
	if v, ok := other.(String); ok {
		return goBoolToBoolObject(value > v), nil
	}

	return False, nil
}

func (value String) greaterOrEquals(_ Context, other Object) (Object, error) {
	if v, ok := other.(String); ok {
		return goBoolToBoolObject(value >= v), nil
	}

	return False, nil
}

func StringEscape(a String, ascii bool) (string, error) {
	s := string(a)
	var out bytes.Buffer
	quote := '\''
	if strings.ContainsRune(s, '\'') && !strings.ContainsRune(s, '"') {
		quote = '"'
	}

	if !ascii {
		out.WriteRune(quote)
	}

	for _, c := range s {
		switch {
		case c < 0x20:
			switch c {
			case '\t':
				out.WriteString(`\t`)
			case '\n':
				out.WriteString(`\n`)
			case '\r':
				out.WriteString(`\r`)
			default:
				_, err := fmt.Fprintf(&out, `\x%02x`, c)
				if err != nil {
					// TODO: convert to ukr error!
					return "", err
				}
			}
		case !ascii && c < 0x7F:
			if c == '\\' || (quote == '\'' && c == '\'') || (quote == '"' && c == '"') {
				out.WriteRune('\\')
			}
			out.WriteRune(c)
		case c < 0x100:
			if ascii || strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				_, err := fmt.Fprintf(&out, "\\x%02x", c)
				if err != nil {
					// TODO: convert to ukr error!
					return "", err
				}
			}
		case c < 0x10000:
			if !ascii && strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				_, err := fmt.Fprintf(&out, "\\u%04x", c)
				if err != nil {
					// TODO: convert to ukr error!
					return "", err
				}
			}
		default:
			if !ascii && strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				_, err := fmt.Fprintf(&out, "\\U%08x", c)
				if err != nil {
					// TODO: convert to ukr error!
					return "", err
				}
			}
		}
	}

	if !ascii {
		out.WriteRune(quote)
	}

	return out.String(), nil
}

func StringNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	if len(args) != 1 {
		return nil, NewErrorf("рядок() приймає 1 аргумент")
	}

	return ToString(ctx, args[0])
}
