package types

var TupleClass = ObjectClass.ClassNew("кортеж", map[string]Object{}, true, TupleNew, nil)

type Tuple []Object

func (value Tuple) Class() *Class {
	return TupleClass
}

func TupleNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	// TODO: add iterators!
	tuple := Tuple{}
	if len(args) == 1 {
		switch arg := args[0].(type) {
		case *List:
			tuple = arg.Values
		default:
			return nil, ErrorNewf("об'єкт '%s' не є ітерованим", arg.Class().Name)
		}
	} else if len(args) > 1 {
		tuple = args
	}

	return tuple, nil
}
