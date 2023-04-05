package types

func Add(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IAdd); ok {
		result, err := v.add(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для +: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Sub(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(ISub); ok {
		result, err := v.sub(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для -: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Div(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IDiv); ok {
		result, err := v.div(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для /: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Mul(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IMul); ok {
		result, err := v.mul(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для *: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Mod(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IMod); ok {
		result, err := v.mod(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для %%: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Pow(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IPow); ok {
		result, err := v.pow(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для **: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Positive(ctx Context, a Object) (Object, error) {
	if v, ok := a.(IPositive); ok {
		result, err := v.positive(ctx)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримуваний тип операнда для унарного +: '%s'",
		a.Class().Name,
	)
}

func Negate(ctx Context, a Object) (Object, error) {
	if v, ok := a.(INegate); ok {
		result, err := v.negate(ctx)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, NewErrorf(
		"непідтримуваний тип операнда для унарного -: '%s'",
		a.Class().Name,
	)
}
