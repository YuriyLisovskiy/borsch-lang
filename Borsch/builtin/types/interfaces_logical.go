package types

type IAnd interface {
	And(ctx Context, other Object) (Object, error)
}

type IOr interface {
	Or(ctx Context, other Object) (Object, error)
}

type INot interface {
	Not(ctx Context) (Object, error)
}
