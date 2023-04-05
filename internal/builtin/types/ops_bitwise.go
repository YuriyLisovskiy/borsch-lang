package types

func ShiftLeft(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IShiftLeft); ok {
		result, err := v.shiftLeft(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для <<: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func ShiftRight(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IShiftRight); ok {
		result, err := v.shiftRight(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для >>: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseOr(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IBitwiseOr); ok {
		result, err := v.bitwiseOr(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для |: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseXor(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IBitwiseXor); ok {
		result, err := v.bitwiseXor(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для ^: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseAnd(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IBitwiseAnd); ok {
		result, err := v.bitwiseAnd(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для &: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Invert(ctx Context, a Object) (Object, error) {
	if v, ok := a.(IInvert); ok {
		result, err := v.invert(ctx)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримуваний тип операнда для унарного ~: '%s'",
		a.Class().Name,
	)
}
