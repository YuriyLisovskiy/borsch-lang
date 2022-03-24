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

type State interface {
	GetParser() Parser
	GetInterpreter() Interpreter
	GetContext() types.Context
	GetCurrentPackage() types.Object
	GetCurrentPackageOrNil() types.Object
	WithContext(types.Context) State
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

type SequentialType interface {
	Length(State) int64
	GetElement(State, int64) (types.Object, error)
	SetElement(State, int64, types.Object) (types.Object, error)
	Slice(State, int64, int64) (types.Object, error)
}

type CallableType interface {
	Call(State, *[]types.Object, *map[string]types.Object) (types.Object, error)
}
