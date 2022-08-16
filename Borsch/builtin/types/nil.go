package types

var (
	NilClass = ObjectClass.ClassNew("нульове", map[string]Object{}, true, nil, nil)

	Nil = NilType{}
)

type NilType struct {
}

func (value NilType) Class() *Class {
	return NilClass
}

func (value NilType) represent(Context) (Object, error) {
	return String("нуль"), nil
}

func (value NilType) string(ctx Context) (Object, error) {
	return value.represent(ctx)
}
