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

func (value NilType) toBool(_ Context) (Object, error) {
	return False, nil
}

func (value NilType) equals(_ Context, other Object) (Object, error) {
	if _, ok := other.(NilType); ok {
		return True, nil
	}

	return False, nil
}

func (value NilType) notEquals(_ Context, other Object) (Object, error) {
	if _, ok := other.(NilType); ok {
		return False, nil
	}

	return True, nil
}
