package common

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/alecthomas/participle/v2/lexer"
)

type Parser interface {
	Parse(filename string, code string) (Evaluatable, error)
}

type Interpreter interface {
	Import(state State, packageName string) (types.Object, error)
	StackTrace() *StackTrace
}

type Context interface {
	PushScope(scope map[string]types.Object)
	PopScope() map[string]types.Object
	TopScope() map[string]types.Object
	GetVar(name string) (types.Object, error)
	SetVar(name string, value types.Object) error
	GetClass(name string) (types.Object, error)
	Derive() Context
}

type State interface {
	GetParser() Parser
	GetInterpreter() Interpreter
	GetContext() Context
	GetCurrentPackage() types.Object
	GetCurrentPackageOrNil() types.Object
	WithContext(Context) State
	WithPackage(types.Object) State
	RuntimeError(message string, statement Statement) error
	Trace(statement Statement, place string)
	PopTrace()
}

type Statement interface {
	String() string
	Position() lexer.Position
}

type Evaluatable interface {
	Evaluate(State) (types.Object, error)
}

type OperatorEvaluatable interface {
	Evaluate(State, types.Object) (types.Object, error)
}
