package types

func Add(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IAdd); ok {
		return v.add(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedAdd); ok {
			return v.reversedAdd(ctx, a)
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
		return v.sub(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedSub); ok {
			return v.reversedSub(ctx, a)
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
		return v.div(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedDiv); ok {
			return v.reversedDiv(ctx, a)
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
		return v.mul(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedMul); ok {
			return v.reversedMul(ctx, a)
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
		return v.mod(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedMod); ok {
			return v.reversedMod(ctx, a)
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
		return v.pow(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedPow); ok {
			return v.reversedPow(ctx, a)
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
		return v.positive(ctx)
	}

	return nil, NewErrorf(
		"непідтримуваний тип операнда для унарного +: '%s'",
		a.Class().Name,
	)
}

func Negate(ctx Context, a Object) (Object, error) {
	if v, ok := a.(INegate); ok {
		return v.negate(ctx)
	}

	return nil, NewErrorf(
		"непідтримуваний тип операнда для унарного -: '%s'",
		a.Class().Name,
	)
}
