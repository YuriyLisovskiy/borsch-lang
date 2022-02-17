package common

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

type Context interface {
	PushScope(scope map[string]Object)
	PopScope() map[string]Object
	TopScope() map[string]Object
	GetVar(name string) (Object, error)
	SetVar(name string, value Object) error
	GetClass(name string) (Object, error)
	Derive() Context
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

type Statement interface {
	String() string
	Position() lexer.Position
}

type Evaluatable interface {
	Evaluate(State) (Object, error)
}

type OperatorEvaluatable interface {
	Evaluate(State, Object) (Object, error)
}

type Object interface {
	String(State) (string, error)
	Representation(State) (string, error)
	AsBool(State) (bool, error)
	GetTypeName() string
	GetOperator(string) (Object, error)
	GetAttribute(string) (Object, error)
	SetAttribute(string, Object) error
	HasAttribute(string) bool
}

type SequentialType interface {
	Length(State) int64
	GetElement(State, int64) (Object, error)
	SetElement(State, int64, Object) (Object, error)
	Slice(State, int64, int64) (Object, error)
}

type CallableType interface {
	Call(State, *[]Object, *map[string]Object) (Object, error)
}
