package types

import "github.com/YuriyLisovskiy/borsch-lang/Borsch/common"

// SequenceTuple converts a sequence object v into a Tuple.
func SequenceTuple(v Object) (Tuple, error) {
	switch x := v.(type) {
	case Tuple:
		return x, nil
	case *List:
		return Tuple(x.Items).Copy(), nil
	default:
		t := Tuple{}
		err := Iterate(
			v, func(item Object) bool {
				t = append(t, item)
				return false
			},
		)
		if err != nil {
			return nil, err
		}

		return t, nil
	}
}

// SequenceList converts a sequence object v into a List.
func SequenceList(v Object) (*List, error) {
	switch x := v.(type) {
	case Tuple:
		return NewListFromItems(x), nil
	case *List:
		return x.Copy(), nil
	default:
		l := NewList()
		err := l.ExtendSequence(v)
		if err != nil {
			return nil, err
		}

		return l, nil
	}
}

// Next calls __наступний__ for the Borsch object.
//
// Returns the next object.
//
// err == StopIteration or subclass when finished.
func Next(self Object) (obj Object, err error) {
	if I, ok := self.(I__next__); ok {
		return I.__next__()
	} else if obj, ok, err = TypeCall0(self, common.NextOperator); ok {
		return obj, err
	}

	return nil, ErrorNewf(TypeError, "'%s' не є ітерованим об'єктом", self.Type().Name)
}

func Iterate(obj Object, fn func(Object) bool) error {
	// Some easy cases
	switch x := obj.(type) {
	case Tuple:
		for _, item := range x {
			if fn(item) {
				break
			}
		}
	case *List:
		for _, item := range x.Items {
			if fn(item) {
				break
			}
		}
	case String:
		for _, item := range x {
			if fn(String(item)) {
				break
			}
		}
	default:
		iterator, err := Iter(obj)
		if err != nil {
			return err
		}
		for {
			item, err := Next(iterator)
			if err == StopIteration {
				break
			}
			if err != nil {
				return err
			}
			if fn(item) {
				break
			}
		}
	}

	return nil
}
