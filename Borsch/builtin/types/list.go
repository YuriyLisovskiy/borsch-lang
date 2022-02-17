// Copyright 2018 The go-python Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2022 The Borsch Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// List objects.

package types

import (
	"sort"
)

var ListType = ObjectType.NewType(
	"список",
	`список() -> новий порожній список
список(ітерований_об_єкт) -> новий список, ініціалізований з елементів ітерованого об'єкта'`,
	ListNew,
	nil,
)

type List struct {
	Items []Object
}

func init() {
	// ListType.Dict["append"] = MustNewMethod(
	// 	"append", func(self Object, args Tuple) (Object, error) {
	// 		listSelf := self.(*List)
	// 		if len(args) != 1 {
	// 			return nil, ErrorNewf(TypeError, "append() takes exactly one argument (%d given)", len(args))
	// 		}
	//
	// 		listSelf.Items = append(listSelf.Items, args[0])
	// 		return NilType{}, nil
	// 	}, 0, "append(item)",
	// )
	// ListType.Dict["sort"] = MustNewMethod(
	// 	"sort", func(self Object, args Tuple, kwargs StringDict) (Object, error) {
	// 		const funcName = "sort"
	// 		l, isList := self.(*List)
	// 		if !isList {
	// 			// method called using `list.sort([], **kwargs)`
	// 			var o Object
	// 			err := UnpackTuple(args, nil, funcName, 1, 1, &o)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	//
	// 			var ok bool
	// 			l, ok = o.(*List)
	// 			if !ok {
	// 				return nil, ErrorNewf(
	// 					TypeError,
	// 					"descriptor 'sort' requires a 'list' object but received a '%s'",
	// 					o.Type(),
	// 				)
	// 			}
	// 		} else {
	// 			// method called using `[].sort(**kargs)`
	// 			err := UnpackTuple(args, nil, funcName, 0, 0)
	// 			if err != nil {
	// 				return nil, err
	// 			}
	// 		}
	//
	// 		err := SortInPlace(l, kwargs, funcName)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	//
	// 		return NilType{}, nil
	// 	}, 0, "sort(key=None, reverse=False)",
	// )
}

// Type of this List object.
func (value *List) Type() *Type {
	return ListType
}

func ListNew(cls *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	var iterable Object
	err = UnpackTuple(args, kwargs, "список", 0, 1, &iterable)
	if err != nil {
		return nil, err
	}

	if iterable != nil {
		return SequenceList(iterable)
	}

	return NewList(), nil
}

// NewList makes a new empty list.
func NewList() *List {
	return &List{}
}

// NewListWithCapacity makes a new empty list with given capacity.
func NewListWithCapacity(n int) *List {
	l := &List{}
	if n != 0 {
		l.Items = make([]Object, 0, n)
	}

	return l
}

// NewListSized makes a list with n nil elements.
func NewListSized(n int) *List {
	l := &List{}
	if n != 0 {
		l.Items = make([]Object, n)
	}

	return l
}

// NewListFromItems makes a new list from an []Object.
//
// The []Object is copied into the list.
func NewListFromItems(items []Object) *List {
	l := NewListSized(len(items))
	copy(l.Items, items)
	return l
}

// NewListFromStrings makes an argv into a tuple.
func NewListFromStrings(items []string) *List {
	l := NewListSized(len(items))
	for i, v := range items {
		l.Items[i] = String(v)
	}

	return l
}

// Copy a list object.
func (value *List) Copy() *List {
	return NewListFromItems(value.Items)
}

// Append an item.
func (value *List) Append(item Object) {
	value.Items = append(value.Items, item)
}

// Resize the list.
func (value *List) Resize(newSize int) {
	value.Items = value.Items[:newSize]
}

// Extend the list with items
func (value *List) Extend(items []Object) {
	value.Items = append(value.Items, items...)
}

// ExtendWithStrings extends the list with strings.
func (value *List) ExtendWithStrings(items []string) {
	for _, item := range items {
		value.Items = append(value.Items, Object(String(item)))
	}
}

// ExtendSequence extends the list with the sequence passed in.
func (value *List) ExtendSequence(seq Object) error {
	return Iterate(
		seq, func(item Object) bool {
			value.Append(item)
			return false
		},
	)
}

// Length of list.
func (value *List) Length() int {
	return len(value.Items)
}

func (value *List) __str__() (Object, error) {
	return value.__represent__()
}

func (value *List) __represent__() (Object, error) {
	return Tuple(value.Items).represent("[", "]")
}

func (value *List) __length__() (Object, error) {
	return Int(len(value.Items)), nil
}

func (value *List) __bool__() (Object, error) {
	return NewBool(len(value.Items) > 0), nil
}

func (value *List) __iter__() (Object, error) {
	return NewIterator(value.Items), nil
}

func (value *List) __get_item__(key Object) (Object, error) {
	i, err := IndexIntCheck(key, len(value.Items))
	if err != nil {
		return nil, err
	}

	return value.Items[i], nil
}

func (value *List) __set_item__(key, item Object) (Object, error) {
	i, err := IndexIntCheck(key, len(value.Items))
	if err != nil {
		return nil, err
	}

	value.Items[i] = item

	return Nil, nil
}

