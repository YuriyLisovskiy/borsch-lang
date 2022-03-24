package types

var (
	BoolClass = ObjectClass.ClassNew("логічний", map[string]Object{}, true, BoolNew, nil)

	True  = Bool(true)
	False = Bool(false)
)

type Bool bool

func (value Bool) Class() *Class {
	return BoolClass
}

func (value Bool) represent(ctx Context) (Object, error) {
	return value.string(ctx)
}

func (value Bool) string(Context) (Object, error) {
	if value {
		return String("істина"), nil
	}

	return String("хиба"), nil
}

func NewBool(value bool) Bool {
	if value {
		return True
	}

	return False
}

func BoolNew(ctx Context, cls *Class, args Tuple) (Object, error) {
	if len(args) != 1 {
		return nil, ErrorNewf("логічний() приймає 1 аргумент")
	}

	return ToBool(ctx, args[0])
}
