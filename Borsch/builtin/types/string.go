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

func (value String) add(_ Context, other Object) (Object, error) {
	if s, ok := other.(String); ok {
		return value + s, nil
	}

	return nil, ErrorNewf("неможливо виконати конкатенацію рядка з об'єктом '%s'", other.Class().Name)
}

func (value String) reversedAdd(_ Context, other Object) (Object, error) {
	if s, ok := other.(String); ok {
		return s + value, nil
	}

	return nil, ErrorNewf("неможливо виконати конкатенацію об'єкта '%s' з рядком", other.Class().Name)
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
		return nil, ErrorNewf("рядок() приймає 1 аргумент")
	}

	return ToString(ctx, args[0])
}
