package common

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

type OperatorEvaluatable interface {
	Evaluate(types.State, types.Object) (types.Object, error)
}
