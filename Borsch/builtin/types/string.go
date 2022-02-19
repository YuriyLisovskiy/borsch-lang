// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// String objects
//
// Note that string objects in Python are arrays of unicode
// characters. However we are using the native Go string which is
// UTF-8 encoded.  This makes very little difference most of the time,
// but care is needed when indexing, slicing or iterating through
// strings.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type String string

var StringType = ObjectClass.NewClass(
	"рядок",

	// FIXME: translate the end of the doc.
	`рядок(об_єкт="") -> рядок
рядок(байти_або_буфер[, кодування[, помилки]]) -> рядок

Створює новий об'єкт рядкового типу із переданого об'єкта.
Якщо вказані кодування або помилки, тоді об'єкт повинен надати буфер даних
які будуть декодовані за допомогою заданого кодування та обробника помилок.
Інакше, повертає результат виклику об_єкт.__рядок__() (якщо він визначений)
або представлення(об_єкт).
encoding defaults to sys.getdefaultencoding().
errors defaults to 'strict'.`, StrNew, nil,
)

// StringEscape escapes the String.
func StringEscape(a String, ascii bool) string {
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
				fmt.Fprintf(&out, `\x%02x`, c)
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
				fmt.Fprintf(&out, "\\x%02x", c)
			}
		case c < 0x10000:
			if !ascii && strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				fmt.Fprintf(&out, "\\u%04x", c)
			}
		default:
			if !ascii && strconv.IsPrint(c) {
				out.WriteRune(c)
			} else {
				fmt.Fprintf(&out, "\\U%08x", c)
			}
		}
	}

	if !ascii {
		out.WriteRune(quote)
	}

	return out.String()
}

// Standard golang strings.Fields doesn't have a 'first N' argument
func fieldsN(s string, n int) []string {
	var out []string
	var cur []rune
	for _, c := range s {
		// until we have covered the first N elements, multiple white-spaces are 'merged'
		if n < 0 || len(out) < n {
			if unicode.IsSpace(c) {
				if len(cur) > 0 {
					out = append(out, string(cur))
					cur = []rune{}
				}
			} else {
				cur = append(cur, c)
			}
			// until we see the next letter, after collecting the first N fields, continue to merge whitespaces
		} else if len(out) == n && len(cur) == 0 {
			if !unicode.IsSpace(c) {
				cur = append(cur, c)
			}
			// now that enough words have been collected, just copy into the last element
		} else {
			cur = append(cur, c)
		}
	}

	if len(cur) > 0 {
		out = append(out, string(cur))
	}

	return out
}

// func init() {
// 	StringType.Dict["split"] = MustNewMethod(
// 		"split", func(self Object, args Tuple) (Object, error) {
// 			selfStr := self.(String)
// 			var value Object = Nil
// 			zeroRemove := true
// 			if len(args) > 0 {
// 				if _, ok := args[0].(NilType); !ok {
// 					value = args[0]
// 					zeroRemove = false
// 				}
// 			}
//
// 			var maxSplit = -2
// 			if len(args) > 1 {
// 				if m, ok := args[1].(Int); ok {
// 					maxSplit = int(m)
// 				}
// 			}
//
// 			var valArray []string
// 			if valStr, ok := value.(String); ok {
// 				valArray = strings.SplitN(string(selfStr), string(valStr), maxSplit+1)
// 			} else if _, ok := value.(NilType); ok {
// 				valArray = fieldsN(string(selfStr), maxSplit)
// 			} else {
// 				return nil, ErrorNewf(TypeError, "Неможливо явно перетворити об'єкт '%s' у рядок", value.Type())
// 			}
//
// 			o := List{}
// 			for _, j := range valArray {
// 				if len(j) > 0 || !zeroRemove {
// 					o.Items = append(o.Items, String(j))
// 				}
// 			}
//
// 			return &o, nil
// 		}, 0, "split(sub) -> split string with sub.",
// 	)
//
// 	StringType.Dict["startswith"] = MustNewMethod(
// 		"startswith", func(self Object, args Tuple) (Object, error) {
// 			selfStr := string(self.(String))
// 			var prefix []string
// 			if len(args) > 0 {
// 				if s, ok := args[0].(String); ok {
// 					prefix = append(prefix, string(s))
// 				} else if s, ok := args[0].(Tuple); ok {
// 					for _, t := range s {
// 						if v, ok := t.(String); ok {
// 							prefix = append(prefix, string(v))
// 						}
// 					}
// 				} else {
// 					return nil, ErrorNewf(
// 						TypeError,
// 						"startswith first arg must be str, unicode, or tuple, not %s",
// 						args[0].Type(),
// 					)
// 				}
// 			} else {
// 				return nil, ErrorNewf(TypeError, "startswith() takes at least 1 argument (0 given)")
// 			}
//
// 			if len(args) > 1 {
// 				if s, ok := args[1].(Int); ok {
// 					selfStr = selfStr[s:]
// 				}
// 			}
//
// 			for _, s := range prefix {
// 				if strings.HasPrefix(selfStr, s) {
// 					return Bool(true), nil
// 				}
// 			}
//
// 			return Bool(false), nil
// 		}, 0, "startswith(prefix[, start[, end]]) -> bool",
// 	)
//
// 	StringType.Dict["endswith"] = MustNewMethod(
// 		"endswith", func(self Object, args Tuple) (Object, error) {
// 			selfStr := string(self.(String))
// 			var suffix []string
// 			if len(args) > 0 {
// 				if s, ok := args[0].(String); ok {
// 					suffix = append(suffix, string(s))
// 				} else if s, ok := args[0].(Tuple); ok {
// 					for _, t := range s {
// 						if v, ok := t.(String); ok {
// 							suffix = append(suffix, string(v))
// 						}
// 					}
// 				} else {
// 					return nil, ErrorNewf(
// 						TypeError,
// 						"endswith first arg must be str, unicode, or tuple, not %s",
// 						args[0].Type(),
// 					)
// 				}
// 			} else {
// 				return nil, ErrorNewf(TypeError, "endswith() takes at least 1 argument (0 given)")
// 			}
//
// 			for _, s := range suffix {
// 				if strings.HasSuffix(selfStr, s) {
// 					return Bool(true), nil
// 				}
// 			}
//
// 			return Bool(false), nil
// 		}, 0, "endswith(suffix[, start[, end]]) -> bool",
// 	)
// }

