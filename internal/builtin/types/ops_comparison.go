package types

func Equals(ctx Context, a, b Object) (Object, error) {
	if v, ok := a.(IEquals); ok {
		result, err := v.equals(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IEquals); ok {
			result, err := v.equals(ctx, a)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
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
		result, err := v.notEquals(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	if a.Class() != b.Class() {
		if v, ok := b.(INotEquals); ok {
			result, err := v.notEquals(ctx, a)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
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
		result, err := v.less(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	if a.Class() != b.Class() {
		if v, ok := b.(ILess); ok {
			result, err := v.less(ctx, a)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
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
		result, err := v.lessOrEquals(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	if a.Class() != b.Class() {
		if v, ok := b.(ILessOrEquals); ok {
			result, err := v.lessOrEquals(ctx, a)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
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
		result, err := v.greater(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IGreater); ok {
			result, err := v.greater(ctx, a)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
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
		result, err := v.greaterOrEquals(ctx, b)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	if a.Class() != b.Class() {
		if v, ok := b.(IGreaterOrEquals); ok {
			result, err := v.greaterOrEquals(ctx, a)
			if err != nil {
				return nil, err
			}

			if result != nil {
				return result, nil
			}
		}
	}

	return nil, NewErrorf(
		"непідтримувані типи операндів для >=: '%s' та '%s'",
		a.Class().Name,
		b.Class().Name,
	)
}
