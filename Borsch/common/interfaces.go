package common

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type Parser interface {
	Parse(filename string, code string) (Evaluatable, error)
}

type Interpreter interface {
	Import(state State, packageName string) (Value, error)
	Trace(pos lexer.Position, place string, statement string)
	StackTrace() *StackTrace
}

type Context interface {
	PushScope(scope map[string]Value)
	PopScope() map[string]Value
	TopScope() map[string]Value
	GetVar(name string) (Value, error)
	SetVar(name string, value Value) error
	GetClass(name string) (Value, error)
	GetChild() Context
}

type State interface {
	GetParser() Parser
	GetInterpreter() Interpreter
	GetContext() Context
	GetCurrentPackage() Value
	GetCurrentPackageOrNil() Value
	WithContext(Context) State
	WithPackage(Value) State
	RuntimeError(string, Statement) error
}

type Statement interface {
	String() string
	Position() lexer.Position
}

type Evaluatable interface {
	Evaluate(State) (Value, error)
}

type OperatorEvaluatable interface {
	Evaluate(State, Value) (Value, error)
}

type Value interface {
	String(State) (string, error)
	Representation(State) (string, error)
	AsBool(State) (bool, error)
	GetTypeName() string
	GetOperator(string) (Value, error)
	GetAttribute(string) (Value, error)
	SetAttribute(string, Value) error
	HasAttribute(string) bool
}

type SequentialType interface {
	Length(State) int64
	GetElement(State, int64) (Value, error)
	SetElement(State, int64, Value) (Value, error)
	Slice(State, int64, int64) (Value, error)
}

type CallableType interface {
	Call(State, *[]Value, *map[string]Value) (Value, error)
}
