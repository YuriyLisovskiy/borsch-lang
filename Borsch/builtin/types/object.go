package types

// ObjectIsTrue returns whether the object is True or not.
func ObjectIsTrue(o Object) bool {
	if o == True {
		return true
	}

	if o == False {
		return false
	}

	if o == Nil {
		return false
	}

	if I, ok := o.(I__bool__); ok {
		cmp, err := I.__bool__()
		if err == nil && cmp == True {
			return true
		} else if err == nil && cmp == False {
			return false
		}
	}

	return false
}
