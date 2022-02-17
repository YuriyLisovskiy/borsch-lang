// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.
//
// Dict and StringDict type.
//
// The idea is that most dictionaries just have strings for keys,
// so we use the simpler StringDict and promote it into a Dict
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
	кінець
dict(**kwargs) -> new dictionary initialized with the name=value pairs
    in the keyword argument list.  For example:  dict(one=1, two=2)`

var (
	StringDictType = NewType("словник", dictDoc)
	DictType       = NewType("словник", dictDoc)
	expectingDict  = ErrorNewf(TypeError, "необхідний словник")
)

// func init() {
// 	StringDictType.Dict["items"] = MustNewMethod(
// 		"items", func(self Object, args Tuple) (Object, error) {
// 			err := UnpackTuple(args, nil, "items", 0, 0)
// 			if err != nil {
// 				return nil, err
// 			}
// 			sMap := self.(StringDict)
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
// 			sMap := self.(StringDict)
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

// StringDict is String to Object dictionary.
//
// Used for variables, etc. where the keys can only be strings.
type StringDict map[string]Object

// Type of this StringDict object
func (d StringDict) Type() *Type {
	return StringDictType
}

// NewStringDict makes a new dictionary.
func NewStringDict() StringDict {
	return make(StringDict)
}

// NewStringDictSized makes a new dictionary with reservation for n entries.
func NewStringDictSized(n int) StringDict {
	return make(StringDict, n)
}

// DictCheckExact checks that obj is exactly a dictionary and returns an error if not.
func DictCheckExact(obj Object) (StringDict, error) {
	dict, ok := obj.(StringDict)
	if !ok {
		return nil, expectingDict
	}

	return dict, nil
}

// DictCheck checks that obj is exactly a dictionary and returns an error if not.
func DictCheck(obj Object) (StringDict, error) {
	// FIXME should be checking subclasses
	return DictCheckExact(obj)
}

func (d StringDict) Copy() StringDict {
	e := make(StringDict, len(d))
	for k, v := range d {
		e[k] = v
	}

	return e
}

func (d StringDict) __str__() (Object, error) {
	return d.__represent__()
}

func (d StringDict) __represent__() (Object, error) {
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
func (d StringDict) __iter__() (Object, error) {
	o := make([]Object, 0, len(d))
	for k := range d {
		o = append(o, String(k))
	}

	return NewIterator(o), nil
}

func (d StringDict) __get_item__(key Object) (Object, error) {
	str, ok := key.(String)
	if ok {
		res, ok := d[string(str)]
		if ok {
			return res, nil
		}
	}

	return nil, ErrorNewf(KeyError, "%v", key)
}

func (d StringDict) __set_item__(key, value Object) (Object, error) {
	str, ok := key.(String)
	if !ok {
		return nil, ErrorNewf(KeyError, "FIXME: can only have string keys!: %v", key)
	}

	d[string(str)] = value
	return Nil, nil
}

func (d StringDict) __equal__(other Object) (Object, error) {
	b, ok := other.(StringDict)
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

func (d StringDict) __not_equal__(other Object) (Object, error) {
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

func (d StringDict) __contains__(other Object) (Object, error) {
	key, ok := other.(String)
	if !ok {
		return nil, ErrorNewf(KeyError, "FIXME can only have string keys!: %v", key)
	}

	if _, ok := d[string(key)]; ok {
		return True, nil
	}

	return False, nil
}

func (d StringDict) GetDict() StringDict {
	return d
}

var _ IGetDict = (*StringDict)(nil)
