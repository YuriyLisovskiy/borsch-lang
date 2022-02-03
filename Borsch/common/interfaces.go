package common

import "github.com/alecthomas/participle/v2/lexer"

type Parser interface {
	Parse(filename string, code string) (Evaluatable, error)
}

type Interpreter interface {
	Import(state State, packageName string) (Type, error)
	Trace(pos lexer.Position, place string, statement string)
	StackTrace() *StackTrace
}

type Context interface {
	PushScope(scope map[string]Type)
	PopScope() map[string]Type
	TopScope() map[string]Type
	GetVar(name string) (Type, error)
	SetVar(name string, value Type) error
	GetClass(name string) (Type, error)
	GetChild() Context
}

type State interface {
	GetParser() Parser
	GetInterpreter() Interpreter
	GetContext() Context
	GetCurrentPackage() Type
	GetCurrentPackageOrNil() Type
	WithContext(Context) State
	WithPackage(p Type) State
}

type Evaluatable interface {
	Evaluate(State) (Type, error)
}

type OperatorEvaluatable interface {
	Evaluate(State, Type) (Type, error)
}

type Type interface {
	String(State) (string, error)
	Representation(State) (string, error)
	AsBool(State) (bool, error)
	GetTypeName() string
	GetOperator(string) (Type, error)
	GetAttribute(string) (Type, error)
	SetAttribute(string, Type) error
	HasAttribute(string) bool
}

type SequentialType interface {
	Length(State) int64
	GetElement(State, int64) (Type, error)
	SetElement(State, int64, Type) (Type, error)
	Slice(State, int64, int64) (Type, error)
}

type CallableType interface {
	Call(State, *[]Type, *map[string]Type) (Type, error)
}