func (value String) Class() *Class {
	return StringType
}

func StrNew(cls *Class, args Tuple) (Object, error) {
	var sObj Object
	if err := ParseExactArgs(args, "рядок|a", &sObj); err != nil {
		return nil, err
	}

	return Str(sObj)
}

func (value String) __str__() (Object, error) {
	return value, nil
}

func (value String) __represent__() (Object, error) {
	out := StringEscape(value, false)
	return String(out), nil
}

func (value String) __bool__() (Object, error) {
	return NewBool(len(value) > 0), nil
}

// len returns length of the string in unicode characters
func (value String) length() int {
	return utf8.RuneCountInString(string(value))
}

func (value String) __length__() (Object, error) {
	return Int(value.length()), nil
}

func (value String) __add__(other Object) (Object, error) {
	if b, ok := other.(String); ok {
		return value + b, nil
	}

	return NotImplemented, nil
}

func (value String) __reversed_add__(other Object) (Object, error) {
	if b, ok := other.(String); ok {
		return b + value, nil
	}

	return NotImplemented, nil
}

func (value String) __in_place_add__(other Object) (Object, error) {
	return value.__add__(other)
}

func (value String) __mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		if b < 0 {
			b = 0
		}

		var out bytes.Buffer
		for i := 0; i < int(b); i++ {
			out.WriteString(string(value))
		}

		return String(out.String()), nil
	}

	return NotImplemented, nil
}

func (value String) __reversed_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value String) __in_place_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

// Convert an Object to a String.
//
// Returns ok if the conversion worked or not.
func convertToString(other Object) (String, bool) {
	switch b := other.(type) {
	case String:
		return b, true
	}

	return "", false
}

// Comparison

func (value String) __less_than__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(value < b), nil
	}

	return NotImplemented, nil
}

func (value String) __less_or_equal__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(value <= b), nil
	}

	return NotImplemented, nil
}

func (value String) __equal__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(value == b), nil
	}

	return NotImplemented, nil
}

func (value String) __not_equal__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(value != b), nil
	}

	return NotImplemented, nil
}

func (value String) __greater_than__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(value > b), nil
	}

	return NotImplemented, nil
}

func (value String) __greater_or_equal__(other Object) (Object, error) {
	if b, ok := convertToString(other); ok {
		return NewBool(value >= b), nil
	}

	return NotImplemented, nil
}

// Returns position in string of n-th character
//
// returns end of string if not found
func (value String) pos(n int) int {
	characterNumber := 0
	for i := range value {
		if characterNumber == n {
			return i
		}

		characterNumber++
	}

	return len(value)
}

// slice returns the slice of this string using character positions
//
// length should be the length of the string in unicode characters
func (value String) slice(start, stop, length int) String {
	if start >= stop {
		return ""
	}

	if length == len(value) {
		return value[start:stop] // ascii only
	}

	if start <= 0 && stop >= length {
		return value
	}

	startI := value.pos(start)
	stopI := value[startI:].pos(stop-start) + startI
	return value[startI:stopI]
}

func (value String) __get_item__(key Object) (Object, error) {
	length := value.length()
	asciiOnly := length == len(value)
	i, err := IndexIntCheck(key, length)
	if err != nil {
		return nil, err
	}

	if asciiOnly {
		return value[i : i+1], nil
	}

	newValue := value[value.pos(i):]
	_, runeSize := utf8.DecodeRuneInString(string(newValue))
	return newValue[:runeSize], nil
}

func (value String) __contains__(item Object) (Object, error) {
	needle, ok := item.(String)
	if !ok {
		return nil, ErrorNewf(TypeError, "'in <string>' requires string as left operand, not %s", item.Class().Name)
	}

	return NewBool(strings.Contains(string(value), string(needle))), nil
}

// Check interface is satisfied
var _ iComparison = String("")
var _ sequenceArithmetic = String("")
var _ I__length__ = String("")
var _ I__bool__ = String("")
var _ I__get_item__ = String("")
var _ I__contains__ = String("")
