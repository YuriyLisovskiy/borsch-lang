package types

func Equals(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IEquals); ok {
		return v.equals(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IEquals); ok {
			return v.equals(ctx, a)
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для ==: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func NotEquals(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(INotEquals); ok {
		return v.notEquals(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(INotEquals); ok {
			return v.notEquals(ctx, a)
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для !=: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Less(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(ILess); ok {
		return v.less(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(ILess); ok {
			return v.less(ctx, a)
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для <: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func LessOrEquals(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(ILessOrEquals); ok {
		return v.lessOrEquals(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(ILessOrEquals); ok {
			return v.lessOrEquals(ctx, a)
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для <=: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func Greater(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IGreater); ok {
		return v.greater(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IGreater); ok {
			return v.greater(ctx, a)
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для >: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}

func GreaterOrEquals(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IGreaterOrEquals); ok {
		return v.greaterOrEquals(ctx, b)
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IGreaterOrEquals); ok {
			return v.greaterOrEquals(ctx, a)
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для >=: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}
