package interpreter

import (
	types2 "github.com/YuriyLisovskiy/borsch-lang/internal/builtin/types"
	common2 "github.com/YuriyLisovskiy/borsch-lang/internal/common"
)

type Parser interface {
	Parse(filename string, code string) (Evaluatable, error)
}

type Interpreter interface {
	Import(packageName string) (types2.Object, error)
	Evaluate(packageName, code string, parentPkg *types2.Package) (types2.Object, error)
	StackTrace() *common2.StackTrace
}

type OperatorEvaluatable interface {
	Evaluate(State, types2.Object) (types2.Object, error)
}

type Evaluatable interface {
	Evaluate(State) (types2.Object, error)
}

type State interface {
	Parent() State
	NewChild() State
	Context() types2.Context
	Package() types2.Object
	StackTrace() *common2.StackTrace
	PackageOrNil() types2.Object
	WithContext(types2.Context) State
	WithPackage(types2.Object) State
	RuntimeError(message string, statement common2.Statement) error
	Trace(statement common2.Statement, place string)
	PopTrace()
}
