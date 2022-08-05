package types

func ShiftLeft(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IShiftLeft); ok {
		return v.shiftLeft(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedShiftLeft); ok {
			return v.reversedShiftLeft(ctx, a)
		}
	}

	return nil, ErrorNewf(
		"непідтримувані типи операндів для <<: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func ShiftRight(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IShiftRight); ok {
		return v.shiftRight(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedShiftRight); ok {
			return v.reversedShiftRight(ctx, a)
		}
	}

	return nil, ErrorNewf(
		"непідтримувані типи операндів для >>: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseOr(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IBitwiseOr); ok {
		return v.BitwiseOr(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedBitwiseOr); ok {
			return v.reversedBitwiseOr(ctx, a)
		}
	}

	return nil, ErrorNewf(
		"непідтримувані типи операндів для |: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseXor(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IBitwiseXor); ok {
		return v.BitwiseXor(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedBitwiseXor); ok {
			return v.reversedBitwiseXor(ctx, a)
		}
	}

	return nil, ErrorNewf(
		"непідтримувані типи операндів для ^: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func BitwiseAnd(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IBitwiseAnd); ok {
		return v.BitwiseAnd(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IReversedBitwiseAnd); ok {
			return v.reversedBitwiseAnd(ctx, a)
		}
	}

	return nil, ErrorNewf(
		"непідтримувані типи операндів для &: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Invert(ctx Context, a Object) (Object, error) {
	if v, ok := a.(IInvert); ok {
		return v.invert(ctx)
	}

	return nil, ErrorNewf(
		"непідтримуваний тип операнда для унарного ~: '%s'",
		a.Class().Name,
	)
}
