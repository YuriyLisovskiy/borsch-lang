package types

func And(ctx Context, a, b Object) (Object, error) {
	aBool, err := ToBool(ctx, a)
	if err != nil {
		return nil, err
	}

	bBool, err := ToBool(ctx, b)
	if err != nil {
		return nil, err
	}

	return aBool.(Bool) && bBool.(Bool), nil
}

func Or(ctx Context, a, b Object) (Object, error) {
	aBool, err := ToBool(ctx, a)
	if err != nil {
		return nil, err
	}

	bBool, err := ToBool(ctx, b)
	if err != nil {
		return nil, err
	}

	return aBool.(Bool) || bBool.(Bool), nil
}

func Not(ctx Context, a Object) (Object, error) {
	return ToBool(ctx, a)
}
