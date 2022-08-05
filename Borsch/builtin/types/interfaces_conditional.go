package types

type IEquals interface {
	equals(ctx Context, other Object) (Object, error)
}

type INotEquals interface {
	notEquals(ctx Context, other Object) (Object, error)
}

type ILess interface {
	less(ctx Context, other Object) (Object, error)
}

type ILessOrEquals interface {
	lessOrEquals(ctx Context, other Object) (Object, error)
}

type IGreater interface {
	greater(ctx Context, other Object) (Object, error)
}

type IGreaterOrEquals interface {
	greaterOrEquals(ctx Context, other Object) (Object, error)
}
