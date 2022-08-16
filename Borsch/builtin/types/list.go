package types

import "fmt"

var ListClass = ObjectClass.ClassNew("список", map[string]Object{}, true, ListNew, nil)

type List struct {
	Values []Object
}

func checkIndex(index, length Int, indexOfWhat string) error {
	if index >= 0 && index < length {
		return nil
	}

	return NewIndexOutOfRangeErrorf("індекс %s за межами діапазону", indexOfWhat)
}

func NewList() *List {
	return &List{Values: nil}
}

func (value *List) Class() *Class {
	return ListClass
}

func ListNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	return &List{Values: args}, nil
}

func (value *List) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *List) string(ctx Context) (Object, error) {
	str := String("")
	vLen := len(value.Values)
	for i, item := range value.Values {
		itemStr, err := ToString(ctx, item)
		if err != nil {
			return nil, err
		}

		str += itemStr.(String)
		if i < vLen-1 {
			str += ", "
		}
	}

	return String(fmt.Sprintf("[%s]", str)), nil
}

func (value *List) Length(_ Context) (Int, error) {
	return Int(len(value.Values)), nil
}

func (value *List) GetElement(ctx Context, index Int) (Object, error) {
	length, err := value.Length(ctx)
	if err != nil {
		return nil, err
	}

	if err = checkIndex(index, length, "списку"); err != nil {
		return nil, err
	}

	return value.Values[index], nil
}

func (value *List) SetElement(ctx Context, index Int, item Object) (Object, error) {
	length, err := value.Length(ctx)
	if err != nil {
		return nil, err
	}

	if err = checkIndex(index, length, "списку"); err != nil {
		return nil, err
	}

	value.Values[index] = item
	return nil, nil
}

func (value *List) Slice(ctx Context, leftBound, rightBound Int) (Object, error) {
	length, err := value.Length(ctx)
	if err != nil {
		return nil, err
	}

	if err = checkIndex(leftBound, length, "списку"); err != nil {
		if leftBound < 0 {
			leftBound = 0
		}
	}

	if err = checkIndex(rightBound, length+1, "списку"); err != nil {
		if rightBound > length {
			rightBound = length
		}
	}

	list := NewList()
	if leftBound > rightBound {
		return list, nil
	}

	slicedList := value.Values[leftBound:rightBound]
	list.Values = make([]Object, len(slicedList))
	copy(list.Values, slicedList)
	return list, nil
}
