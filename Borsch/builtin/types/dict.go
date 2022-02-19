// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.
//
// Dict and Dict type.
//
// The idea is that most dictionaries just have strings for keys,
// so we use the simpler Dict and promote it into a Dict
// when necessary.

package types

import "bytes"

// FIXME: translate the end of the doc
const dictDoc = `словник() -> новий порожній словник
словник(пари) -> новий словник із пар об'єктів.
    пари (ключ, значення)
словник(ітерований_об_єкт) -> новий словник, ініціалізовний подібно до:
    д = {};
    цикл к, з : ітерований_об_єкт:
        д[к] = з;
	кінець`

var (
	StringDictClass = NewClass("словник", dictDoc)
	DictClass       = NewClass("словник", dictDoc)
	expectingDict   = ErrorNewf(TypeError, "необхідний словник")
)

// func init() {
// 	StringDictType.Dict["items"] = MustNewMethod(
// 		"items", func(self Object, args Tuple) (Object, error) {
// 			err := UnpackTuple(args, nil, "items", 0, 0)
// 			if err != nil {
// 				return nil, err
// 			}
// 			sMap := self.(Dict)
// 			o := make([]Object, 0, len(sMap))
// 			for k, v := range sMap {
// 				o = append(o, Tuple{String(k), v})
// 			}
// 			return NewIterator(o), nil
// 		}, 0, "items() -> list of D's (key, value) pairs, as 2-tuples",
// 	)
//
// 	StringDictType.Dict["get"] = MustNewMethod(
// 		"get", func(self Object, args Tuple) (Object, error) {
// 			var length = len(args)
// 			switch {
// 			case length == 0:
// 				return nil, ExceptionNewf(TypeError, "%s expected at least 1 arguments, got %d", "items()", length)
// 			case length > 2:
// 				return nil, ExceptionNewf(TypeError, "%s expected at most 2 arguments, got %d", "items()", length)
// 			}
// 			sMap := self.(Dict)
// 			if str, ok := args[0].(String); ok {
// 				if res, ok := sMap[string(str)]; ok {
// 					return res, nil
// 				}
//
// 				switch length {
// 				case 2:
// 					return args[1], nil
// 				default:
// 					return None, nil
// 				}
// 			}
// 			return nil, ExceptionNewf(KeyError, "%v", args[0])
// 		}, 0, "gets(key, default) -> If there is a val corresponding to key, return val, otherwise default",
// 	)
// }

// Dict is String to Object dictionary.
//
// Used for variables, etc. where the keys can only be strings.
// FIXME: change type to map[string]Object
type Dict map[string]Object

func (d Dict) Class() *Class {
	return StringDictClass
}

// NewStringDict makes a new dictionary.
func NewStringDict() Dict {
	return make(Dict)
}

// NewStringDictSized makes a new dictionary with reservation for n entries.
func NewStringDictSized(n int) Dict {
	return make(Dict, n)
}

// DictCheckExact checks that obj is exactly a dictionary and returns an error if not.
func DictCheckExact(obj Object) (Dict, error) {
	dict, ok := obj.(Dict)
	if !ok {
		return nil, expectingDict
	}

	return dict, nil
}

// DictCheck checks that obj is exactly a dictionary and returns an error if not.
func DictCheck(obj Object) (Dict, error) {
	// FIXME should be checking subclasses
	return DictCheckExact(obj)
}

func (d Dict) Copy() Dict {
	e := make(Dict, len(d))
	for k, v := range d {
		e[k] = v
	}

	return e
}

func (d Dict) __str__() (Object, error) {
	return d.__represent__()
}

func (d Dict) __represent__() (Object, error) {
	var out bytes.Buffer
	out.WriteRune('{')
	spacer := false
	for key, value := range d {
		if spacer {
			out.WriteString(", ")
		}

		keyStr, err := RepresentAsString(String(key))
		if err != nil {
			return nil, err
		}

		valueStr, err := RepresentAsString(value)
		if err != nil {
			return nil, err
		}

		out.WriteString(keyStr)
		out.WriteString(": ")
		out.WriteString(valueStr)
		spacer = true
	}

	out.WriteRune('}')
	return String(out.String()), nil
}

// Returns a list of keys from the dict
func (d Dict) __iter__() (Object, error) {
	panic("unreachable")
	// TODO:
	// o := make([]Object, 0, len(d))
	// for k := range d {
	// 	o = append(o, String(k))
	// }
	//
	// return NewIterator(o), nil
}

func (d Dict) __get_item__(key Object) (Object, error) {
	str, ok := key.(String)
	if ok {
		res, ok := d[string(str)]
		if ok {
			return res, nil
		}
	}

	return nil, ErrorNewf(KeyError, "%v", key)
}

func (d Dict) __set_item__(key, value Object) (Object, error) {
	str, ok := key.(String)
	if !ok {
		return nil, ErrorNewf(KeyError, "FIXME: can only have string keys!: %v", key)
	}

	d[string(str)] = value
	return Nil, nil
}

func (d Dict) __equal__(other Object) (Object, error) {
	b, ok := other.(Dict)
	if !ok {
		return NotImplemented, nil
	}

	if len(d) != len(b) {
		return False, nil
	}

	for k, av := range d {
		bv, ok := b[k]
		if !ok {
			return False, nil
		}

		res, err := Equal(av, bv)
		if err != nil {
			return nil, err
		}

		if res == False {
			return False, nil
		}
	}

	return True, nil
}

func (d Dict) __not_equal__(other Object) (Object, error) {
	res, err := d.__equal__(other)
	if err != nil {
		return nil, err
	}

	if res == NotImplemented {
		return res, nil
	}

	if res == True {
		return False, nil
	}

	return True, nil
}

func (d Dict) __contains__(other Object) (Object, error) {
	key, ok := other.(String)
	if !ok {
		return nil, ErrorNewf(KeyError, "FIXME can only have string keys!: %v", key)
	}

	if _, ok := d[string(key)]; ok {
		return True, nil
	}

	return False, nil
}

func (d Dict) GetDict() Dict {
	return d
}

var _ IGetDict = (*Dict)(nil)
