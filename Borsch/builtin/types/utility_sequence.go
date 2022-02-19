package types

// SequenceTuple converts a sequence object v into a Tuple.
func SequenceTuple(v Object) (Tuple, error) {
	switch x := v.(type) {
	case Tuple:
		return x, nil
	case *List:
		return Tuple(x.Items).Copy(), nil
	default:
		panic("unreachable")
		// TODO:
		// t := Tuple{}
		// err := Iterate(
		// 	v, func(item Object) bool {
		// 		t = append(t, item)
		// 		return false
		// 	},
		// )
		// if err != nil {
		// 	return nil, err
		// }
		//
		// return t, nil
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
