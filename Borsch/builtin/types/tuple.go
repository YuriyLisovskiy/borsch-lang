package types

import "fmt"

var TupleClass = ObjectClass.ClassNew("кортеж", map[string]Object{}, true, TupleNew, nil)

type Tuple []Object

func (value *Tuple) Class() *Class {
	return TupleClass
}

func TupleNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	// TODO: add iterators!
	tuple := &Tuple{}
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case *List:
			*tuple = arg.Values
		default:
			return nil, ErrorNewf("об'єкт '%s' не є ітерованим", arg.Class().Name)
		}
	} else if len(args) > 1 {
		*tuple = args
	}

	return tuple, nil
}

func (value *Tuple) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value *Tuple) string(ctx Context) (Object, error) {
	str := String("")
	vLen := len(*value)
	for i, item := range *value {
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

func (value *Tuple) Length(_ Context) (Int, error) {
	return Int(len(*value)), nil
}

func (value *Tuple) GetElement(ctx Context, index Int) (Object, error) {
	length, err := value.Length(ctx)
	if err != nil {
		return nil, err
	}

	if err = checkIndex(index, length, "кортежу"); err != nil {
		return nil, err
	}

	return (*value)[index], nil
}

func (value *Tuple) SetElement(_ Context, _ Int, _ Object) (Object, error) {
	return nil, TypeErrorNewf("об'єкт з типом 'кортеж' не підтримує присвоєння елементів за індексом")
}

func (value *Tuple) Slice(ctx Context, leftBound, rightBound Int) (Object, error) {
	length, err := value.Length(ctx)
	if err != nil {
		return nil, err
	}

	if err = checkIndex(leftBound, length, "кортежу"); err != nil {
		if leftBound < 0 {
			leftBound = 0
		}
	}

	if err = checkIndex(rightBound, length+1, "кортежу"); err != nil {
		if rightBound > length {
			rightBound = length
		}
	}

	list := NewList()
	if leftBound > rightBound {
		return list, nil
	}

	slicedList := (*value)[leftBound:rightBound]
	list.Values = make([]Object, len(slicedList))
	copy(list.Values, slicedList)
	return list, nil
}
