package types

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Parser interface {
	Parse(filename string, code string) (Evaluatable, error)
}

type Interpreter interface {
	Import(state State, packageName string) (Object, error)
	StackTrace() *StackTrace
}

type Evaluatable interface {
	Evaluate(State) (Object, error)
}

type State interface {
	GetParser() Parser
	GetInterpreter() Interpreter
	GetContext() Context
	GetCurrentPackage() Object
	GetCurrentPackageOrNil() Object
	WithContext(Context) State
	WithPackage(Object) State
	RuntimeError(message string, statement Statement) error
	Trace(statement Statement, place string)
	PopTrace()
}

type Context interface {
	PushScope(scope map[string]Object)
	PopScope() map[string]Object
	TopScope() map[string]Object
	GetVar(name string) (Object, error)
	SetVar(name string, value Object) error
	GetClass(name string) (Object, error)
	Derive() Context
}

type Statement interface {
	String() string
	Position() lexer.Position
}
