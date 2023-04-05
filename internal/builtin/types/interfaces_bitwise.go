package types

type IShiftLeft interface {
	shiftLeft(ctx Context, other Object) (Object, error)
}

type IReversedShiftLeft interface {
	reversedShiftLeft(ctx Context, other Object) (Object, error)
}

type IShiftRight interface {
	shiftRight(ctx Context, other Object) (Object, error)
}

type IReversedShiftRight interface {
	reversedShiftRight(ctx Context, other Object) (Object, error)
}

type IBitwiseOr interface {
	bitwiseOr(ctx Context, other Object) (Object, error)
}

type IReversedBitwiseOr interface {
	reversedBitwiseOr(ctx Context, other Object) (Object, error)
}

type IBitwiseXor interface {
	bitwiseXor(ctx Context, other Object) (Object, error)
}

type IReversedBitwiseXor interface {
	reversedBitwiseXor(ctx Context, other Object) (Object, error)
}

type IBitwiseAnd interface {
	bitwiseAnd(ctx Context, other Object) (Object, error)
}

type IReversedBitwiseAnd interface {
	reversedBitwiseAnd(ctx Context, other Object) (Object, error)
}

type IInvert interface {
	invert(ctx Context) (Object, error)
}
