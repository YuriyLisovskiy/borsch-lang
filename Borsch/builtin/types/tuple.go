package types

import "bytes"

var TupleType = ObjectType.NewType(
	"кортеж",
	`кортеж() -> порожній кортеж
кортеж(ітерований_об_єкт) -> кортеж, ініціалізований з елементів ітерованого об'єкта

Якщо аргумент є кортежем, результатом буде цей ж об'єкт.`,
	TupleNew,
	nil,
)

type Tuple []Object

func (t Tuple) Type() *Type {
	return TupleType
}

func TupleNew(cls *Type, args Tuple, kwargs StringDict) (res Object, err error) {
	var iterable Object
	err = UnpackTuple(args, kwargs, "кортеж", 0, 1, &iterable)
	if err != nil {
		return nil, err
	}
	if iterable != nil {
		return SequenceTuple(iterable)
	}

	return Tuple{}, nil
}

func (t Tuple) Copy() Tuple {
	newT := make(Tuple, len(t))
	copy(newT, t)
	return newT
}

// output the tuple to out, using fn to transform the tuple to out
// start and end brackets
func (t Tuple) represent(start, end string) (Object, error) {
	var out bytes.Buffer
	out.WriteString(start)
	for i, obj := range t {
		if i != 0 {
			out.WriteString(", ")
		}

		str, err := RepresentAsString(obj)
		if err != nil {
			return nil, err
		}

		out.WriteString(str)
	}

	out.WriteString(end)
	return String(out.String()), nil
}

func (t Tuple) __str__() (Object, error) {
	return t.__represent__()
}

func (t Tuple) __represent__() (Object, error) {
	return t.represent("(", ")")
}

func (t Tuple) __length__() (Object, error) {
	return Int(len(t)), nil
}

func (t Tuple) __bool__() (Object, error) {
	return NewBool(len(t) > 0), nil
}

func (t Tuple) __add__(other Object) (Object, error) {
	if b, ok := other.(Tuple); ok {
		newTuple := make(Tuple, len(t)+len(b))
		copy(newTuple, t)
		copy(newTuple[len(b):], b)
		return newTuple, nil
	}

	return NotImplemented, nil
}

func (t Tuple) __reversed_add__(other Object) (Object, error) {
	if b, ok := other.(Tuple); ok {
		return b.__add__(t)
	}

	return NotImplemented, nil
}

func (t Tuple) __in_place_add__(other Object) (Object, error) {
	return t.__add__(other)
}

func (t Tuple) __mul__(other Object) (Object, error) {
	if b, ok := convertToInt(other); ok {
		m := len(t)
		n := int(b) * m
		if n < 0 {
			n = 0
		}

		newTuple := make(Tuple, n)
		for i := 0; i < n; i += m {
			copy(newTuple[i:i+m], t)
		}

		return newTuple, nil
	}

	return NotImplemented, nil
}

func (t Tuple) __reversed_mul__(other Object) (Object, error) {
	return t.__mul__(other)
}

func (t Tuple) __in_place_mul__(other Object) (Object, error) {
	return t.__mul__(other)
}

func (t Tuple) __equal__(other Object) (Object, error) {
	b, ok := other.(Tuple)
	if !ok {
		return NotImplemented, nil
	}

	if len(t) != len(b) {
		return False, nil
	}

	for i := range t {
		eq, err := Equal(t[i], b[i])
		if err != nil {
			return nil, err
		}

		if eq == False {
			return False, nil
		}
	}

	return True, nil
}

func (t Tuple) __not_equal__(other Object) (Object, error) {
	b, ok := other.(Tuple)
	if !ok {
		return NotImplemented, nil
	}

	if len(t) != len(b) {
		return True, nil
	}

	for i := range t {
		eq, err := Equal(t[i], b[i])
		if err != nil {
			return nil, err
		}

		if eq == False {
			return True, nil
		}
	}

	return False, nil
}

// Check interface is satisfied
var _ sequenceArithmetic = Tuple(nil)
var _ I__str__ = Tuple(nil)
var _ I__represent__ = Tuple(nil)
var _ I__length__ = Tuple(nil)
var _ I__bool__ = Tuple(nil)
var _ I__equal__ = Tuple(nil)
var _ I__not_equal__ = Tuple(nil)
