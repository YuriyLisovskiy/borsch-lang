package types

var (
	NilClass = ObjectClass.ClassNew("нульовий", map[string]Object{}, true, nil, nil)

	Nil = nilType{}
)

type nilType struct {
}

func (value nilType) Class() *Class {
	return NilClass
}

func (value nilType) represent(Context) (Object, error) {
	return String("нуль"), nil
}

func (value nilType) string(ctx Context) (Object, error) {
	return value.represent(ctx)
}
