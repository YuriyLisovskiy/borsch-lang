package types

type Object interface {
	Class() *Class
}

type IString interface {
	string(ctx Context) (Object, error)
}

type IRepresent interface {
	represent(ctx Context) (Object, error)
}

type IBool interface {
	toBool(ctx Context) (Object, error)
}

type IInt interface {
	toInt(ctx Context) (Object, error)
}

type IReal interface {
	toReal(ctx Context) (Object, error)
}

type IGoInt interface {
	toGoInt(ctx Context) (int, error)
}

type ICall interface {
	call(args Tuple) (Object, error)
}

type IGetAttribute interface {
	getAttribute(ctx Context, name string) (Object, error)
}

type ISetAttribute interface {
	setAttribute(ctx Context, name string, value Object) error
}

type IDeleteAttribute interface {
	deleteAttribute(ctx Context, name string) (Object, error)
}

type IAdd interface {
	add(ctx Context, other Object) (Object, error)
}

type IReversedAdd interface {
	reversedAdd(ctx Context, other Object) (Object, error)
}

type ISub interface {
	sub(ctx Context, other Object) (Object, error)
}

type IReversedSub interface {
	reversedSub(ctx Context, other Object) (Object, error)
}

type IDiv interface {
	div(ctx Context, other Object) (Object, error)
}

type IReversedDiv interface {
	reversedDiv(ctx Context, other Object) (Object, error)
}

type IMul interface {
	mul(ctx Context, other Object) (Object, error)
}

type IReversedMul interface {
	reversedMul(ctx Context, other Object) (Object, error)
}

type IMod interface {
	mod(ctx Context, other Object) (Object, error)
}

type IReversedMod interface {
	reversedMod(ctx Context, other Object) (Object, error)
}

type IPow interface {
	pow(ctx Context, other Object) (Object, error)
}

type IReversedPow interface {
	reversedPow(ctx Context, other Object) (Object, error)
}

type IPositive interface {
	positive(ctx Context) (Object, error)
}

type INegate interface {
	negate(ctx Context) (Object, error)
}