// DeleteItem removes the item at i.
func (value *List) DeleteItem(i int) {
	value.Items = append(value.Items[:i], value.Items[i+1:]...)
}

// __delete_item__ removes items from a list.
func (value *List) __delete_item__(key Object) (Object, error) {
	i, err := IndexIntCheck(key, len(value.Items))
	if err != nil {
		return nil, err
	}

	value.DeleteItem(i)
	return Nil, nil
}

func (value *List) __add__(other Object) (Object, error) {
	if b, ok := other.(*List); ok {
		newList := NewListSized(len(value.Items) + len(b.Items))
		copy(newList.Items, value.Items)
		copy(newList.Items[len(value.Items):], b.Items)
		return newList, nil
	}

	return NotImplemented, nil
}

func (value *List) __reversed_add__(other Object) (Object, error) {
	if b, ok := other.(*List); ok {
		return b.__add__(value)
	}

	return NotImplemented, nil
}

func (value *List) __in_place_add__(other Object) (Object, error) {
	if b, ok := other.(*List); ok {
		value.Extend(b.Items)
		return value, nil
	}

	return NotImplemented, nil
}

func (value *List) __mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		m := len(value.Items)
		n := int(b) * m
		if n < 0 {
			n = 0
		}

		newList := NewListSized(n)
		for i := 0; i < n; i += m {
			copy(newList.Items[i:i+m], value.Items)
		}

		return newList, nil
	}

	return NotImplemented, nil
}

func (value *List) __reversed_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

func (value *List) __in_place_mul__(other Object) (Object, error) {
	return value.__mul__(other)
}

// Check interface is satisfied
var _ sequenceArithmetic = (*List)(nil)
var _ I__str__ = (*List)(nil)
var _ I__represent__ = (*List)(nil)
var _ I__length__ = (*List)(nil)
var _ I__bool__ = (*List)(nil)
var _ I__iter__ = (*List)(nil)
var _ I__get_item__ = (*List)(nil)
var _ I__set_item__ = (*List)(nil)

func (value *List) __equal__(other Object) (Object, error) {
	b, ok := other.(*List)
	if !ok {
		return NotImplemented, nil
	}

	if len(value.Items) != len(b.Items) {
		return False, nil
	}

	for i := range value.Items {
		eq, err := Equal(value.Items[i], b.Items[i])
		if err != nil {
			return nil, err
		}

		if eq == False {
			return False, nil
		}
	}

	return True, nil
}

func (value *List) __not_equal__(other Object) (Object, error) {
	b, ok := other.(*List)
	if !ok {
		return NotImplemented, nil
	}

	if len(value.Items) != len(b.Items) {
		return True, nil
	}

	for i := range value.Items {
		eq, err := Equal(value.Items[i], b.Items[i])
		if err != nil {
			return nil, err
		}

		if eq == False {
			return True, nil
		}
	}

	return False, nil
}

type sortable struct {
	l        *List
	keyFunc  Object
	reverse  bool
	firstErr error
}

type ptrSortable struct {
	s *sortable
}

func (s ptrSortable) Len() int {
	return s.s.l.Length()
}

func (s ptrSortable) Swap(i, j int) {
	itemI, err := s.s.l.__get_item__(Int(i))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return
	}
	itemJ, err := s.s.l.__get_item__(Int(j))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return
	}
	_, err = s.s.l.__set_item__(Int(i), itemJ)
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
	}
	_, err = s.s.l.__set_item__(Int(j), itemI)
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
	}
}

func (s ptrSortable) Less(i, j int) bool {
	itemI, err := s.s.l.__get_item__(Int(i))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return false
	}
	itemJ, err := s.s.l.__get_item__(Int(j))
	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return false
	}

	if s.s.keyFunc != Nil {
		itemI, err = Call(s.s.keyFunc, Tuple{itemI}, nil)
		if err != nil {
			if s.s.firstErr == nil {
				s.s.firstErr = err
			}
			return false
		}
		itemJ, err = Call(s.s.keyFunc, Tuple{itemJ}, nil)
		if err != nil {
			if s.s.firstErr == nil {
				s.s.firstErr = err
			}
			return false
		}
	}

	var cmpResult Object
	if s.s.reverse {
		cmpResult, err = LessThan(itemJ, itemI)
	} else {
		cmpResult, err = LessThan(itemI, itemJ)
	}

	if err != nil {
		if s.s.firstErr == nil {
			s.s.firstErr = err
		}
		return false
	}

	if boolResult, ok := cmpResult.(Bool); ok {
		return bool(boolResult)
	}

	return false
}

// SortInPlace sorts the given List in place using a stable sort.
// kwargs can have the keys "ключ" and "зворотний".
func SortInPlace(l *List, kwargs StringDict, funcName string) error {
	var keyFunc Object
	var reverse Object
	err := ParseTupleAndKeywords(nil, kwargs, "|$OO:"+funcName, []string{"ключ", "зворотний"}, &keyFunc, &reverse)
	if err != nil {
		return err
	}

	if keyFunc == nil {
		keyFunc = Nil
	}

	if reverse == nil {
		reverse = False
	}

	s := ptrSortable{&sortable{l, keyFunc, ObjectIsTrue(reverse), nil}}
	sort.Stable(s)
	return s.s.firstErr
}
