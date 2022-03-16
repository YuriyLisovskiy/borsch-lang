package interpreter

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
)

type StateImpl struct {
	parser         types.Parser
	interpreter    types.Interpreter
	context        types.Context
	currentPackage types.Object
}

func NewState(
	parser types.Parser,
	interpreter types.Interpreter,
	context types.Context,
	currentPackage *types.Package,
) *StateImpl {
	return &StateImpl{
		parser:         parser,
		interpreter:    interpreter,
		context:        context,
		currentPackage: currentPackage,
	}
}

func (s *StateImpl) GetParser() types.Parser {
	if s.parser != nil {
		return s.parser
	}

	panic("state: parser is nil")
}

func (s *StateImpl) GetInterpreter() types.Interpreter {
	if s.interpreter != nil {
		return s.interpreter
	}

	panic("state: interpreter is nil")
}

func (s *StateImpl) GetContext() types.Context {
	if s.context != nil {
		return s.context
	}

	panic("state: context is nil")
}

func (s *StateImpl) GetCurrentPackage() types.Object {
	if s.currentPackage != nil {
		return s.currentPackage
	}

	panic("state: current package is nil")
}

func (s *StateImpl) GetCurrentPackageOrNil() types.Object {
	return s.currentPackage
}

func (s *StateImpl) WithContext(ctx types.Context) types.State {
	return &StateImpl{
		parser:         s.parser,
		interpreter:    s.interpreter,
		context:        ctx,
		currentPackage: s.currentPackage,
	}
}

func (s *StateImpl) WithPackage(pkg types.Object) types.State {
	return &StateImpl{
		parser:         s.parser,
		interpreter:    s.interpreter,
		context:        s.context,
		currentPackage: pkg,
	}
}

func (s *StateImpl) RuntimeError(message string, statement types.Statement) error {
	if statement != nil {
		s.Trace(statement, "")
	}

	return errors.New(fmt.Sprintf("Помилка виконання: %s", message))
}

func (s *StateImpl) Trace(statement types.Statement, place string) {
	s.interpreter.StackTrace().Push(types.NewTraceRow(statement.Position(), statement.String(), place))
}

func (s *StateImpl) PopTrace() {
	s.interpreter.StackTrace().Pop()
}
