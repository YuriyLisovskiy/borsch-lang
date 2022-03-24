package types

var ListClass = ObjectClass.ClassNew("список", map[string]Object{}, true, ListNew, nil)

type List struct {
	Values []Object
}

func NewList() *List {
	return &List{Values: nil}
}

func (value *List) Class() *Class {
	return ListClass
}

func (value *List) Length() int64 {
	return int64(len(value.Values))
}

func ListNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	return &List{Values: args}, nil
}
