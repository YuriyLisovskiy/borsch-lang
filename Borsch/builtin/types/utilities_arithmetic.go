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

	return nil, ErrorNewf(
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

	return nil, ErrorNewf(
		"непідтримувані типи операндів для -: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Div(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
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

	return nil, ErrorNewf(
		"непідтримувані типи операндів для *: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Mod(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для %: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Pow(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для **: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func ShiftLeft(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для <<: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func ShiftRight(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для >>: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseOr(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для |: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseXor(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для ^: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseAnd(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для &: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func And(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для &&: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Or(ctx Context, a, b Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримувані типи операндів для ||: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Positive(ctx Context, a Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримуваний тип операнда для унарного +: '%s'",
		a.Class().Name,
	)
}

func Negate(ctx Context, a Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримуваний тип операнда для унарного -: '%s'",
		a.Class().Name,
	)
}

func Invert(ctx Context, a Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримуваний тип операнда для унарного ~: '%s'",
		a.Class().Name,
	)
}

func Not(ctx Context, a Object) (Object, error) {
	return nil, ErrorNewf(
		"непідтримуваний тип операнда для унарного !: '%s'",
		a.Class().Name,
	)
}
