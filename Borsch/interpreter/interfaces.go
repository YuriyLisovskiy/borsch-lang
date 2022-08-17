package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
)

type Parser interface {
	Parse(filename string, code string) (Evaluatable, error)
}

type Interpreter interface {
	Import(packageName string) (types.Object, error)
	StackTrace() *common.StackTrace
}

type OperatorEvaluatable interface {
	Evaluate(State, types.Object) (types.Object, error)
}

type Evaluatable interface {
	Evaluate(State) (types.Object, error)
}

type State interface {
	Parent() State
	NewChild() State
	Context() types.Context
	Package() types.Object
	StackTrace() *common.StackTrace
	PackageOrNil() types.Object
	WithContext(types.Context) State
	WithPackage(types.Object) State
	RuntimeError(message string, statement common.Statement) error
	Trace(statement common.Statement, place string)
	PopTrace()
}
